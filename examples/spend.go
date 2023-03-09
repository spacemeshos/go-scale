package examples

//go:generate scalegen

type Spend struct {
	Type uint8
	Body SpendBody
}

type SpendBody struct {
	Address  [20]byte
	Selector uint8
	Payload  SpendPayload
}

type SpendPayload struct {
	Arguments SpendArguments
	Nonce     SpendNonce
	GasPrice  uint32
	Signature [64]byte
}

type SpendArguments struct {
	Recipient [20]byte
	Amount    uint64
}

type SpendNonce struct {
	Counter  uint32
	Bitfield uint64
}
