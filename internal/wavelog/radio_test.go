package wavelog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"waveloggate/internal/config"
)

func TestUpdateRadioStatusPostsXCATFrequencyAndMode(t *testing.T) {
	received := make(chan radioPayload, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/radio" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("User-Agent"); got != "Flex-WavelogGate/test" {
			t.Errorf("User-Agent = %q", got)
		}
		var payload radioPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		received <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	profile := config.Profile{
		WavelogURL:       server.URL,
		WavelogKey:       "test-key",
		WavelogRadioname: "Maestro",
	}
	client := New(&profile, "test")
	if err := client.UpdateRadioStatus(RadioData{Frequency: 7139200, Mode: "LSB"}); err != nil {
		t.Fatal(err)
	}
	payload := <-received
	if payload.Frequency != 7139200 || payload.Mode != "LSB" || payload.Radio != "Maestro" {
		t.Fatalf("unexpected radio payload: %+v", payload)
	}
	if payload.Key != "test-key" {
		t.Fatalf("v1 API key missing from radio payload")
	}
}
