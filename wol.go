package gowol

import (
	"fmt"
	"log"
	"net"
)

// SendMagicPacket sends a Wake-on-LAN (WoL) magic packet to the target device.
func SendMagicPacket(macAddress string, ipAddress string) error {
	if ipAddress == "" {
		ipAddress = "255.255.255.255"
	}

	// Parse the MAC address
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		return fmt.Errorf("error parsing MAC address: %v", err)
	}

	// Create the magic packet
	magicPacket := make([]byte, 102)
	for i := 0; i < 6; i++ {
		magicPacket[i] = 0xff
	}
	for i := 1; i <= 16; i++ {
		copy(magicPacket[i*6:], mac)
	}

	// Create a UDP connection
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ipAddress, 9))
	if err != nil {
		return fmt.Errorf("error dialing UDP: %v", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	// Send the magic packet
	_, err = conn.Write(magicPacket)
	if err != nil {
		return fmt.Errorf("error sending magic packet: %v", err)
	}
	return nil
}
