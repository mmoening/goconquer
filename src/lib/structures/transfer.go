package structures

import "time"

// Transfer defines authentication transfer between the Account Server and Game 
// Server. Usually, transfer is sent through the client which directly exposes the 
// session token and account id. If done incorrectly, this opens vulnerability for 
// session hijacking or, even worse, bypassing authentication for any user. 
// GoConquer makes attempts to avoid this by sending this structure over a backend
// channel to the other server.
type Transfer struct {
	Account   Account
	IPAddress string
	Requested time.Time
}
