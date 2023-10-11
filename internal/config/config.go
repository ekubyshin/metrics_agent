package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/ekubyshin/metrics_agent/internal/utils"
)

var regExp = regexp.MustCompile(`^\w+:[0-9]{2,5}$`)

// Это промежуточная структура, которая принимает int
// тк в тестах интервалы передаются в int, а не в duration
// А мне хочется хранить в конфиге именно duration
// поэтому я читаю int, а потом возвращаю уже нужную конфигу с duration
type config struct {
	Address        *Address `env:"ADDRESS"`
	ReportInterval *int     `env:"REPORT_INTERVAL"`
	PollInterval   *int     `env:"POLL_INTERVAL"`
}

type Config struct {
	Address        Address       `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

type Address struct {
	Host string
	Port int
}

type Builder struct {
	config Config
}

func NewBuilder() Builder {
	return Builder{}
}

func (b Builder) WithAddress(address Address) Builder {
	b.config.Address = address
	return b
}

func (b Builder) WithHost(host string) Builder {
	b.config.Address.Host = host
	return b
}

func (b Builder) WithPort(port int) Builder {
	b.config.Address.Port = port
	return b
}

func (b Builder) WithReportInterval(t int) Builder {
	b.config.ReportInterval = utils.IntToDuration(t)
	return b
}

func (b Builder) WithPollInterval(t int) Builder {
	b.config.PollInterval = utils.IntToDuration(t)
	return b
}

func (b Builder) Build() Config {
	return b.config
}

func (a *Address) ToString() string {
	return fmt.Sprintf("%v:%v", a.Host, a.Port)
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
	*a = Address{Host: host, Port: port}
	return nil
}

func NewConfigFromENV() Config {
	tcfg := &config{}
	builder := NewBuilder()
	if err := env.Parse(tcfg); err != nil {
		return builder.Build()
	}
	if tcfg.Address != nil {
		builder = builder.WithAddress(*tcfg.Address)
	}
	if tcfg.PollInterval != nil && *tcfg.PollInterval > 0 {
		builder = builder.WithPollInterval(*tcfg.PollInterval)
	}
	if tcfg.ReportInterval != nil && *tcfg.ReportInterval > 0 {
		builder = builder.WithReportInterval(*tcfg.ReportInterval)
	}
	return builder.Build()
}

func NewConfigFromFlags() Config {
	endpoint := flag.String("a", "localhost:8080", "endpoint address")
	reportInterval := flag.Int("r", 10, "report interval")
	pollInterval := flag.Int("p", 2, "poll interval")
	flag.Parse()
	builer := NewBuilder()
	address := &Address{}
	_ = address.UnmarshalText([]byte(*endpoint))
	return builer.
		WithAddress(*address).
		WithPollInterval(*pollInterval).
		WithReportInterval(*reportInterval).
		Build()
}

func AutoLoad() Config {
	if _, ok := os.LookupEnv("ADDRESS"); ok {
		return NewConfigFromENV()
	}
	return NewConfigFromFlags()
}
