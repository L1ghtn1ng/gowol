package gowol

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestSendMagicPacket(test *testing.T) {
	test.Run("success", func(t *testing.T) {
		conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
		if err != nil {
			t.Fatalf("listen UDP: %v", err)
		}
		defer conn.Close()

		if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			t.Fatalf("set read deadline: %v", err)
		}

		errCh := make(chan error, 1)
		go func() {
			addr := conn.LocalAddr().(*net.UDPAddr)
			// Send to the listener's random port so the test never needs privileged UDP port 9.
			errCh <- sendMagicPacket("00:11:22:33:44:55", addr.IP.String(), addr.Port)
		}()

		packet := make([]byte, 102)
		n, _, err := conn.ReadFromUDP(packet)
		if err != nil {
			t.Fatalf("read UDP packet: %v", err)
		}

		if err := <-errCh; err != nil {
			t.Fatalf("SendMagicPacket returned error: %v", err)
		}
		assertMagicPacket(t, packet[:n], []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
	})

	test.Run("builds expected packet", func(t *testing.T) {
		packet, err := newMagicPacket("00:11:22:33:44:55")
		if err != nil {
			t.Fatalf("newMagicPacket returned error: %v", err)
		}
		assertMagicPacket(t, packet, []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
	})

	test.Run("invalid MAC address", func(t *testing.T) {
		err := SendMagicPacket("invalid", "255.255.255.255")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRemoteUDPAddr(t *testing.T) {
	t.Run("defaults empty IP", func(t *testing.T) {
		addr, err := remoteUDPAddr("", 9)
		if err != nil {
			t.Fatalf("remoteUDPAddr returned error: %v", err)
		}

		if got := addr.String(); got != "255.255.255.255:9" {
			t.Fatalf("remote address = %q, want %q", got, "255.255.255.255:9")
		}
	})

	t.Run("rejects invalid IP", func(t *testing.T) {
		if _, err := remoteUDPAddr("invalid", 9); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("rejects invalid port", func(t *testing.T) {
		if _, err := remoteUDPAddr("127.0.0.1", 65536); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func assertMagicPacket(t *testing.T, packet []byte, mac []byte) {
	t.Helper()

	if len(packet) != 102 {
		t.Fatalf("packet length = %d, want 102", len(packet))
	}

	if !bytes.Equal(packet[:6], bytes.Repeat([]byte{0xff}, 6)) {
		t.Fatalf("packet header = %x, want six 0xff bytes", packet[:6])
	}

	for i := 1; i <= 16; i++ {
		offset := i * 6
		if !bytes.Equal(packet[offset:offset+6], mac) {
			t.Fatalf("packet MAC copy %d = %x, want %x", i, packet[offset:offset+6], mac)
		}
	}
}
