package signalPing

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const MTU = 1500

type packet struct {
	bytes  []byte
	nbytes int
	peer   net.Addr
}

type SignalPing struct {
	quit chan bool
}

func newService() *SignalPing {
	return &SignalPing{quit: make(chan bool)}
}

func NewService() *SignalPing {
	s := newService()
	return s
}

func (s *SignalPing) Run() {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	recv := make(chan *packet)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	wg.Add(1)
	go s.recvICMP(conn, recv, &wg)

	for {
		select {
		case <-s.quit:
			wg.Wait()
			fmt.Println("finishing task")
			time.Sleep(time.Second)
			fmt.Println("task done")
			s.quit <- true
			return
		case rb := <-recv:
			fmt.Println("process an icmp packet")
			dst, err := s.processPacket(rb)
			if err != nil {
				log.Fatal(err)
			}
			s.sendICMPReply(conn, dst)
		}
	}
}

func (s *SignalPing) Stop() {
	fmt.Println("SignalPing stopping")
	s.quit <- true
	<-s.quit
	fmt.Println("SignalPing stopped")
	close(s.quit)
}

func (s *SignalPing) recvICMP(
	conn *icmp.PacketConn,
	recv chan<- *packet,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-s.quit:
			return
		default:
			rb := make([]byte, MTU)
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
			n, peer, err := conn.ReadFrom(rb)
			if err != nil {
				if neterr, ok := err.(*net.OpError); ok {
					if neterr.Timeout() {
						continue
					} else {
						s.Stop()
						return
					}
				}
			}
			fmt.Println("received an icmp packet")
			recv <- &packet{bytes: rb, nbytes: n, peer: peer}
		}
	}
}

func (s *SignalPing) processPacket(recv *packet) (net.Addr, error) {

	rb := recv.bytes
	rm, err := icmp.ParseMessage(1, rb[:recv.nbytes])
	if err != nil {
		log.Fatal(err)
	}

	switch rm.Type {
	case ipv4.ICMPTypeEcho:
		fmt.Println("icmptypeecho")
		mb, err := rm.Body.Marshal(1)
		if err != nil {
			log.Fatal(err)
		}
		msg := string(mb)
		if msg[16:] == "SIGTERMING_REQUEST" {
			fmt.Println("streaming request get")
			handler := s.OnRecv
			if handler != nil {
				go handler(recv.peer)
			}
			return recv.peer, nil
		} else {
			return nil, fmt.Errorf("Error, invalid ICMP echo reply. Body message: %s", msg)
		}
	default:
		return nil, fmt.Errorf("Error, invalid ICMP echo reply. ICMP type: %T", rm.Type)
	}
	return nil, nil
}

func (s *SignalPing) sendICMPReply(conn *icmp.PacketConn, dstAddr net.Addr) error {
	wb, err := (&icmp.Message{
		Type: ipv4.ICMPTypeEchoReply,
		Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Intn(65535),
			Seq:  1,
			Data: []byte("STREAMING_REQUEST_ACK"),
		},
	}).Marshal(nil)
	if err != nil {
		return err
	}

	dst, err := net.ResolveUDPAddr("udp4", dstAddr.String())
	if _, err := conn.WriteTo(wb, dst); err != nil {
		log.Fatal(err)
	}

	return err
}
