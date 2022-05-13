package types

import "github.com/spacemeshos/go-scale"

//go:generate scalegen -pkg types -file multi_scale.go -types SpawnMulti,SpawnMultiBody,SpawnMultiPayload,MultiSig,SpendMulti,SpendMultiBody,SpendMultiPayload -imports github.com/spacemeshos/go-scale/transactions/types

type SpawnMulti struct {
	Type uint8
	Body SpawnMultiBody
}

type SpawnMultiBody struct {
	Address  scale.Address
	Selector uint8
	Payload  SpawnMultiPayload
}

type SpawnMultiPayload struct {
	Template   scale.Address
	Keys       []scale.PublicKey `scale-type:"StructArray"`
	GasPrice   uint32
	Signatures MultiSig
}

type MultiSig struct {
	SigConf    uint8
	Signatures []scale.Signature `scale-type:"StructArray"`
}

type SpendMulti struct {
	Type uint8
	Body SpendMultiBody
}

type SpendMultiBody struct {
	Address  scale.Address
	Selector uint8
	Payload  SpendMultiPayload
}

type SpendMultiPayload struct {
	Arguments  SpendMethodArguments
	Nonce      SpendNonceFields
	GasPrice   uint64
	Signatures MultiSig
}
