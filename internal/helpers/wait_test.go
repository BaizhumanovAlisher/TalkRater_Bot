package helpers

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWaitSignal(t *testing.T) {
	// Create a channel to receive the signal
	signalReceived := make(chan os.Signal, 1)

	// Start waiting for signals in a goroutine
	go func() {
		signalReceived <- WaitSignal()
	}()

	// Simulate sending a SIGTERM signal after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	// Verify that the received signal is SIGTERM
	select {
	case sig := <-signalReceived:
		if sig != syscall.SIGTERM {
			t.Errorf("Received unexpected signal. Expected: %v, Got: %v", syscall.SIGTERM, sig)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for signal")
	}
}
