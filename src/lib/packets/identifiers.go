// Packets implements a simple translation between golang type structures and 
// Conquer Online binary packet structures, which are sequences of bytes for data 
// types of both fixed and varying sizes. This package favors simplicity over 
// efficiency. Applications which require high-performance serialization, 
// especially for large data structures, should look at more advanced solutions 
// such as protocol buffers. 
package packets

// PacketHeader is the binary header structure for NetDragon Websoft packets. If
// your client crashes or disconnects, it's likely due to an invalid length sent
// in this header.
type PacketHeader struct {
	Length, Identifier uint16
}

// Identifiers for packet structures.
const (
	MSGREGISTER  = 1001
	MSGTALK      = 1004
	MSGUSERINFO	 = 1006
	MSGITEM		 = 1009
	MSGACTION	 = 1010
	MSGACCOUNT   = 1051
	MSGCONNECT   = 1052
	MSGCONNECTEX = 1055
)
