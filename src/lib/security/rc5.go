package security
import "encoding/binary"

const ( // Definitions for manipulating the algorithm.
	RC5_Q uint = 0x61C88647
	RC5_P uint = 0xB7E15163
)

// RC5 is implemented in the account server for password verification. It was 
// removed around patch 5528, but can be reimplemented using hooks. Keys found in 
// this implementation are targeted for the Conquer Online game client for patches 
// lower than around 5176.
type RC5 struct {
	Key [4]uint32 
	Sub [26]uint32
}

// Init is called during authentication routine in the account server. Does not
// require input seed or prior initialization during OnConnect event unless using
// patches higher than around 5175.
func (cipher *RC5) Init() {
	
	// Initialize key and substitution box.
	cipher.Key = [...]uint32 { 0xE8FEDC3C, 0x7ED654C4, 0x1AF8A616, 0xBE38D0E8 }
	cipher.Sub[0] = 0xB7E15163
	for i := 1; i < 26; i++ { 
		cipher.Sub[i] = cipher.Sub[i - 1] - uint32(RC5_Q) 
	}
	
	// Generate the key vector according to RC5 specifications.
	// c = 3 * max(t, c). For this cipher, c = 3 * 26, c = 78.
	var A, B, i, j uint32 = 0, 0, 0, 0
	for c := 0; c < 78; c++ {
		A = rotl(cipher.Sub[i] + A + B, 3)
		cipher.Sub[i] = A
		B = rotl(cipher.Key[j] + A + B, A + B)
		cipher.Key[j] = B
		
		i = (i + 1) % 26
		j = (j + 1) % 4
	}
}

// Decrypt decrypts a password sent from the client to the account server. This 
// method should be avoided, unless you would like to remove the RC5 wrapping and 
// replace it with a more secure encryption or hash algorithm.
func (c *RC5) Decrypt(buffer []byte) {
	for block := 0; block < len(buffer) / 8; block++ {
		A := binary.LittleEndian.Uint32(buffer[8 * block:])
		B := binary.LittleEndian.Uint32(buffer[8 * block + 4:])
		
		for i := 12; i > 0; i-- {
			B = rotr(B - c.Sub[2 * i + 1], A) ^ A
			A = rotr(A - c.Sub[2 * i], B) ^ B
		}
		binary.LittleEndian.PutUint32(buffer[8 * block:],  A - c.Sub[0])
		binary.LittleEndian.PutUint32(buffer[8 * block + 4:], B - c.Sub[1])
	}
}

// Encrypt, you guested it, encrypts plain text into cipher text for client to 
// server, or encryption and verification on the server side. It can be compared
// to the cipher text sent my the client. This can be called by the web server
// to encrypt the password for storage. Should be hashed as well afterwards since
// RC5 isn't a strong encryption.
func (c *RC5) Encrypt(buffer []byte) {
	for block := 0; block < len(buffer) / 8; block++ {
		A := binary.LittleEndian.Uint32(buffer[8 * block:]) + c.Sub[0]
		B := binary.LittleEndian.Uint32(buffer[8 * block + 4:]) + c.Sub[1]
		
		for i := 12; i > 0; i-- {
			A = rotl(A ^ B, B) + c.Sub[2 * i]
			B = rotl(B ^ A, A) + c.Sub[2 * i + 1]
		}
		binary.LittleEndian.PutUint32(buffer[8 * block:],  A)
		binary.LittleEndian.PutUint32(buffer[8 * block + 4:], B)
	}
}