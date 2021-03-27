package config

import (
	"flag"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/jinzhu/configor"
	"io/ioutil"
	"os"
	"strings"
)

const (
	KeyUser   = "user"
	KeyClient = "client"
)

const (
	ClientIdPrefixLength = 8
)

type EmailPasswordTuple struct {
	Email    string
	Password string
}

type mailConfig struct {
	Domain string `yaml:"domain" env:"MW_MAIL_DOMAIN"`
	SPF    spfConfig
}

type spfConfig struct {
	Check              bool     `env:"MW_MAIL_SPF_CHECK"`
	AuthorizedIncludes []string `yaml:"authorized_includes"`
}

type smtpConfig struct {
	Host     string `env:"MW_SMTP_HOST"`
	Port     uint   `env:"MW_SMTP_PORT"`
	Username string `env:"MW_SMTP_USER"`
	Password string `env:"MW_SMTP_PASS"`
	TLS      bool   `env:"MW_SMTP_TLS"`
}

type webConfig struct {
	ListenV4    string   `yaml:"listen_v4" default:"127.0.0.1:3000" env:"MW_WEB_LISTEN_V4"`
	CorsOrigins []string `yaml:"cors_origins" env:"MW_WEB_CORS_ORIGINS"`
}

type storeConfig struct {
	Path string `default:"data.gob.db" env:"MW_STORE_PATH"`
}

type securityConfig struct {
	Pepper    string               `env:"MW_SECURITY_PEPPER"`
	SeedUsers []EmailPasswordTuple `yaml:"seed_users"`
}

type Config struct {
	Env      string `default:"dev" env:"MW_ENV"`
	Version  string
	Mail     mailConfig
	Web      webConfig
	Smtp     smtpConfig
	Store    storeConfig
	Security securityConfig
}

var cfg *Config

func Get() *Config {
	return cfg
}

func Set(config *Config) {
	cfg = config
}

func Load() *Config {
	config := &Config{}

	flag.Parse()

	if err := configor.New(&configor.Config{}).Load(config, "config.yml"); err != nil {
		logbuch.Fatal("failed to read config: %v", err)
	}

	config.Version = readVersion()

	if config.Web.ListenV4 == "" {
		logbuch.Fatal("config option 'listen4' must be specified")
	}

	if config.Security.SeedUsers == nil || len(config.Security.SeedUsers) == 0 {
		logbuch.Fatal("you need to create at least one initial user via config.yml as there currently is no way to dynamically create users at runtime")
	}

	Set(config)
	return Get()
}

func (c *smtpConfig) ConnStr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) IsDev() bool {
	return c.Env == "dev" || c.Env == "development"
}

func readVersion() string {
	file, err := os.Open("version.txt")
	if err != nil {
		logbuch.Fatal("failed to read version: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		logbuch.Fatal(err.Error())
	}

	return strings.TrimSpace(string(bytes))
}
