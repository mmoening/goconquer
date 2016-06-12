package structures

// Account is instantiated when a user logs into the account server. The structure 
// is populated by decoding the account JSON file from the database. 
type Account struct {
	Identity  uint32
	Username  string
	Password  string
	Authority uint32
	Status    uint32
}
	
const (
	ACCTSTATUS_OK = 0
	ACCTSTATUS_LOCKED = 1
	ACCTSTATUS_BANNED = 2
	ACCTSTATUS_LIKESUPERHELLABANNED = 3 // Matt is great.
)