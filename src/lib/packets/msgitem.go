package packets

// MsgItem is sent from the game client to the game server. Item actions, market
// interface actions, and some interface updates are handled through this packet. 
// Ping response is handled through this packet.
// http://conquer.wiki/doku.php?id=msgitem
type MsgItem struct {
	PacketHeader
	Identity, Argument, Action, Timestamp uint32
}