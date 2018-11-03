package config

import (
	"container/list"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/robfig/config"
)

type LogConfig struct {
	LogLevel string
	LogPath  string
}

type MongoDBConfig struct {
	MongoUri string
	MongoDb  string
}

type HttpConfig struct {
	Address     string
	Domain      string
	Cert        string
	Key         string
	SessionTime time.Duration
}

var (
	Config        *config.Config
	logConfig     *LogConfig
	mongoDBConfig *MongoDBConfig
	httpConfig    *HttpConfig
)

func GetLogConfig() *LogConfig {
	return logConfig
}

func GetMongoDBConfig() *MongoDBConfig {
	return mongoDBConfig
}

func GetHttpConfig() *HttpConfig {
	return httpConfig
}

func LoadConfig(filename string) error {
	var err error
	Config, err = config.ReadDefault(filename)
	if err != nil {
		return err
	}
	messages := list.New()

	RequireSection(messages, "logging")
	RequireOption(messages, "logging", "level")
	if messages.Len() > 0 {
		fmt.Fprintln(os.Stderr, "Error(s) validating configuration:")
		for e := messages.Front(); e != nil; e = e.Next() {
			fmt.Fprintln(os.Stderr, " -", e.Value.(string))
		}
		return fmt.Errorf("Failed to validate configuration")
	}

	if err = ParseLoggingConfig(); err != nil {
		return err
	}
	if err = ParseMongoDBConfig(); err != nil {
		return err
	}
	if err = ParseHttpConfig(); err != nil {
		return err
	}

	return nil
}

func ParseLoggingConfig() error {
	logConfig = new(LogConfig)
	section := "logging"

	option := "level"
	str, err := Config.String(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	switch strings.ToUpper(str) {
	case "TRACE", "INFO", "WARN", "ERROR":
	default:
		return fmt.Errorf("Invalid value provided for [%v]%v: '%v'", section, option, str)
	}
	logConfig.LogLevel = str

	option = "path"
	str, err = Config.String(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	logConfig.LogPath = str

	return nil
}

func ParseMongoDBConfig() error {
	mongoDBConfig = new(MongoDBConfig)
	section := "mongodbstore"

	option := "mongo.uri"
	str, err := Config.String(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	mongoDBConfig.MongoUri = str

	option = "mongo.db"
	str, err = Config.String(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	mongoDBConfig.MongoDb = str

	return nil
}

func ParseHttpConfig() error {
	httpConfig = new(HttpConfig)
	section := "http"

	option := "address"
	str, err := Config.String(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	httpConfig.Address = str

	option = "domain"
	str, err = Config.String(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	httpConfig.Domain = str

	option = "cert"
	str, err = Config.String(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	httpConfig.Cert = str

	option = "key"
	str, err = Config.String(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	httpConfig.Key = str

	option = "sessionTime"
	number, err := Config.Int(section, option)
	if err != nil {
		return fmt.Errorf("Failed to parse [%v]%v: '%v'", section, option, err)
	}
	httpConfig.SessionTime = time.Duration(number) * time.Minute

	return nil
}

func RequireSection(messages *list.List, section string) {
	if !Config.HasSection(section) {
		messages.PushBack(fmt.Sprintf("Config section [%v] is required", section))
	}
}

func RequireOption(messages *list.List, section string, option string) {
	if !Config.HasOption(section, option) {
		messages.PushBack(fmt.Sprintf("Config option '%v' is required in section [%v]", option, section))
	}
}
