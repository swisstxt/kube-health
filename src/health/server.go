package health

import (
	"time"
	"path"
	"strings"
	"net/http"
	"encoding/json"
)

const (
	HttpStatusOk = http.StatusOK
	HttpStatusWarning = 280
	HttpStatusCritical = 290
)

type CheckResult struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Check CheckConfiguration `json:"check"`
}

type Result struct {
	Status string `json:"status"`
	Results []CheckResult `json:"results"`
}

type Server struct {
	Config *Configuration
}

func SanitisePath(line string) string {
	return strings.ToLower(path.Clean(line))
}

func (server *Server) ProcessCheck() ([]byte, int, error) {
	results := make([]CheckResult, len(server.Config.Checks))
	warning := false
	critical := false
	
	for i, check := range server.Config.Checks {
		result := CheckResult{
			Check: check,
		}
		switch check.Type {
			case "ping":
				message, err := ProcessPing(check.Url, time.Second * time.Duration(check.Timeout), check.Ping.Count, float64(check.Ping.Warning), float64(check.Ping.Error))
				result.Message = message
				if err != nil {
					if err.(ErrorWithLevel).IsCritical() {
						result.Status = "critical"
						critical = true
					} else {
						result.Status = "warning"
						warning = true
					}
				} else {
					result.Status = "ok"
				}
			case "http":
				message, err := ProcessHttp(check.Url, time.Second * time.Duration(check.Timeout), check.Http.Status, check.Http.InvertStatus, check.Http.Contains, check.Http.InvertMatch)
				result.Message = message
				if err != nil {
					if err.(ErrorWithLevel).IsCritical() {
						result.Status = "critical"
						critical = true
					} else {
						result.Status = "warning"
						warning = true
					}
				} else {
					result.Status = "ok"
				}
			default:
				result.Status = "unknown"
				result.Message = "Invalid check type"
		}
		results[i] = result
	}
	
	var code int
	result := Result{
		Results: results,
	}
	if critical {
		result.Status = "critical"
		code = HttpStatusCritical
	} else if warning {
		result.Status = "warning"
		code = HttpStatusWarning
	} else {
		result.Status = "ok"
		code = HttpStatusOk
	}
	
	encoded, err := json.Marshal(result)
	if err != nil {
		return nil, 0, err
	}
	return encoded, code, nil
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch SanitisePath(request.URL.Path) {
		case "/":
			content, status, err := server.ProcessCheck()
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
