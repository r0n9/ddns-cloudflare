package cmd

import (
	"encoding/json"
	"github.com/r0n9/ddns-cloudflare/cmd/flags"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Conf struct {
	Email   string   `json:"email"`
	ZoneId  string   `json:"zoneId"`
	ApiKey  string   `json:"apiKey"`
	Domains []string `json:"domains"`

	SendKey string `json:"sendKey"`
}

var Config *Conf

func Init() {
	formatter := log.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		TimestampFormat:           "2006-01-02 15:04:05",
		FullTimestamp:             true,
	}
	log.SetFormatter(&formatter)

	// init configurations
	configPath := flags.Conf
	log.Infof("loading config file: %s", configPath)
	if !Exists(configPath) {
		log.Fatalf("config file not exists")
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("reading config file error: %+v", err)
	}
	Config = &Conf{}
	err = json.Unmarshal(configBytes, Config)
	if err != nil {
		log.Fatalf("load config error: %+v", err)
	}
	log.Printf("loaded config file: %+v ", Config)
}

var pid = -1
var pidFile string

func initDaemon() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	_ = os.MkdirAll(filepath.Join(exPath, "daemon"), 0700)
	pidFile = filepath.Join(exPath, "daemon/pid")
	if Exists(pidFile) {
		bytes, err := os.ReadFile(pidFile)
		if err != nil {
			log.Fatal("failed to read pid file", err)
		}
		id, err := strconv.Atoi(string(bytes))
		if err != nil {
			log.Fatal("failed to parse pid data", err)
		}
		pid = id
	}
}

// Exists determine whether the file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
