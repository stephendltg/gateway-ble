package store

import (
	"gateway-ble/beacon"
	"sync"
)

// Vars
var (
	DB      = make(map[string]string)
	Beacons = make(map[string]beacon.Beacon)
	Mutex   = sync.RWMutex{}
)

// Set in Store
func Set(index string, data string) {
	DB[index] = data
}

// Get in Store
func Get(index string, defaut *string) string {

	n := ""
	if defaut != nil {
		n = *defaut
	}

	if len(DB[index]) == 0 {
		return n
	} else {
		return DB[index]
	}
}

// Get All
func All() map[string]string {
	return DB
}

// Add beacon in list
func AddBeacon(index string, data beacon.Beacon) {
	Mutex.Lock()
	Beacons[index] = data
	Mutex.Unlock()
}

// Get All beacon
func BeaconsList() map[string]beacon.Beacon {
	// write & read IO
	Mutex.RLock()
	Beaconslist := Beacons
	Beacons = make(map[string]beacon.Beacon)
	Mutex.RUnlock()
	return Beaconslist
}
