package radio

import (
	"bufio"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"

	"waveloggate/internal/config"
)

func xcatTestServer(t *testing.T, handler func(net.Conn) error) (string, string, <-chan error) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	done := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()
		done <- handler(conn)
	}()
	host, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	return host, port, done
}

func TestXCATGetStatusUsesSupportedSubset(t *testing.T) {
	var commands []string
	host, port, done := xcatTestServer(t, func(conn net.Conn) error {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			command := scanner.Text()
			commands = append(commands, command)
			switch command {
			case "f":
				_, _ = fmt.Fprint(conn, "14225000\n")
			case "m":
				_, _ = fmt.Fprint(conn, "USB\n2400\n")
				return nil
			default:
				return fmt.Errorf("unexpected command %q", command)
			}
		}
		return scanner.Err()
	})

	status, err := NewXCAT(host, port).GetStatus()
	if err != nil {
		t.Fatal(err)
	}
	if status.FreqA != 14225000 || status.Mode != "USB" {
		t.Fatalf("unexpected status: %+v", status)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if want := []string{"f", "m"}; !reflect.DeepEqual(commands, want) {
		t.Fatalf("commands = %v, want %v", commands, want)
	}
}

func TestXCATSetFreqModeChecksReplies(t *testing.T) {
	var commands []string
	host, port, done := xcatTestServer(t, func(conn net.Conn) error {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			commands = append(commands, scanner.Text())
			_, _ = fmt.Fprint(conn, "RPRT 0\n")
			if len(commands) == 2 {
				return nil
			}
		}
		return scanner.Err()
	})

	if err := NewXCAT(host, port).SetFreqMode(7139200, "lsb"); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	want := []string{"F 7139200", "M LSB -1"}
	if !reflect.DeepEqual(commands, want) {
		t.Fatalf("commands = %v, want %v", commands, want)
	}
}

func TestXCATRejectsUnsafeMode(t *testing.T) {
	err := NewXCAT("127.0.0.1", "4532").SetFreqMode(7139200, "USB\nF 1")
	if err == nil || !strings.Contains(err.Error(), "invalid xCAT mode") {
		t.Fatalf("expected invalid mode error, got %v", err)
	}
}

func TestNormalizeXCATMode(t *testing.T) {
	tests := []struct {
		name string
		mode string
		hz   float64
		want string
	}{
		{name: "voice lower sideband", mode: "LSB", hz: 7128000, want: "LSB"},
		{name: "voice upper sideband", mode: "USB", hz: 14225000, want: "USB"},
		{name: "cw reverse", mode: "CWR", hz: 7025000, want: "CW"},
		{name: "ft8 data mode", mode: "PKTUSB", hz: 7074000, want: "FT8"},
		{name: "ft8 tuning tolerance", mode: "DATA-U", hz: 14075200, want: "FT8"},
		{name: "ft4 data mode", mode: "PKTUSB", hz: 7047500, want: "FT4"},
		{name: "generic data mode", mode: "PKTUSB", hz: 7100000, want: "USB"},
		{name: "digital lower sideband fallback", mode: "DIGL", hz: 7100000, want: "LSB"},
		{name: "packet fm", mode: "PKTFM", hz: 145500000, want: "FM"},
		{name: "rtty reverse", mode: "RTTYR", hz: 7045000, want: "RTTY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeXCATMode(tt.mode, tt.hz); got != tt.want {
				t.Fatalf("normalizeXCATMode(%q, %.0f) = %q, want %q", tt.mode, tt.hz, got, tt.want)
			}
		})
	}
}

func TestBuildClientOnlyUsesXCAT(t *testing.T) {
	client := buildClient(&config.Profile{
		XCATEna:   true,
		XCATHost:  "127.0.0.1",
		XCATPort:  "4532",
		FlrigEna:  true,
		HamlibEna: true,
	})
	if _, ok := client.(*XCATClient); !ok {
		t.Fatalf("buildClient() = %T, want *XCATClient", client)
	}
}
