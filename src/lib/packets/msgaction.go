package packets

// MsgAction is used to control and verify many different types of actions
// between the game client and game server. It contains the ID of an entity
// performing the action, a type of action being performed and a variable amount
// of information, often with various fields being used.
// http://conquer.wiki/doku.php?id=msgaction
type MsgAction struct {
	PacketHeader
	Timestamp, Identity, Data uint32
	X, Y, Direction           uint16
	Action                    MsgActionType
}

func NewMsgAction() *MsgAction {
	p := new(MsgAction)
	p.Identifier = MSGACTION
	return p
}

type MsgActionType uint16
const (
	ACTION_SETLOCATION  MsgActionType = 74
	ACTION_SETEQUIPMENT MsgActionType = 75
	ACTION_SETFRIENDS   MsgActionType = 76
	ACTION_SETSKILLS    MsgActionType = 77
	ACTION_SETSPELLS    MsgActionType = 78
	ACTION_SETDIRECTION MsgActionType = 79
	ACTION_SETPOSE    	MsgActionType = 81
	ACTION_USEPORTAL    MsgActionType = 85
	ACTION_USETELEPORT	MsgActionType = 86
	ACTION_SETLEVEL		MsgActionType = 92
	ACTION_USEXPSKILLS	MsgActionType = 93
	ACTION_USEREVIVE	MsgActionType = 94
	ACTION_DELETECHAR	MsgActionType = 95
	ACTION_SETPKMODE	MsgActionType = 96
	ACTION_SETRESPAWN	MsgActionType = 102
	ACTION_SETPOSITION	MsgActionType = 108
	ACTION_BOOTHSETUP	MsgActionType = 111
	ACTION_BOOTHSUSPEND	MsgActionType = 112
	ACTION_BOOTHRESUME	MsgActionType = 113
	ACTION_BOOTHLEAVE	MsgActionType = 114
	ACTION_OFFLINEMSGS	MsgActionType = 132
	ACTION_REMOVESPAWN	MsgActionType = 135
	ACTION_JUMP			MsgActionType = 137
)
