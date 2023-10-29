package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddress_UnmarshalText(t *testing.T) {
	host := "localhost"
	port := 8080
	tests := []struct {
		name    string
		env     string
		wantErr bool
		want    *Address
	}{
		{
			"should parse correct",
			"localhost:8080",
			false,
			&Address{
				host,
				port,
			},
		},
		{
			"should return nil",
			"localhost:abv",
			true,
			&Address{},
		},
		{
			"should return nil",
			"abv",
			true,
			&Address{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Address{}
			err := a.UnmarshalText([]byte(tt.env))
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want, a)
		})
	}
}

func TestAddress_Config(t *testing.T) {
	host := "localhost"
	port := 8080
	poll := time.Duration(2) * time.Second
	report := time.Duration(10) * time.Second
	tests := []struct {
		name    string
		env     map[string]string
		want    Config
		wantErr bool
	}{
		{
			"should parse correct",
			map[string]string{
				"ADDRESS":         "localhost:8080",
				"POLL_INTERVAL":   "2",
				"REPORT_INTERVAL": "10",
			},
			Config{
				Address: Address{
					Host: host,
					Port: port,
				},
				PollInterval:   poll,
				ReportInterval: report,
			},
			false,
		},
		{
			"should parse correct 2",
			map[string]string{
				"ADDRESS":         "localhost:8080",
				"POLL_INTERVAL":   "2",
				"REPORT_INTERVAL": "10",
			},
			Config{
				Address: Address{
					Host: host,
					Port: port,
				},
				PollInterval:   poll,
				ReportInterval: report,
			},
			false,
		},
		{
			"should return empty",
			map[string]string{},
			Config{
				PollInterval:   poll,
				ReportInterval: report,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			cfg := NewConfigFromENV()
			assert.Equal(t, tt.want, cfg)
			os.Clearenv()
		})
	}
}
