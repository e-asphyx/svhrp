package main

import (
	"encoding/json"
	"io/ioutil"
)

type Route struct {
	Host     string
	CertFile string
	CertPEM  string
	KeyFile  string
	KeyPEM   string
}

type Config struct {
	Listen                string
	HttpsRedirectorListen string
	Routes                map[string]Route
}

func NewConfig(path string) (*Config, error) {
	config := Config{
		Listen:                ":10443",
		HttpsRedirectorListen: ":8080",
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
