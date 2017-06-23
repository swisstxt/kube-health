package health

import (
	"fmt"
	"time"
	"strings"
	"io/ioutil"
	"net/http"
)

const (
	errHttpParse = Error("parsing HTTP URL")
	errHttpRequest = Error("requesting HTTP")
	errHttpRead = Error("HTTP read")
	errHttpStatus = Error("HTTP status code")
	errHttpContent = Error("HTTP response content")
)

func ProcessHttp(rawurl string, timeout time.Duration, statusCode int, notCode bool, substring string, notMatch bool) (string, error) {
	client := &http.Client{
		Timeout: timeout,
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
			critical = !strings.Contains(string(body), substring)
		} else {
			critical = strings.Contains(string(body), substring)
		}
	}
	if critical {
		message := fmt.Sprintf("Incorrect response content, was: %s", body)
		return message, errHttpContent
	}
	
	message := fmt.Sprintf("Response ok: %s", body)
	return message, nil
}
