package packets

// MsgUserInfo is sent from the game server to the game client. It contains 
// character information from the game database. In response to this packet, the 
// client will send a character location request.
// http://conquer.wiki/doku.php?id=msguserinfo
type MsgUserInfo struct {
	PacketHeader
	Identity, Mesh                      uint32
	Hairstyle, _                        uint16
	Silver               				uint32
	Experience, _                    	uint64 // TutorExp & MercenaryExp unused.
	_ 									uint32 // Potential unused.
	Strength, Agility, Vitality, Spirit uint16
	Attributes, Health, Mana, PkPoints  uint16
	Level, Class						byte
	Autoallot							bool
	Rebirths							byte
	ShowName							bool
	Strings								[]string
}

func NewMsgUserInfo() *MsgUserInfo {
	p := &MsgUserInfo{}
	p.Identifier = MSGUSERINFO
	p.ShowName = true
	return p
}