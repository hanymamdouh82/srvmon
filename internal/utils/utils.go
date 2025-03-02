package utils

import (
	"encoding/json"
	"fmt"
)

// Logs JSON representation of any objects on stdout
func LogJSON(val interface{}) {
	j, _ := json.MarshalIndent(val, "", "    ")
	fmt.Printf("%s\n", j)
}
