package security

const ( // Definitions for manipulating the algorithm.
	TQCS_KEY1 uint = 0x13FA0F9D
	TQCS_KEY2 uint = 0x6D5C7962
)

// TQCipher is an in-house asymmetric xor-cipher implementation for the account 
// and game server's packet transactions. This legacy cipher was used in older 
// clients until patch 5018 (excluding Conquer Online 1.0 Alpha), where it was 
// replaced by Blowfish for the game server.
type TQCipher struct {
	key2, key1   [0x200]byte
	key          *[0x200]byte
	encryptcount uint64
	decryptcount uint64
}

// Init generates the first key vector which is later used to generate the second 
// key vector and necessary for decrypting and encrypting packets for the account 
// server.
func (c *TQCipher) Init() {
	
	// Initialize variables for generating the first key.
	a, b := TQCS_KEY1, TQCS_KEY2
	p, g := make([]byte, 4), make([]byte, 4)
	p[3] = byte(a >> 24);
	p[2] = byte(a >> 16);
	p[1] = byte(a >> 8);
	p[0] = byte(a);
	
	g[3] = byte(b >> 24);
	g[2] = byte(b >> 16);
	g[1] = byte(b >> 8);
	g[0] = byte(b);
	
	// Generate the initialization vector.
	for i := 0; i < 0x100; i++ {
		
		c.key1[i] = p[0]
		c.key1[i + 0x100] = g[0]
		p[0] = (p[1] + byte(p[0] * p[2])) * p[0] + p[3]
		g[0] = (g[1] - byte(g[0] * g[2])) * g[0] + g[3]
	}
	c.key = &c.key1;
}

// Generate is called to generate the key vector after decrypting the first packet 
// for the game server. The first packet should be MsgConnect, which contains the 
// cipher's token and identity values (used in generating the master key and 
// squared master key).
func (c *TQCipher) Generate(token, identity uint32) {
	
	// These are awful. Initialize variables for key generation. Please never 
	// write your own cipher algorithm.
	temp1 := int(((token + identity) ^ 0x4321) ^ token)
	temp2 := int(temp1 * temp1)
	
	xorkey, squared := make([]byte, 4), make([]byte, 4)
	xorkey[3] = byte(temp1 >> 24);
	xorkey[2] = byte(temp1 >> 16);
	xorkey[1] = byte(temp1 >> 8);
	xorkey[0] = byte(temp1);
	
	squared[3] = byte(temp2 >> 24);
	squared[2] = byte(temp2 >> 16);
	squared[1] = byte(temp2 >> 8);
	squared[0] = byte(temp2);
	
	// Generate the key vector.
	for i := uint(0); i < 0x100; i++ {
		c.key2[i] = c.key1[i] ^ xorkey[i % 4]
		c.key2[i + 0x100] = c.key1[i + 0x100] ^ squared[i % 4]
	}
	c.key = &c.key2;
}

// Decrypt requires that init is called before first time use. After the first 
// decrypt on the game server, the generate method must be called to continue 
// decrypting client packets. 
func (c *TQCipher) Decrypt(buffer []byte) {
	for i := 0; i < len(buffer); i++ {
		
		buffer[i] ^= 0xAB
		buffer[i] = buffer[i] >> 4 | buffer[i] << 4
		buffer[i] ^= c.key[byte(c.decryptcount & 0xFF)]
		buffer[i] ^= c.key[int(byte(c.decryptcount >> 8)) + 0x100]
		c.decryptcount++
	}
}

// Encrypt requires that init is called before first time use. The game server 
// will not use the key vector in this method; therefore, the client can send 
// packets to the game server to reject the connection if the token and identity 
// for the generate method are invalid.
func (c *TQCipher) Encrypt(buffer []byte) {
	for i := 0; i < len(buffer); i++ {
		
		buffer[i] ^= 0xAB
		buffer[i] = buffer[i] >> 4 | buffer[i] << 4
		buffer[i] ^= c.key1[byte(c.encryptcount & 0xFF)]
		buffer[i] ^= c.key1[int(byte(c.encryptcount >> 8)) + 0x100]
		c.encryptcount++
	}
}