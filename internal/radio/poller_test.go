package radio

import (
	"testing"

	"waveloggate/internal/config"
)

type fixedRadioClient struct {
	status RigStatus
}

func (c fixedRadioClient) GetStatus() (RigStatus, error) { return c.status, nil }
func (fixedRadioClient) SetFreqMode(int64, string) error { return nil }
func (fixedRadioClient) SetTxFreq(int64) error           { return nil }
func (fixedRadioClient) GetModes() ([]string, error)     { return nil, nil }

func TestPollPublishesCurrentStatusEvenWhenUnchanged(t *testing.T) {
	profile := &config.Profile{XCATEna: true}
	callbackCount := 0
	poller := &Poller{
		client: fixedRadioClient{status: RigStatus{FreqA: 7139200, Mode: "LSB"}},
		cfg:    profile,
		onStatus: func(RigStatus) {
			callbackCount++
		},
	}

	poller.poll()
	poller.poll()
	if callbackCount != 2 {
		t.Fatalf("status callback count = %d, want 2", callbackCount)
	}
}
