package radio

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const xcatTimeout = 3 * time.Second

const digitalDialToleranceHz = 2000

var digitalDialFrequencies = []struct {
	mode        string
	frequencies []float64
}{
	{mode: "FT8", frequencies: []float64{
		1840000, 3573000, 7074000, 10136000, 14074000, 18100000,
		21074000, 24915000, 28074000, 50313000, 144174000, 432174000,
	}},
	{mode: "FT4", frequencies: []float64{
		3575000, 7047500, 10140000, 14080000, 18104000,
		21140000, 24919000, 28180000, 50318000,
	}},
}

// XCATClient implements the small rigctld subset exposed by xCAT. Unlike the
// generic Hamlib client it deliberately does not probe SAT mode, split state,
// or RF power: xCAT 2.0 answers those probes with RPRT -4.
type XCATClient struct {
	host    string
	port    string
	timeout time.Duration
}

func NewXCAT(host, port string) *XCATClient {
	if strings.TrimSpace(host) == "" {
		host = "127.0.0.1"
	}
	if strings.TrimSpace(port) == "" {
		port = "4532"
	}
	return &XCATClient{host: host, port: port, timeout: xcatTimeout}
}

func (c *XCATClient) connect() (net.Conn, *bufio.Reader, error) {
	address := net.JoinHostPort(c.host, c.port)
	conn, err := net.DialTimeout("tcp", address, c.timeout)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to xCAT RigCtlD at %s: %w", address, err)
	}
	if err := conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		conn.Close()
		return nil, nil, err
	}
	return conn, bufio.NewReader(conn), nil
}

func readXCATLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if err != nil && line == "" {
		return "", err
	}
	if strings.HasPrefix(line, "RPRT") {
		return "", fmt.Errorf("xCAT RigCtlD error: %s", line)
	}
	if line == "" {
		return "", fmt.Errorf("empty response from xCAT RigCtlD")
	}
	return line, nil
}

func sendXCATSet(conn net.Conn, reader *bufio.Reader, command string) error {
	if _, err := fmt.Fprint(conn, command); err != nil {
		return err
	}
	line, err := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if err != nil && line == "" {
		return err
	}
	if line != "RPRT 0" {
		return fmt.Errorf("xCAT RigCtlD rejected %q: %s", strings.TrimSpace(command), line)
	}
	return nil
}

func (c *XCATClient) GetStatus() (RigStatus, error) {
	var status RigStatus
	conn, reader, err := c.connect()
	if err != nil {
		return status, err
	}
	defer conn.Close()

	if _, err := fmt.Fprint(conn, "f\n"); err != nil {
		return status, err
	}
	freq, err := readXCATLine(reader)
	if err != nil {
		return status, err
	}
	status.FreqA, err = strconv.ParseFloat(freq, 64)
	if err != nil || status.FreqA <= 0 {
		return RigStatus{}, fmt.Errorf("invalid frequency from xCAT RigCtlD: %q", freq)
	}

	if _, err := fmt.Fprint(conn, "m\n"); err != nil {
		return RigStatus{}, err
	}
	status.Mode, err = readXCATLine(reader)
	if err != nil {
		return RigStatus{}, err
	}
	// rigctld mode replies contain a second passband line. Consume it so the
	// protocol remains aligned even though Wavelog does not need the value.
	if _, err := readXCATLine(reader); err != nil {
		return RigStatus{}, err
	}
	status.Mode = normalizeXCATMode(status.Mode, status.FreqA)
	return status, nil
}

// normalizeXCATMode translates xCAT/Hamlib radio modes into values Wavelog can
// select. A radio can only report its data-sideband mode (for example PKTUSB),
// not the application protocol. For the conventional FT8 and FT4 dial
// frequencies, use the protocol name; otherwise retain the underlying sideband.
func normalizeXCATMode(mode string, frequencyHz float64) string {
	mode = strings.ToUpper(strings.TrimSpace(mode))
	switch mode {
	case "PKTUSB", "DIGU", "DATA-U", "USB-D":
		if digitalMode := digitalModeForFrequency(frequencyHz); digitalMode != "" {
			return digitalMode
		}
		return "USB"
	case "PKTLSB", "DIGL", "DATA-L", "LSB-D":
		if digitalMode := digitalModeForFrequency(frequencyHz); digitalMode != "" {
			return digitalMode
		}
		return "LSB"
	case "PKTFM", "FM-D":
		return "FM"
	case "CWR", "CW-R":
		return "CW"
	case "RTTYR", "RTTY-R":
		return "RTTY"
	default:
		return mode
	}
}

func digitalModeForFrequency(frequencyHz float64) string {
	for _, candidate := range digitalDialFrequencies {
		for _, dialFrequency := range candidate.frequencies {
			if frequencyHz >= dialFrequency-digitalDialToleranceHz && frequencyHz <= dialFrequency+digitalDialToleranceHz {
				return candidate.mode
			}
		}
	}
	return ""
}

func (c *XCATClient) SetFreqMode(hz int64, mode string) error {
	if hz <= 0 {
		return fmt.Errorf("frequency must be greater than zero")
	}
	mode = strings.ToUpper(strings.TrimSpace(mode))
	for _, r := range mode {
		if !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') && r != '-' {
			return fmt.Errorf("invalid xCAT mode %q", mode)
		}
	}

	conn, reader, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := sendXCATSet(conn, reader, fmt.Sprintf("F %d\n", hz)); err != nil {
		return err
	}
	if mode != "" {
		return sendXCATSet(conn, reader, fmt.Sprintf("M %s -1\n", mode))
	}
	return nil
}

func (c *XCATClient) SetTxFreq(int64) error {
	return fmt.Errorf("xCAT split-VFO control is not supported")
}

func (c *XCATClient) GetModes() ([]string, error) {
	return []string{
		"LSB", "USB", "CW", "CWR", "RTTY", "RTTYR", "AM", "FM", "WFM",
		"PKTLSB", "PKTUSB", "PKTFM", "SAM",
	}, nil
}
