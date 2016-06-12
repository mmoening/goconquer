package packets

// MsgName is a variable length packet sent between the client and game server to
// request, set and modify certain strings (such as spouse name). It can also be 
// used to display client assets in game (effects, sounds, etc).
// http://conquer.wiki/doku.php?id=msgname
type MsgName struct { 
	PacketHeader
	Identity uint32
	Action byte
	Strings []string
}