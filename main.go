package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gzook/obnodeman/lib/nodeman"

	"gopkg.in/natefinch/lumberjack.v2"
)

var nm *nodeman.Client

func main() {

	var (
		httpAddr = flag.String("http", ":3080", "HTTP service address (:3080 is default)")
		logfile  = flag.String("log", "obnodeman.log", "log file path (blank for no logging)")
	)
	flag.Parse()
	prepareLogger(*logfile)

	osExit := make(chan os.Signal, 1)
	signal.Notify(osExit, os.Interrupt)
	signal.Notify(osExit, syscall.SIGTERM)
	go func() {
		<-osExit
		shutdown()
		os.Exit(1)
	}()

	nm = nodeman.New()
	nm.Start()
	serve(*httpAddr)
}

func shutdown() {
	if nm != nil {
		nm.Stop()
	}
}
func serve(httpAddr string) {

	writeLog([]string{`msg`, `service started`, `address`, httpAddr})

	http.HandleFunc("/", handler)
	http.ListenAndServe(httpAddr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	writeLog([]string{`path`, r.URL.Path})

	lowerPath := strings.ToLower(r.URL.Path)
	switch lowerPath {
	case `/`:
		var message string
		if nm.Running() {
			message = "running"
		} else {
			message = "stopped"
		}
		writeAPIResponse(w, nil, message)
		break

	case `/stop`:
		err := nm.Stop()
		writeAPIResponse(w, err, "")
		break
	case `/start`:
		err := nm.Start()
		writeAPIResponse(w, err, "")
		break
	case `/restart`:
		message := ""
		err := nm.Stop()
		if err != nil {
			message = fmt.Sprintf("failed to stop: %s", err.Error())
		}

		err = nm.Start()
		writeAPIResponse(w, err, message)
		break
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("404 not found: %s", r.URL.Path)))
		break
	}
}

func writeAPIResponse(w http.ResponseWriter, err error, message string) {
	var r apiSimpleResponse
	if err != nil {
		r = apiSimpleResponse{Success: false, Error: err.Error(), Message: message}
	} else {
		r = apiSimpleResponse{Success: true, Message: message}
	}
	b := r.ToJSON()
	w.Write(*b)
}

func writeLog(values []string) {
	var buffer bytes.Buffer
	buffer.WriteString(time.Now().Format(`ts="2006-01-02 15:04:05.999"`))
	for i, v := range values {
		if i%2 == 0 {
			buffer.WriteString(` `)
			buffer.WriteString(v)
			buffer.WriteString(`=`)
		} else {
			buffer.WriteString(strconv.Quote(v))
		}
	}
	log.Println(buffer.String())
}

func prepareLogger(logFile string) {

	var logWriter io.Writer
	{
		if logFile == "" {
			logWriter = os.Stderr
		} else {
			lj := &lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    100, // megabytes
				MaxBackups: 3,
				MaxAge:     31, //days
			}

			logWriter = io.MultiWriter(os.Stderr, lj)
		}
	}

	log.SetOutput(logWriter)
}
