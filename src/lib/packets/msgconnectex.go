package packets

// MsgConnectEx is sent during the client's authentication state to either accept
// or reject the client's attempt to login to the selected game server. The 
// authentication code and identity are XORed by 0x4321 in the original client; 
// however, this method of “protection” will not suffice. The game server IP 
// address cannot be a localhost address (127.x.x.x). Local network and public 
// network addresses are accepted. http://conquer.wiki/doku.php?id=msgconnectex 
type MsgConnectEx struct { 
	PacketHeader
	Identity, Token uint32
	Address [16]byte
	Port uint32
}

func NewMsgConnectEx() *MsgConnectEx { 
	p := new(MsgConnectEx)
	p.Identifier = MSGCONNECTEX
	return p
}

// These constants are encoded using GB2312 (Simplified Chinese). Modifying the 
// message will not work, as the client string matches these messages with the 
// received rejection id. This method was used in all older clients with the old 
// login interface (before the animated flash login dialog was introduced around 
// patch 5028). After that, it only uses the rejection id. 
var MSGCONNECTEX_INVALID_ACCOUNT = []byte { // 帐号名或口令错 (1)
	0xD5, 0xca, 0xba, 0xc5, 0xc3, 0xfb, 0xbb, 0xf2,
	0xbf, 0xda, 0xc1, 0xee, 0xb4, 0xed, 0x00, 0x00 }
var MSGCONNECTEX_SERVER_DOWN = []byte { // 服务器未启动 (10)
	0xb7, 0xfe, 0xce, 0xf1, 0xc6, 0xf7, 0xce, 0xb4, 
	0xc6, 0xf4, 0xb6, 0xaf }
var MSGCONNECTEX_LOGIN_LATER = []byte { // 请稍后重新登录 (11)
	0xc7, 0xeb, 0xc9, 0xd4, 0xba, 0xf3, 0xd6, 0xd8, 
	0xd0, 0xc2, 0xb5, 0xc7, 0xc2, 0xbc, 0x00, 0x00 }
var MSGCONNECTEX_BANNED_ACCOUNT = []byte { // 该帐号被封号 (12)
	0xb8, 0xc3, 0xd5, 0xca, 0xba, 0xc5, 0xb1, 0xbb,
	0xb7, 0xe2, 0xba, 0xc5, 0x00, 0x00, 0x00, 0x00 }