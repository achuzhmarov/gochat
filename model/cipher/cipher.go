package cipher

type CipherType int

const (
	FirstLetterCipher CipherType = 1 + iota
	HalfWordRoundUpCipher
	HalfWordRoundDownCipher
	ShuffleCipher
)
