package db

import (
	"bufio"
	"encoding/json"
	"os"
)

// Configuration is a privately defined variable which encapsulates the server's
// configuration details. The struct is loaded in on runtime from a JSON file 
// found in the same directory as the executable.
var Configuration configuration
type configuration struct {
	Host string
}

// Decode is called from the main function to load the server's json configuration
// file. It uses a decoding stream to parse the file into a configuration 
// structure, globally defined as Configuration.
func (c *configuration) Decode(path string) error {
	
	// Open the configuration file and read stream.	
	file, err := os.Open(path)
	if err != nil { return err }
	reader := bufio.NewReader(file)
	
	// Decode the JSON file into the structure passed.
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(c)
	if err != nil { return err }
	return nil
}