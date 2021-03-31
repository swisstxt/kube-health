package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	errResolveDns     = Error("resolving domain name")
	errNoDnsNameMatch = Error("no name match")
)

func ProcessDns(name string, expected string, timeout time.Duration) (string, error) {
	addrs, err := net.LookupHost(name)
	if err != nil {
		message := fmt.Sprintf("Cannot resolve: %s", err.Error())
		return message, errResolveDns
	}

	var match bool
	if expected == "" {
		message := fmt.Sprintf("DNS name resolved to [%s]", strings.Join(addrs, ","))
		return message, nil
	} else {
		for _, addr := range addrs {
			if addr == expected {
				match = true
			}
		}
	}

	if match {
		message := fmt.Sprintf("Address %s found in [%s]", expected, strings.Join(addrs, ","))
		return message, nil
	} else {
		message := fmt.Sprintf("Address %s not found in [%s]", expected, strings.Join(addrs, ","))
		return message, errNoDnsNameMatch
	}
}
