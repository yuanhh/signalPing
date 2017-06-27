package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const MTU = 1500

func main() {
	switch runtime.GOOS {
	case "darwin":
	case "linux":
	default:
		log.Println("not supported on", runtime.GOOS)
		return
	}

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port\n", os.Args[0])
		os.Exit(1)
	}

	// listen for control message reply
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// construct a control message
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("STREAMING_REQUEST"),
		},
	}

	// the returned message always contains the calculated checksum field
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	dst, err := net.ResolveUDPAddr("udp4", os.Args[1])
	// send a control message to the endpoint
	if _, err := c.WriteTo(wb, dst); err != nil {
		log.Fatal("c.WriteTo:", err)
	}

	// receive icmp echo reply
	rb := make([]byte, MTU)
	n, peer, err := c.ReadFrom(rb)
	if err != nil {
		log.Fatal(err)
	}
	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		log.Fatal(err)
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		log.Printf("got reflection from %v", peer)
	default:
		log.Printf("got %+v; want echo reply", rm)
	}
}
