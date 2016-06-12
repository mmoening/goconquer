package db

import "lib/structures"

// Kernel is an anonymously defined variable which contains global variable 
// definitions and collections. These global collections pool server information
// and information from the flat-file database, both used during server processing.
var Kernel struct {
	GameServers map[string]*structures.GameServer
}