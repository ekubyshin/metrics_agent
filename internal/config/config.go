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
	"github.com/ekubyshin/metrics_agent/internal/pointer"
)

var regExp = regexp.MustCompile(`^\w+:[0-9]{2,5}$`)

const (
	defaultPath           = "/tmp/metrics-db.json"
	defaultReportInterval = 10
	defaultPollInterval   = 2
	defaultStoreInterval  = 300
	shouldRestore         = true
)

type Config struct {
	Address         Address `env:"ADDRESS"`
	ReportInterval  int     `env:"REPORT_INTERVAL"`
	PollInterval    int     `env:"POLL_INTERVAL"`
	StoreInterval   *int    `env:"STORE_INTERVAL"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
	Restore         *bool   `env:"RESTORE"`
	Env             string  `env:"Env"`
	DBDSN           *string `env:"DATABASE_DSN"`
}

func (c Config) ReportDuration() time.Duration {
	return time.Duration(c.ReportInterval) * time.Second
}

func (c Config) PollDuration() time.Duration {
	return time.Duration(c.PollInterval) * time.Second
}

func (c Config) StoreDuration() time.Duration {
	if c.StoreInterval == nil {
		return time.Duration(defaultStoreInterval) * time.Second
	}
	return time.Duration(*c.StoreInterval) * time.Second
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
	b.config.ReportInterval = t
	return b
}

func (b Builder) WithPollInterval(t int) Builder {
	b.config.PollInterval = t
	return b
}

func (b Builder) WithStoreInterval(t int) Builder {
	b.config.StoreInterval = pointer.From[int](t)
	return b
}

func (b Builder) WithStoreFilePath(p string) Builder {
	b.config.FileStoragePath = pointer.From[string](p)
	return b
}

func (b Builder) WithRestore(r bool) Builder {
	b.config.Restore = &r
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

func NewServerConfigFromENV() (cfg Config) {
	if err := env.Parse(&cfg); err != nil {
		return
	}
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = pointer.From[int](defaultStoreInterval)
	}
	if cfg.FileStoragePath == nil {
		cfg.FileStoragePath = pointer.From[string](defaultPath)
	}
	if cfg.Restore == nil {
		cfg.Restore = pointer.From[bool](shouldRestore)
	}
	return
}

func NewAgentConfigFromENV() (cfg Config) {
	if err := env.Parse(&cfg); err != nil {
		return
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = defaultPollInterval
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = defaultReportInterval
	}
	return
}

func NewAgentConfigFromFlags() Config {
	endpoint := flag.String("a", "localhost:8080", "endpoint address")
	reportInterval := flag.Int("r", defaultReportInterval, "report interval")
	pollInterval := flag.Int("p", defaultPollInterval, "poll interval")
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

func NewServerConfigFromFlags() Config {
	endpoint := flag.String("a", "localhost:8080", "endpoint address")
	storeInterval := flag.Int("i", defaultStoreInterval, "store interval")
	fileStorage := flag.String("f", defaultPath, "store db file path")
	restore := flag.Bool("r", shouldRestore, "should restore db")
	flag.Parse()
	builer := NewBuilder()
	address := &Address{}
	_ = address.UnmarshalText([]byte(*endpoint))
	return builer.
		WithAddress(*address).
		WithStoreInterval(*storeInterval).
		WithStoreFilePath(*fileStorage).
		WithRestore(*restore).
		Build()
}

func AutoLoadAgent() Config {
	if _, ok := os.LookupEnv("ADDRESS"); ok {
		return NewAgentConfigFromENV()
	}
	return NewAgentConfigFromFlags()
}

func AutoLoadServer() Config {
	if _, ok := os.LookupEnv("ADDRESS"); ok {
		return NewServerConfigFromENV()
	}
	return NewServerConfigFromFlags()
}
