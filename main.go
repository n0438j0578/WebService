package main

import (
	"flag"
	"fmt"
	golog "log"
	"net/http"
	"os"
	"strings"
	"time"

	"WebService/config"
	"WebService/controller"
	"WebService/data"
	"WebService/log"
)

const TAG = "Primes"

var (
	VERSION    = "0.0.1"
	help       = flag.Bool("help", false, "Displays this help")
	version    = flag.Bool("version", false, "Print version information and exit")
	configfile = flag.String("config", "etc/primes.conf", "Path to the configuration file")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	if *version {
		fmt.Printf("%s version %s\n", TAG, VERSION)
		return
	}
	err := config.LoadConfig(*configfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config: %v\n", err)
		os.Exit(1)
	}
	golog.SetFlags(golog.LstdFlags)
	log.SetLogLevel(config.GetLogConfig().LogLevel)
	logPath := config.GetLogConfig().LogPath
	if logPath != "stderr" {
		if logPath == "stdout" {
			golog.SetOutput(os.Stdout)
		} else {
			index := strings.Index(logPath, ":")
			if index < 1 {
				fmt.Fprintln(os.Stderr, "Write out log into the stderr or stdout or file(file:test.log)")
			}
			protocol := logPath[:index]
			path := logPath[index+1:]
			if protocol == "file" {
				err := log.OpenLogFile(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v", err)
					os.Exit(1)
				}
				defer log.CloseLogFile()
			}
		}
	}
	log.LogInfo("%s %v starting...", TAG, VERSION)
	log.LogInfo("pidof (%v) ", os.Getpid())

	ds := data.NewDataStore()
	ds.StorageConnect()

	s := &http.Server{
		Addr:           ds.HttpConfig.Address,
		Handler:        controller.Router(ds),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.LogInfo(ds.HttpConfig.Domain)
	if ds.HttpConfig.Cert != "" && ds.HttpConfig.Key != "" {
		log.LogInfo("Listen And Serve TLS")
		err = s.ListenAndServeTLS(ds.HttpConfig.Cert, ds.HttpConfig.Key)
	} else {
		log.LogInfo("Listen And Serve")
		err = s.ListenAndServe()
	}
	if err != nil {
		log.LogError("Https Server Error: %s", err)
	}
}
