package packets

// MsgConnect is created by the client after it processes the MsgConnectEx packet,
// disconnecting from the account server, and connecting to the game server. It 
// sends this packet to the game server to pass authenticated information from the
// account server. It also sends a second packet to the account server with the 
// contents of Res.dat. http://conquer.wiki/doku.php?id=msgconnect
type MsgConnect struct {
	PacketHeader
	Identity, Token uint32
	Version string `len:"4"`
	Language string `len:"12"`
}