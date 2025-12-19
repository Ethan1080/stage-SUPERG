package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	device := "wlp58s0" // souvent eth0 ou wlan0
	snapshotLen := int32(1600)
	promiscuous := true
	timeout := pcap.BlockForever

	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	fmt.Println("Écoute du réseau sur", device)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		handlePacket(packet)
	}
}

func handlePacket(packet gopacket.Packet) {
	// IP
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}
	ip := ipLayer.(*layers.IPv4)

	// TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return
	}
	tcp := tcpLayer.(*layers.TCP)

	fmt.Printf(
		"\n📦 %s:%d → %s:%d\n",
		ip.SrcIP, tcp.SrcPort,
		ip.DstIP, tcp.DstPort,
	)

	if len(tcp.Payload) > 0 {
		fmt.Println("DATA:")
		fmt.Println(string(tcp.Payload))
	}
}
