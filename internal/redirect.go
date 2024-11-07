package internal

import (
	"log"
	"net"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Redirect sends an ICMP redirect packet with a fake Ping packet to the specified destination
func Redirect(srcIP, dstIP string) {
	// Create a raw socket
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Fatalf("Failed to create raw socket: %v", err)
	}
	defer syscall.Close(socket)

	// Generate ICMP Redirect packets
	buf := gopacket.NewSerializeBuffer()
	ipLayer := &layers.IPv4{
		Version:  4,
		IHL:      5,
		SrcIP:    net.ParseIP(srcIP),
		DstIP:    net.ParseIP(dstIP),
		Protocol: layers.IPProtocolICMPv4,
		TTL:      64,
	}

	icmpRedirectLayer := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeRedirect, 0),
		Id:       0x1234,
		Seq:      1,
	}

	// Generate fake Ping (ICMP Echo Request) packet
	fakeIPLayer := &layers.IPv4{
		Version:  4,
		IHL:      5,
		SrcIP:    net.ParseIP("10.10.10.10"),
		DstIP:    net.ParseIP("20.20.20.20"),
		Protocol: layers.IPProtocolICMPv4,
		TTL:      64,
	}

	fakeICMPLayer := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       0x5678,
		Seq:      1,
	}

	// Serialize and write to buffer
	err = gopacket.SerializeLayers(buf, gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true},
		ipLayer,
		icmpRedirectLayer,
		fakeIPLayer,
		fakeICMPLayer,
		gopacket.Payload([]byte("")),
	)
	if err != nil {
		log.Fatalf("Failed to serialize packet: %v", err)
	}

	packetData := buf.Bytes()

	// Create of destination address structure
	addr := &syscall.SockaddrInet4{}
	copy(addr.Addr[:], ipLayer.DstIP.To4())

	// Send packets via raw socket
	err = syscall.Sendto(socket, packetData, 0, addr)
	if err != nil {
		log.Fatalf("Failed to send packet: %v", err)
	}

	log.Printf("Sent ICMP Redirect packet to %s", dstIP)
}
