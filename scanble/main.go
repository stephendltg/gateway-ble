package scanble

import (
	"context"
	"encoding/hex"
	"fmt"
	"gateway-ble/broker"
	"gateway-ble/store"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Vars
type Beacon struct {
	Datetime         string
	Mac              string
	Rssi             int
	Name             string
	Connectable      bool
	TxPower          int
	UUID             []string
	DATA             []string
	Services         []string
	ManufacturerData string
	SolicitedService []string
	OverflowService  []string
}

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

	// if a.Addr().String() != "ac:23:3f:a2:7e:ce" {
	// 	return
	// }

	Data := []string{}
	UUID := []string{}
	Services := []string{}
	ManufacturerData := ""
	SolicitedService := []string{}
	OverflowService := []string{}

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

	m := Beacon{
		time.Now().Format("2006-01-02T15:04:05.000Z"),
		a.Addr().String(),
		a.RSSI(),
		a.LocalName(),
		a.Connectable(),
		a.TxPowerLevel(),
		UUID,
		Data,
		Services,
		ManufacturerData,
		SolicitedService,
		OverflowService,
	}

	debugger.Info(m)
	broker.Publish(m)

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
func Write(c Beacon) {

	// debugger
	debugger := log.WithFields(log.Fields{"package": "DB"})

	client := influxdb2.NewClient(store.Get("db:host", nil), fmt.Sprintf("%s:%s", store.Get("db:user", nil), store.Get("db:pass", nil)))

	writeAPI := client.WriteAPIBlocking("", "epyo/autogen")

	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "rssi", "mac": c.Mac, "Name": c.Name},
		map[string]interface{}{"rssi": c.Rssi, "txpower": c.TxPower},
		time.Now())
	// Write data
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		debugger.Warn("Write error: ", err.Error())
	}

	// Close client
	client.Close()

}
