package main

import (
	"flag"
	"fmt"
	"gateway-ble/broker"
	"gateway-ble/database"
	"gateway-ble/scanble"
	"gateway-ble/store"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	metrics "github.com/tevjef/go-runtime-metrics"
)

// Vars cli
var (
	device   = flag.String("device", "default", "Implementation of ble")
	duration = flag.Duration("du", 0, "Scanning duration") // ex: 5-time.Second
	dup      = flag.Bool("dup", true, "Allow duplicate reported")
	mqtt     = flag.String("mqtt", "", "MQTT hostname(:port<option>)")
	interval = flag.String("interval", "5s", "Interval duration publish messsage")
	influxdb = flag.String("db", "", "InfluxDB host (<http://127.0.0.1:8086>)")
	user     = flag.String("u", "gateway", "Username influxDB")
	pass     = flag.String("p", "gatewayis100%MAGIC", "Password influxDB")
	debug    = flag.Bool("debug", false, "Mode debug")
	rssi     = flag.String("rssi", "130", "Beacon RSSI filter")
	name     = flag.String("name", "", "Beacon name filter")
	mac      = flag.String("mac", "", "Beacon Mac adress filter")
	collect  = flag.Bool("collect", false, "Metrics runtime")
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

	debugger := log.WithFields(log.Fields{"package": "MAIN"})

	if *collect {
		err := metrics.RunCollector(metrics.DefaultConfig)

		if err != nil {
			debugger.Error(err)
		}
	}

	fmt.Println("Initialize ....")

	// Mqtt cron job
	store.Set("mqtt:interval", *interval)

	// Start broker if available
	mqttParams := len(strings.Split(*mqtt, ":"))
	if mqttParams > 1 {
		broker.Connect(&strings.Split(*mqtt, ":")[0], &strings.Split(*mqtt, ":")[1], nil, nil)
		fmt.Println("Broker: ", *mqtt)
	}

	// DB if available && check database
	db := len(strings.Split(*influxdb, ":"))
	if db > 1 {

		// Check database in influxDb
		client := &http.Client{}
		req, err := http.NewRequest("POST", *influxdb+"/query", strings.NewReader("q=CREATE%20DATABASE%20%22gateway%22"))
		if err != nil {
			debugger.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(req)
		if err != nil {
			debugger.Fatal(err)
		}
		if res.StatusCode == 200 {
			debugger.Info("Database gateway checked")
		} else {
			debugger.Warn("Database gateway error", res.StatusCode)
		}

		// Prepare configuration
		store.Set("db:host", *influxdb)
		store.Set("db:user", *user)
		store.Set("db:pass", *pass)
		fmt.Println("DB influxDB : ", *influxdb)

	}

	// Filter
	store.Set("filter:rssi", *rssi)
	store.Set("filter:name", *name)
	store.Set("filter:mac", strings.ToLower(*mac))

	// Connection database SQL
	database.Connect(*debug)

	// Don't Exit
	var wg sync.WaitGroup
	wg.Add(1)
	// Start Scan ble
	time.AfterFunc(2*time.Second, func() { scanble.Start(*device, *duration, *dup) })
	wg.Wait()

}
