package health

import (
	"fmt"
	"time"
	"strings"
	"regexp"
	"strconv"
	"bytes"
	"os/exec"
)

const (
	errResolve = Error("resolving address")
	warnPacketLoss = Warning("packet loss warning")
	errPacketLoss = Error("packet loss critical")
	errExec = Error("calling ping")
	errParse = Error("parsing output")
	errPing = Error("error returned")
)

var (
	pingStatMatcher = regexp.MustCompile("([0-9]+) packets transmitted, ([0-9]+) received, ([0-9]+)% packet loss, time ([0-9]+)ms")
)

func ProcessPing(addr string, timeout time.Duration, count int, warnLost float64, errLost float64) (string, error) {
	cmd := exec.Command("ping", "-c", strconv.FormatInt(int64(count), 10), "-n", "-q", "-W", strconv.FormatInt(int64(timeout.Seconds()), 10), addr)
	
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	
	
	
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			message := fmt.Sprintf("Execute error: %s", err.Error())
			return message, errExec
		}
		// non-zero return codes are treated as errors by Go, but we'll handle them below
	}
	
	found := false
	//sent := -1
	//received := -1
	loss := 100.0
	rtt := 0
	for _, line := range strings.Split(stdout.String(), "\n") {
		fields := pingStatMatcher.FindStringSubmatch(line)
		if fields != nil {
			//sent, _ = strconv.Atoi(fields[1])
			//received, _ = strconv.Atoi(fields[2])
			loss, _ = strconv.ParseFloat(fields[3], 64)
			rtt, _ = strconv.Atoi(fields[4])
			found = true
		}
	}
	
	if found {
		if loss >= errLost {
			message := fmt.Sprintf("Critical packet loss: %d%%", int(loss))
			return message, errPacketLoss
		}
		if loss >= warnLost {
			message := fmt.Sprintf("Packet loss: %d%%", int(loss))
			return message, warnPacketLoss
		}
		message := fmt.Sprintf("No packet loss, RTT %dms", rtt)
		return message, nil
	}
	
	message := fmt.Sprintf("Ping error: %s", stderr.String())
	return message, errPing
}
