package scenjsonmodel

// DCDTTransfer models the transfer of tokens in a tx
type DCDTTxData struct {
	TokenIdentifier JSONBytesFromString
	Nonce           JSONUint64
	Value           JSONBigInt
}

// DCDTInstance models an instance of an NFT/SFT, with its own nonce
type DCDTInstance struct {
	Nonce   JSONUint64
	Balance JSONBigInt
	Uri     JSONBytesFromTree
}

// DCDTData models an account holding an DCDT token
type DCDTData struct {
	TokenIdentifier JSONBytesFromString
	Instances       []*DCDTInstance
	LastNonce       JSONUint64
	Roles           []string
	Frozen          JSONUint64
}

// CheckDCDTInstance checks an instance of an NFT/SFT, with its own nonce
type CheckDCDTInstance struct {
	Nonce   JSONCheckUint64
	Balance JSONCheckBigInt
	Uri     JSONCheckBytes
}

// CheckDCDTData checks the DCDT tokens held by an account
type CheckDCDTData struct {
	TokenIdentifier JSONBytesFromString
	Instances       []*CheckDCDTInstance
	LastNonce       JSONCheckUint64
	Roles           []string
	Frozen          JSONCheckUint64
}
