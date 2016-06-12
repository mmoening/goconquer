package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lib/structures"
	"os"
)

// LoadGameServers loads game server forwarding information on server startup,
// used for transferring account information and redirecting the player's client 
// to the correct game server. Returns false on error.
func LoadGameServers() bool {
	
	// Read all files from the servers directory.
	fmt.Println("Loading game servers...")
	Kernel.GameServers = make(map[string]*structures.GameServer)
	files, err := ioutil.ReadDir("./servers")
	if err != nil { fmt.Println(err); return false }
	
	// Read from each file.
	for _, f := range files {
		file, err := os.Open(fmt.Sprintf("./servers/%s", f.Name()))
		if err != nil { fmt.Println(err); return false }
		reader := bufio.NewReader(file)
		
		// Decode the JSON file.
		server := &structures.GameServer {}
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(server)
		if err != nil { fmt.Println(err); return false }
		
		// Add to map of available servers.
		Kernel.GameServers[server.Name] = server
		go server.Connect()
	}
	return true
}