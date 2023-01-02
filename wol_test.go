package gowol

import (
	"testing"
)

func TestSendMagicPacket(test *testing.T) {
	test.Run("success", func(t *testing.T) {
		// This test should pass if the SendMagicPacket function returns no error
		err := SendMagicPacket("00:11:22:33:44:55", "255.255.255.255")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	test.Run("invalid MAC address", func(t *testing.T) {
		// This test should pass if the SendMagicPacket function returns an error
		// when given an invalid MAC address
		err := SendMagicPacket("invalid", "255.255.255.255")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
