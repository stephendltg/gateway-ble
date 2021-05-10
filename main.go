package main

import (
	"flag"
	"fmt"
	"gateway-ble/broker"
	"gateway-ble/scanble"
	"gateway-ble/store"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Vars cli
var (
	device   = flag.String("device", "default", "implementation of ble")
	duration = flag.Duration("du", 0, "scanning duration") // ex: 5-time.Second
	dup      = flag.Bool("dup", true, "allow duplicate reported")
	mqtt     = flag.String("mqtt", "", "hostname(:port<option>)")
	influxdb = flag.String("db", "", "host(<http://127.0.0.1:8086>)")
	user     = flag.String("u", "epyo", "Username influxDB")
	pass     = flag.String("p", "epyois100%MAGIC", "Password influxDB")
	debug    = flag.Bool("debug", false, "Mode debug")
)

// Init
func init() {

	// Datetime format logger
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02T15:04:05.000Z"
	logrus.SetFormatter(&log.JSONFormatter{})
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Parse cli
	flag.Parse()

	// Debug mode
	if !*debug {
		log.SetLevel(log.WarnLevel)
	}

	// log.Info(reflect.TypeOf(duration))
}

func main() {

	fmt.Println("Initialize ....")

	// Start broker if available
	mqttParams := len(strings.Split(*mqtt, ":"))
	if mqttParams > 1 {
		broker.Connect(&strings.Split(*mqtt, ":")[0], &strings.Split(*mqtt, ":")[1], nil, nil)
		fmt.Println("Broker: ", *mqtt)
	}

	// DB if available
	db := len(strings.Split(*influxdb, ":"))
	if db > 1 {
		store.Set("db:host", *influxdb)
		store.Set("db:user", *user)
		store.Set("db:pass", *pass)
		fmt.Println("DB influxDB : ", *influxdb)
	}

	// Don't Exit
	var wg sync.WaitGroup
	wg.Add(1)
	// Start Scan ble
	time.AfterFunc(2*time.Second, func() { scanble.Start(*device, *duration, *dup) })
	wg.Wait()
}
