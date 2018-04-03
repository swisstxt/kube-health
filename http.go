package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	errHttpParse      = Error("parsing HTTP URL")
	errHttpRequest    = Error("requesting HTTP")
	errHttpRead       = Error("HTTP read")
	errHttpStatus     = Error("HTTP status code")
	errHttpContent    = Error("HTTP response content")
	errHttpBundleLoad = Error("HTTP CA bundle load")
)

func ProcessHttp(rawurl string, timeout time.Duration, statusCode int, notCode bool, substring string, notMatch bool, caBundle string) (string, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	if caBundle != "" {
		var err error
		var bundle *os.File
		var pem []byte
		var certs *x509.CertPool

		bundle, err = os.Open(caBundle)
		if err == nil {
			pem, err = ioutil.ReadAll(bundle)
		}
		if err == nil {
			certs = x509.NewCertPool()
		}
		if !certs.AppendCertsFromPEM(pem) {
			err = errors.New("invalid CA bundle")
		}

		if err != nil {
			message := fmt.Sprintf("Error loading certificates: %s", err.Error())
			return message, errHttpBundleLoad
		}

		// from http.DefaultTransport
		client.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				RootCAs: certs,
			},
		}
	}

	request, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		message := fmt.Sprintf("Invalid URL %s: %s", rawurl, err.Error())
		return message, errHttpParse
	}
	response, err := client.Do(request)
	if err != nil {
		message := fmt.Sprintf("Error connecting to %s: %s", request.URL.Host, err.Error())
		return message, errHttpRequest
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		message := fmt.Sprintf("Error reading response: %s", err.Error())
		return message, errHttpRead
	}
	response.Body.Close()

	critical := false

	if statusCode > 0 {
		if notCode {
			critical = response.StatusCode == statusCode
		} else {
			critical = response.StatusCode != statusCode
		}
	}
	if critical {
		message := fmt.Sprintf("Incorrect status code %d, response was: %s", response.StatusCode, body)
		return message, errHttpStatus
	}

	if substring != "" {
		if notMatch {
			critical = strings.Contains(string(body), substring)
		} else {
			critical = !strings.Contains(string(body), substring)
		}
	}
	if critical {
		message := fmt.Sprintf("Incorrect response content, was: %s", body)
		return message, errHttpContent
	}

	message := fmt.Sprintf("Response ok: %s", body)
	return message, nil
}
