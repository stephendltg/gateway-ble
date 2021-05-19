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
}
