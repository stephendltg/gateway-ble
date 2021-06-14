package beacon

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
	Temperature      float64
	Humidity         float64
	Battery          int64
}
