package store

// Vars
var DB = make(map[string]string)

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
