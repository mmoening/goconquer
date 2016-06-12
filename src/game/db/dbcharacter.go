package db

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"lib/structures"
	"os"
	"strconv"
	"sync"
)

// Characters are playable entities on the server, controlled by players.
// This structure defines actions for managing character files in the
// flat-file database.
var Characters characters
type characters struct { }
var indexlock sync.Mutex

// Save encodes a character to JSON and saves it to the flat-file database.
func (_ *characters) Save(c *structures.Character) bool {
	
	// BUG(Gareth): Race condition could occur with server shutdown.
	file, err := os.Create("./characters/" + c.Name + ".json")
	defer file.Close()
	if err != nil { 
		fmt.Printf("error: open character file for %s\n", c.Name)
		return false
	}
	
	// Encode to the new file.
	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(c)
	defer writer.Flush()
	if err != nil { fmt.Println(err); return false }
	return true
} 

// Load opens a character file from the flat-file database after performing a
// lookup from the character index, which maps character ids to file names.
func (_ *characters) Load(c *structures.Client) (bool, error) {
	
	name := Kernel.CharacterIndex.Get(c.Identity)
	if name == nil { return false, nil }
	file, err := os.Open("./characters/" + name.(string) + ".json")
	if err != nil { 
		fmt.Printf("error: open character file for %s\n", name.(string))
		return false, err 
	}
	
	c.Character = new(structures.Character)
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(c.Character)
	if err != nil { 
		fmt.Printf("error: parse character file for %s\n", name.(string))
		return false, err 
	} else { return true, nil }
}

// LoadIndex initializes the character index map from a file in the flat-file
// database. This file is required for proper indexing of character names
// by character identites in memory. The file is modified on character 
// creation and referenced in memory on login.
func (_ *characters) LoadIndex() bool {
	
	// Open the file and read all entries.
	file, err := os.Open("./characters/index.csv")
	if err != nil { 
	
		// Attempt to create the file.
		file, err = os.Create("./characters/index.csv")
		if err != nil { 
			fmt.Println("error: could not open character index.csv")
			return false
		}
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	
	// Add entries to the map.
	for _, row := range records {
		identity, _ := strconv.Atoi(row[0])
		name := row[1]
		Kernel.CharacterIndex.Add(uint32(identity), name)
	} 
	return true
}

// AppendIndex adds a new index entry to the end of the index file, then adds
// the new mapping to memory.
func (_ *characters) AppendIndex(id uint32, name string) bool {
	indexlock.Lock()
	file, err := os.OpenFile("./characters/index.csv", 
		os.O_WRONLY | os.O_APPEND, 0660)
	if err != nil { fmt.Println(err); return false }
	
	// Append to the file and close.
	fmt.Fprintf(file, "%d,%s\r\n", id, name)
	file.Close()
	indexlock.Unlock()
	
	// Add to the character index map in memory.
	Kernel.CharacterIndex.Add(id, name)
	return true
}