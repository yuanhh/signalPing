package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const MTU = 1500

func main() {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("listen error", err)
		return
	}
	rb := make([]byte, 1500)

	for {
		n, raddr, err := c.ReadFrom(rb)
		if err != nil {
			log.Fatal(err)
		}
		rm, err := icmp.ParseMessage(1, rb[:n])
		if err != nil {
			log.Fatal(err)
			return
		}

		switch rm.Type {
		case ipv4.ICMPTypeEcho:
			fmt.Printf(raddr.String())
			// log.Printf(raddr.String())
			// log.Printf("got %+v; want echo reply", rm)
			// stream(raddr.String())
			return
		default:
		}
	}
}

func stream(addr string) {
	path, err := exec.LookPath("startStreamingFileSinkLinux.sh")
	if err != nil {
		log.Fatal("streaming path not found!")
	}
	fmt.Printf("startStreamingFileSinkLinux.sh is available at %s\n", path)

	host, port, err := net.SplitHostPort(addr)
	cmd := exec.Command("startStreamingFileSinkLinux.sh", "-f signal.mp4", "-h ", host, "-p", port)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}
