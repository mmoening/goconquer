package packets

// MsgRegister is sent from the game client to the game server. It contains 
// customized character information from the character creation screen. This 
// information needs to be verified with the server to ensure that the name is not
// taken and contains valid characters, and that the body and class are valid.
// http://conquer.wiki/doku.php?id=msgregister
type MsgRegister struct {
	PacketHeader
	Account, Name, Password string `len:"16"`
	Model, Class uint16
	Identity uint32
}