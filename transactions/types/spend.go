package types

import "github.com/spacemeshos/go-scale"

type Spend struct {
	Type uint8
	Body SpendBody
}

type SpendBody struct {
	Adress   scale.Address
	Selector uint8
	Payload  SpendPayload
}

type SpendPayload struct {
	Arguments SpendMethodArguments
	Nonce     SpendNonceFields
	GasPrice  uint32
	Signature scale.Signature
}

type SpendMethodArguments struct {
	Recipient scale.Address
	Amount    uint64
}

type SpendNonceFields struct {
	Counter  uint32
	Bitfield uint64
}
