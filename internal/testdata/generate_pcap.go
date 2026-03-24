//go:build ignore

// This program generates a minimal valid pcap file for testing.
// Run: go run internal/testdata/generate_pcap.go
package main

import (
	"encoding/binary"
	"os"
)

func main() {
	f, err := os.Create("testdata/sample.pcap")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Global header (24 bytes)
	binary.Write(f, binary.LittleEndian, uint32(0xa1b2c3d4)) // magic
	binary.Write(f, binary.LittleEndian, uint16(2))           // version major
	binary.Write(f, binary.LittleEndian, uint16(4))           // version minor
	binary.Write(f, binary.LittleEndian, int32(0))            // thiszone
	binary.Write(f, binary.LittleEndian, uint32(0))           // sigfigs
	binary.Write(f, binary.LittleEndian, uint32(65535))       // snaplen
	binary.Write(f, binary.LittleEndian, uint32(1))           // network (Ethernet)

	// Write 3 packets: a simple Ethernet + IPv4 + TCP SYN-like frame
	for i := 0; i < 3; i++ {
		pkt := makeEthernetIPv4TCP(i)
		// Packet header (16 bytes)
		binary.Write(f, binary.LittleEndian, uint32(1711300000+uint32(i))) // ts_sec
		binary.Write(f, binary.LittleEndian, uint32(i*1000))               // ts_usec
		binary.Write(f, binary.LittleEndian, uint32(len(pkt)))             // incl_len
		binary.Write(f, binary.LittleEndian, uint32(len(pkt)))             // orig_len
		f.Write(pkt)
	}
}

func makeEthernetIPv4TCP(seq int) []byte {
	pkt := make([]byte, 54) // 14 Ethernet + 20 IPv4 + 20 TCP

	// Ethernet header
	// dst MAC: ff:ff:ff:ff:ff:ff
	for i := 0; i < 6; i++ {
		pkt[i] = 0xff
	}
	// src MAC: 00:11:22:33:44:55
	pkt[6] = 0x00
	pkt[7] = 0x11
	pkt[8] = 0x22
	pkt[9] = 0x33
	pkt[10] = 0x44
	pkt[11] = 0x55
	// EtherType: IPv4
	pkt[12] = 0x08
	pkt[13] = 0x00

	// IPv4 header (20 bytes)
	pkt[14] = 0x45       // version + IHL
	pkt[15] = 0x00       // DSCP
	pkt[16] = 0x00       // Total length high
	pkt[17] = 40         // Total length low (20 IPv4 + 20 TCP)
	pkt[18] = 0x00       // ID
	pkt[19] = byte(seq)  // ID low
	pkt[20] = 0x40       // Flags (Don't Fragment)
	pkt[21] = 0x00       // Fragment offset
	pkt[22] = 64         // TTL
	pkt[23] = 6          // Protocol: TCP
	pkt[24] = 0x00       // Checksum (simplified)
	pkt[25] = 0x00       // Checksum
	pkt[26] = 192        // Src IP: 192.168.1.100
	pkt[27] = 168
	pkt[28] = 1
	pkt[29] = 100
	pkt[30] = 93         // Dst IP: 93.184.216.34 (example.com)
	pkt[31] = 184
	pkt[32] = 216
	pkt[33] = 34

	// TCP header (20 bytes)
	pkt[34] = 0xC0                  // Src port high (49152+seq)
	pkt[35] = byte(seq)             // Src port low
	pkt[36] = 0x00                  // Dst port high (80)
	pkt[37] = 80                    // Dst port low
	pkt[38] = 0x00                  // Seq num
	pkt[39] = 0x00
	pkt[40] = 0x00
	pkt[41] = byte(seq)
	pkt[42] = 0x00                  // Ack num
	pkt[43] = 0x00
	pkt[44] = 0x00
	pkt[45] = 0x00
	pkt[46] = 0x50                  // Data offset (5 words)
	pkt[47] = 0x02                  // Flags: SYN
	pkt[48] = 0xFF                  // Window
	pkt[49] = 0xFF
	pkt[50] = 0x00                  // Checksum
	pkt[51] = 0x00
	pkt[52] = 0x00                  // Urgent pointer
	pkt[53] = 0x00

	return pkt
}
