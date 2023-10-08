package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/ekubyshin/metrics_agent/internal/utils"
)

var regExp = regexp.MustCompile(`^\w+:[0-9]{2,5}$`)

type config struct {
	Address        *Address `env:"ADDRESS"`
	ReportInterval *int     `env:"REPORT_INTERVAL"`
	PollInterval   *int     `env:"POLL_INTERVAL"`
}

type Config struct {
	Address        *Address       `env:"ADDRESS"`
	ReportInterval *time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   *time.Duration `env:"POLL_INTERVAL"`
}

type Address struct {
	Host *string
	Port *int
}

func (a *Address) ToString() string {
	if a == nil || a.Host == nil || a.Port == nil {
		return ""
	}
	return fmt.Sprintf("%v:%v", *a.Host, *a.Port)
}

func (a *Address) UnmarshalText(text []byte) error {
	if !regExp.Match(text) {
		return errors.New("invalid endpoint format")
	}
	parts := strings.Split(string(text), ":")
	if len(parts) != 2 {
		return errors.New("invalid endpoint format")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	host := parts[0]
	*a = Address{Host: &host, Port: &port}
	return nil
}

func NewConfig() *Config {
	tcfg := &config{}
	cfg := &Config{}
	if err := env.Parse(tcfg); err != nil {
		return nil
	}
	cfg.Address = tcfg.Address
	if tcfg.PollInterval != nil && *tcfg.PollInterval > 0 {
		t := utils.IntToDuration(*tcfg.PollInterval)
		cfg.PollInterval = &t
	}
	if tcfg.ReportInterval != nil && *tcfg.ReportInterval > 0 {
		t := utils.IntToDuration(*tcfg.ReportInterval)
		cfg.ReportInterval = &t
	}
	return cfg
}
