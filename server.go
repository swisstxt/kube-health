package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	HttpStatusOk       = http.StatusOK
	HttpStatusWarning  = 280
	HttpStatusCritical = 290
	HttpStatusUnknown  = 270
)

type CheckResult struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Check   CheckConfiguration `json:"check"`
}

type Result struct {
	Status  string        `json:"status"`
	Results []CheckResult `json:"results,omitempty"`
}

type Server struct {
	Config *Configuration
}

func sanitisePath(line string) string {
	return strings.ToLower(path.Clean(line))
}

func processError(result *CheckResult, message string, err error, loglevel int, critical, warning *bool) {
	result.Message = message
	if err != nil {
		if err.(ErrorWithLevel).IsCritical() {
			result.Status = "critical"
			*critical = true
		} else {
			result.Status = "warning"
			*warning = true
		}
	} else {
		result.Status = "ok"
	}
}

func (server *Server) processCheck() ([]byte, int, error) {
	results := make([]CheckResult, len(server.Config.Checks))
	warning := false
	critical := false
	unknown := false

	for i, check := range server.Config.Checks {
		result := CheckResult{
			Check: check,
		}
		switch check.Type {
		case "ping":
			message, err := ProcessPing(check.Url, time.Second*time.Duration(check.Timeout), check.Ping.Count, float64(check.Ping.Warning), float64(check.Ping.Error))
			processError(&result, message, err, server.Config.LogLevel, &critical, &warning)
		case "http":
			message, err := ProcessHttp(check.Url, time.Second*time.Duration(check.Timeout), check.Http.Status, check.Http.InvertStatus, check.Http.Contains, check.Http.InvertMatch, check.Http.CaBundle)
			processError(&result, message, err, server.Config.LogLevel, &critical, &warning)
		case "dns":
			message, err := ProcessDns(check.Url, check.Dns.Expect, time.Second*time.Duration(check.Timeout))
			processError(&result, message, err, server.Config.LogLevel, &critical, &warning)
		default:
			unknown = true
			result.Status = "unknown"
			result.Message = "Invalid check type"
		}
		results[i] = result
	}

	var code int
	result := Result{}
	if critical {
		result.Status = "critical"
		code = HttpStatusCritical
		if server.Config.LogLevel >= 1 {
			result.Results = results
		}
	} else if warning {
		result.Status = "warning"
		code = HttpStatusWarning
		if server.Config.LogLevel >= 2 {
			result.Results = results
		}
	} else if unknown {
		result.Status = "unknown"
		code = HttpStatusUnknown
		if server.Config.LogLevel >= 2 {
			result.Results = results
		}
	} else {
		result.Status = "ok"
		code = HttpStatusOk
		if server.Config.LogLevel >= 3 {
			result.Results = results
		}
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		return nil, 0, err
	}
	return encoded, code, nil
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch sanitisePath(request.URL.Path) {
	case "/":
		content, status, err := server.processCheck()
		if err == nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(status)
			writer.Write(content)
		} else {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("{\"status\":\"error\",\"error\":\"internal server error\"}"))
		}
	case "/healthz":
		fallthrough
	case "/live":
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("{\"status\":\"ok\"}"))
	default:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("{\"status\":\"error\",\"error\":\"invalid path\"}"))
	}
}
