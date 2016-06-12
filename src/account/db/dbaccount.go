package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"lib/structures"
	"os"
)

// LoadAccount reads an account file from the database given a username from the 
// player's MsgAccount packet (sent by the client from the login screen). If no
// account exists, the function will return false.
func LoadAccount(acct *structures.Account, username string) bool {
	
	// Open the file and read stream.
	file, err := os.Open(fmt.Sprintf("./accounts/%s.json", username))
	if err != nil { return false }
	reader := bufio.NewReader(file)
	
	// Decode the JSON file into the structure passed.
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(acct)
	if err != nil { fmt.Println("failed to parse account file") }
	return err == nil
}
