package security

// Cipher allows the server to specify different cipher algorithms between servers 
// and states of the server for packet encryption. This is used on patches which 
// implement a second cipher on the game server, such as on patches higher than 
// 5017. 
type Cipher interface {
	Init()
	Generate(token, identity uint32)
	Decrypt(buffer []byte)
	Encrypt(buffer []byte)
}

func rotl(num, bits uint32) uint32 {
	return (num << (bits % 32)) | (num >> (32 - (bits % 32)))
}

func rotr(num, bits uint32) uint32 {
	return (num >> (bits % 32)) | (num << (32 - (bits % 32)))
}