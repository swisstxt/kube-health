package health

import (
	"fmt"
	"time"
	"github.com/sparrc/go-ping"
)

const (
	errResolve = Error("resolving address")
	warnPacketLoss = Warning("packet loss warning")
	errPacketLoss = Error("packet loss critical")
)

func ProcessPing(addr string, timeout time.Duration, count int, warnLost float64, errLost float64) (string, error) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		message := fmt.Sprintf("Error resolving %s: %s", addr, err.Error())
		return message, errResolve
	}
	
	pinger.Timeout = timeout
	pinger.Count = count
	pinger.SetPrivileged(true)
	
	// Blocks while the ping is running
	pinger.Run()
	
	stats := pinger.Statistics()
	if stats.PacketLoss >= errLost {
		message := fmt.Sprintf("Critical packet loss: %d%%", int(stats.PacketLoss))
		return message, errPacketLoss
	}
	if stats.PacketLoss >= warnLost {
		message := fmt.Sprintf("Packet loss: %d%%", int(stats.PacketLoss))
		return message, warnPacketLoss
	}
	
	return "No packet loss", nil
}
