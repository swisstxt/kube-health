package health

import (
	"os"
	"encoding/json"
)

const (
	defaultListenAddress string = "0.0.0.0:8080"
)

type CheckConfiguration struct {
	Type string `json:"type"`
	Url string `json:"url"`
	Timeout int `json:"timeout"`
	Ping struct {
		Count int `json:"count"`
		Warning int `json:"warning"`
		Error int `json:"error"`
	} `json:"ping"`
	Http struct {
		Status int `json:"status"`
		InvertStatus bool `json:"invert_status"`
		Contains string `json:"contains"`
		InvertMatch bool `json:"invert_match"`
		CaBundle string `json:"ca_certificates"`
	} `json:"http"`
}

type Configuration struct {
	Listen string `json:"listen"`
	Checks []CheckConfiguration `json:"checks"`
}

func DefaultConfiguration() *Configuration {
	return &Configuration{
		Listen: defaultListenAddress,
	}
}

func LoadConfiguration(filename string) (*Configuration, error) {
	config := DefaultConfiguration()
	
	fd, err := os.Open(filename)
	if err == nil {
		decoder := json.NewDecoder(fd)
		err = decoder.Decode(&config)
		fd.Close()
	}
	
	return config, err
}