package db

import (
	"lib/threadsafe"
)
// Kernel is an anonymously defined variable which contains global variable 
// definitions and collections. These global collections pool server information
// and information from the flat-file database, both used during server processing.
var Kernel kernel
type kernel struct {
	AuthenticatedClients *threadsafe.SafeMap 
	ConnectedClients *threadsafe.SafeMap
	CharacterCreationPool *threadsafe.SafeMap
	CharacterIndex *threadsafe.SafeMap
}

// Init initializes global collections used by the server.
func (k *kernel) Init() {
	k.AuthenticatedClients = threadsafe.NewSafeMap()
	k.ConnectedClients = threadsafe.NewSafeMap()
	k.CharacterCreationPool = threadsafe.NewSafeMap()
	k.CharacterIndex = threadsafe.NewSafeMap()
}

