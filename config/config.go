package config

import (
	"fmt"
	"go-boilerplate/infrastructure/database"
	"go-boilerplate/infrastructure/jwt"
	"go-boilerplate/infrastructure/redis"
	"go-boilerplate/pkg/ratelimiter"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Credential     CredentialConfig     `json:"credential"`
	App            AppConfig            `json:"app"`
	Infrastructure InfrastructureConfig `json:"infrastructure"`
	Feature        FeatureConfig        `json:"feature"`
}

type AppConfig struct {
	Env     string   `json:"env"`
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Port    int      `json:"port"`
	Host    []string `json:"host"`
}

type CredentialConfig struct {
	Auth AuthCredentialConfig `json:"auth"`
	Jwt  jwt.Config           `json:"jwt"`
	SMTP SMPTConfig           `json:"smtp"`
}

type AuthCredentialConfig struct {
	SecretKey string `json:"secretKey"`
	IV        string `json:"iv"`
}

type FeatureConfig struct {
	GracefullyShutdown time.Duration     `json:"gracefullShutdown"`
	RateLimiter        RateLimiterConfig `json:"rateLimiter"`
}

type RateLimiterConfig struct {
	HTTPLimiter ratelimiter.Config `json:"httpLimiter"`
}

type InfrastructureConfig struct {
	Database database.Config `json:"database"`
	Redis    redis.Config    `json:"redis"`
}

type SysLogConfig struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Stdout   bool   `json:"stdout"`
}

type TdrLogConfig struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Stdout   bool   `json:"stdout"`
}

type SMPTConfig struct {
	MailFrom           string `json:"mailFrom"`
	MailAdminRecipient string `json:"mailAdminRecipient"`
	MailSMTP           string `json:"mailSMTP"`
	MailPort           string `json:"mailPort"`
	MailUsername       string `json:"mailUsername"`
	MailPassword       string `json:"mailPassword"`
}

func New() (cfg *Config) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./resources")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("fatal error config file: config file not found")
		}
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	jwtExpStr := viper.GetString("credential.jwt.jwtExpiredTime")
	jwtExp, err := time.ParseDuration(jwtExpStr)
	if err != nil {
		panic(fmt.Errorf("failed to parse timeout: %s", err))
	}

	cfg.Credential.Jwt.ExpiredTime = jwtExp

	gracefullShutdownStr := viper.GetString("feature.gracefullShutdown")
	gracefullShutdown, err := time.ParseDuration(gracefullShutdownStr)
	if err != nil {
		panic(fmt.Errorf("failed to parse timeout: %s", err))
	}
	cfg.Feature.GracefullyShutdown = gracefullShutdown
	return
}
