//go:build integration

package radio

import "testing"

// TestLiveXCAT performs read-only validation against the local xCAT RigCtlD
// endpoint used by the packaged macOS app. It never sends a tune or PTT command.
func TestLiveXCAT(t *testing.T) {
	status, err := NewXCAT("127.0.0.1", "4532").GetStatus()
	if err != nil {
		t.Fatal(err)
	}
	if status.FreqA <= 0 || status.Mode == "" {
		t.Fatalf("incomplete live xCAT status: %+v", status)
	}
	t.Logf("xCAT live status: %.0f Hz %s", status.FreqA, status.Mode)
}
