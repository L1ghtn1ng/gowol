package gowol

import (
	"context"
	"fmt"
	"net"
	"net/netip"
)

const defaultBroadcastIP = "255.255.255.255"

// SendMagicPacket sends a Wake-on-LAN (WoL) magic packet to the target device.
func SendMagicPacket(macAddress string, ipAddress string) error {
	return sendMagicPacket(macAddress, ipAddress, 9)
}

func sendMagicPacket(macAddress string, ipAddress string, port int) error {
	magicPacket, err := newMagicPacket(macAddress)
	if err != nil {
		return err
	}

	remoteAddr, err := remoteUDPAddr(ipAddress, port)
	if err != nil {
		return err
	}

	var dialer net.Dialer
	// DialUDP uses typed netip addresses in Go 1.26, avoiding host:port string formatting.
	conn, err := dialer.DialUDP(context.Background(), "udp", netip.AddrPort{}, remoteAddr)
	if err != nil {
		return fmt.Errorf("error dialing UDP: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(magicPacket); err != nil {
		return fmt.Errorf("error sending magic packet: %w", err)
	}
	return nil
}

func remoteUDPAddr(ipAddress string, port int) (netip.AddrPort, error) {
	if ipAddress == "" {
		ipAddress = defaultBroadcastIP
	}
	if port < 0 || port > 65535 {
		return netip.AddrPort{}, fmt.Errorf("invalid UDP port: %d", port)
	}

	remoteIP, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return netip.AddrPort{}, fmt.Errorf("error parsing IP address: %w", err)
	}
	return netip.AddrPortFrom(remoteIP, uint16(port)), nil
}

func newMagicPacket(macAddress string) ([]byte, error) {
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		return nil, fmt.Errorf("error parsing MAC address: %w", err)
	}

	// A WoL magic packet is six 0xff bytes followed by the target MAC repeated 16 times.
	magicPacket := make([]byte, 102)
	for i := range 6 {
		magicPacket[i] = 0xff
	}
	for i := 1; i <= 16; i++ {
		copy(magicPacket[i*6:], mac)
	}
	return magicPacket, nil
}
