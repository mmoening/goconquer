package packets

// MsgAccount is sent during the client's authentication state. It's sent to the 
// account server to verify the user's inputted password. It also contains the 
// account name and game server name being requested. 
// http://conquer.wiki/doku.php?id=msgaccount
type MsgAccount struct {
	PacketHeader
	Account  string `len:"16"`
	Password [16]byte
	Server   string `len:"16"`
}