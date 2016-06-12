package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

const (
	 AGILITY = 0
	 SPIRIT = 1
	 STRENGTH = 2
	 VITALITY = 3
)

// Attributes are points added to a character each time it levels up. Generally,
// points are automatically allocated by the server. When a character is reborn,
// the player can allocate points to desired attributes.
var Attributes attributes
type attributes struct {
	Archer  [121][4]uint16
	Toaist  [121][4]uint16
	Trojan  [121][4]uint16
	Warrior [121][4]uint16
}

// Load reads attributes from a JSON file in the flat-file database.
func (a *attributes) Load() bool {	
	fmt.Println("Loading attributes...")
	file, err := os.Open("./attributes.json")
	if err != nil { fmt.Println(err); return false }
	reader := bufio.NewReader(file)
	
	// Decode the JSON file into the structure passed.
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(a)
	if err != nil { fmt.Println(err); return false }
	return true
}

// Get returns an array of attributes for a class at the specified level. Use the
// constants defined above for accessing attributes from the array (AGILITY, 
// SPIRIT, STRENGTH, and VITALITY).
func (a *attributes) Get(class byte, level byte) [4]uint16 {
	if (level > 120) { return [4]uint16{0,0,0,0} }
	switch ((class / 10) * 10) {
		case 10: return a.Trojan[level]
		case 20: return a.Warrior[level]
		case 40: return a.Archer[level]
		case 100, 130, 140: return a.Toaist[level]
		default: return [4]uint16{0,0,0,0}
	}
}