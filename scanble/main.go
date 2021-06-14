package scanble

import (
	"context"
	"encoding/hex"
	"fmt"
	"gateway-ble/beacon"
	"gateway-ble/store"
	"strconv"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Start Scanner BLE
func Start(device string, duration time.Duration, duplicate bool) {

	// debugger
	debugger := log.WithFields(log.Fields{"package": "SCANBLE"})

	d, err := dev.NewDevice(device)
	if err != nil {
		debugger.Fatal("can't new device : ", err)
	}
	ble.SetDefaultDevice(d)

	// Scan for specified durantion, or until interrupted by user.
	// debugger.Info("Scanning for ", duration)
	fmt.Println("Scan BLE start for: ", duration)

	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), duration))

	// If duration  == 0 -> infinity
	if int64(duration/time.Millisecond) == 0 {
		ctx = ble.WithSigHandler(context.Background(), nil)
	}
	chkErr(ble.Scan(ctx, duplicate, advHandler, nil))

}

// Scan ble handler
func advHandler(a ble.Advertisement) {

	// debugger
	debugger := log.WithFields(log.Fields{"package": "SCANBLE:HANDLER"})

	// RSSI filter
	rssi, _ := strconv.Atoi(store.Get("filter:rssi", nil))
	if (rssi * -1) > a.RSSI() {
		return
	}

	// NAME filter
	if !strings.Contains(a.LocalName(), strings.Replace(store.Get("filter:name", nil), "%", " ", -1)) {
		return
	}

	// MAC FILTER
	if !strings.Contains(a.Addr().String(), strings.Replace(store.Get("filter:mac", nil), "%", " ", -1)) {
		return
	}

	// Vars
	Data := []string{}
	UUID := []string{}
	Services := []string{}
	ManufacturerData := ""
	SolicitedService := []string{}
	OverflowService := []string{}
	var Temperature float64
	var Humidity float64
	var Battery int64 = -1

	if len(a.ServiceData()) > 0 {
		for i := range a.ServiceData() {
			UUID = append(UUID, a.ServiceData()[i].UUID.String())
			Data = append(Data, hex.EncodeToString([]byte(a.ServiceData()[i].Data)))
		}
	}
	if len(a.Services()) > 0 {
		for i := range a.Services() {
			Services = append(Services, hex.EncodeToString([]byte(a.Services()[i])))
		}
	}
	if len(a.ManufacturerData()) > 0 {
		ManufacturerData = hex.EncodeToString([]byte(a.ManufacturerData()))
	}
	if len(a.SolicitedService()) > 0 {
		for i := range a.SolicitedService() {
			SolicitedService = append(SolicitedService, hex.EncodeToString([]byte(a.SolicitedService()[i])))
		}
	}
	if len(a.OverflowService()) > 0 {
		for i := range a.OverflowService() {
			OverflowService = append(OverflowService, hex.EncodeToString([]byte(a.OverflowService()[i])))
		}
	}

	// ParserBeacons ELA / MINEW
	if len(UUID) > 0 {

		// Associated data to UUID
		data := Data[0]

		switch UUID[0] {

		case "2a6e":
			// Beacon ELA PT
			// Flip hexa
			hexaFlip := []string{}
			for i := len(data); i > 0; i -= 2 {
				hexaFlip = append(hexaFlip, data[i-2:i])
			}
			data = strings.Join(hexaFlip, "")
			f, _ := strconv.ParseInt(data, 16, 64)
			Temperature = float64(f) / 100

		case "ffe1":
			// Beacon minew
			b, _ := strconv.ParseInt(data[4:6], 16, 64)
			Battery = b
			t := data[2:4]

			//Beacon S1
			if t == "01" {
				f, _ := strconv.ParseInt(data[6:8], 16, 64)
				d, _ := strconv.ParseInt(data[8:10], 16, 64)
				Temperature = float64(f) + (float64(d)*0.4)/100
				h, _ := strconv.ParseInt(data[10:12], 16, 64)
				i, _ := strconv.ParseInt(data[12:14], 16, 64)
				Humidity = float64(h) + (float64(i)*0.4)/100
			}

		default:

		}
	}

	// ParserBeacons manufacture
	if len(ManufacturerData) > 0 {

		// fmt.Println(len(ManufacturerData))

		// Battery radius network E4 / E2 (protocol AltBeacon)
		if !a.Connectable() && ManufacturerData[0:4] == "1801" {
			l := len(ManufacturerData)
			b, _ := strconv.ParseInt(ManufacturerData[l-2:l], 16, 64)
			Battery = b
		}

		// Beacons Minew S3
		// if ManufacturerData[0:4] == "3906" {

		// }
	}

	m := beacon.Beacon{
		Datetime:         time.Now().Format("2006-01-02T15:04:05.000Z"),
		Mac:              a.Addr().String(),
		Rssi:             a.RSSI(),
		Name:             a.LocalName(),
		Connectable:      a.Connectable(),
		TxPower:          a.TxPowerLevel(),
		UUID:             UUID,
		DATA:             Data,
		Services:         Services,
		ManufacturerData: ManufacturerData,
		SolicitedService: SolicitedService,
		OverflowService:  OverflowService,
		Temperature:      Temperature,
		Humidity:         Humidity,
		Battery:          Battery,
	}

	debugger.Info(m)

	// Add beacon to store
	if store.Get("mqtt:status", nil) == "connected" {
		trame := ""
		if a.Connectable() {
			trame = "@info"
		}
		if len(UUID) > 0 {
			trame += "@" + strings.Join(UUID, ":")
		}
		store.AddBeacon(a.Addr().String()+trame, m)
	}

	// Write to DB
	if len(store.Get("db:host", nil)) > 0 {
		go Write(m)
	}

}

// Check error
func chkErr(err error) {

	// debugger
	debugger := log.WithFields(log.Fields{"package": "SCANBLE:DEVICE"})

	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Println("Scanner BLE : Done ....")
	case context.Canceled:
		fmt.Println("Scanner BLE : Canceled ....")
	default:
		debugger.Fatal("can't new device: ", err)
	}
}

// Write InfluxDB
func Write(c beacon.Beacon) {

	// debugger
	debugger := log.WithFields(log.Fields{"package": "DB"})

	client := influxdb2.NewClient(store.Get("db:host", nil), fmt.Sprintf("%s:%s", store.Get("db:user", nil), store.Get("db:pass", nil)))

	writeAPI := client.WriteAPIBlocking("", "gateway/autogen")

	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "rssi", "mac": c.Mac, "Name": c.Name},
		map[string]interface{}{"rssi": c.Rssi, "txpower": c.TxPower, "temperature": c.Temperature, "humidity": c.Humidity, "battery": c.Battery},
		time.Now())
	// Write data
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		debugger.Warn("Write error: ", err.Error())
	}

	// Close client
	client.Close()

}
