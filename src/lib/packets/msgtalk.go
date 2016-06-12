package packets

// MsgTalk is sent between the game client and server to exchange messages. The
// messages are used to control the client during login to the game server, and to
// socialize with other players on the server. The packet includes player meshes in
// patches 5017 and higher for private whisper windows. Suffix is used for offline
// messages, and usually contains the date in the format yyyyMMdd. 
// http://conquer.wiki/doku.php?id=msgtalk
type MsgTalk struct {
	PacketHeader
	Hue         uint32
	Tone, Style uint16
	Identity    uint32
	Strings     []string
}

func NewMsgTalk(identity uint32, sender, recipient, message string, 
	tone uint16) *MsgTalk {
	
	p := new(MsgTalk)
	p.Identifier = MSGTALK
	p.Hue = 0xFFFFFF
	p.Tone = tone
	p.Identity = identity

	p.Strings = make([]string, 4)
	p.Strings[0] = sender
	p.Strings[1] = recipient
	p.Strings[3] = message
	return p
}

const (
	MSGTALK_TALK            = 2000
	MSGTALK_WHISPER         = 2001
	MSGTALK_ACTION          = 2002
	MSGTALK_TEAM            = 2003
	MSGTALK_GUILD           = 2004
	MSGTALK_TOP_LEFT        = 2005
	MSGTALK_CLAN            = 2006
	MSGTALK_SERVER          = 2007
	MSGTALK_YELL            = 2008
	MSGTALK_FRIEND          = 2009
	MSGTALK_GLOBAL          = 2010
	MSGTALK_CENTER          = 2011
	MSGTALK_GHOST           = 2013
	MSGTALK_SERVICE         = 2014
	MSGTALK_TIP             = 2015
	MSGTALK_ENTRANCE    	= 2100
	MSGTALK_REGISTRATION    = 2101
	MSGTALK_SHOP            = 2102
	MSGTALK_PET             = 2103
	MSGTALK_VENDOR_HAWK     = 2104
	MSGTALK_WEBSITE         = 2105
	MSGTALK_TOP_RIGHT_FIRST = 2108
	MSGTALK_TOP_RIGHT_CONT  = 2109
	MSGTALK_GUILD_BULLETIN  = 2111
	MSGTALK_TRADE_BOARD     = 2201
	MSGTALK_FRIEND_BOARD    = 2202
	MSGTALK_TEAM_BOARD      = 2203
	MSGTALK_GUILD_BOARD     = 2204
	MSGTALK_OTHERS_BOARD    = 2205
)
