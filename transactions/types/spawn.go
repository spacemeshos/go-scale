package types

import "github.com/spacemeshos/go-scale"

type SelfSpawn struct {
	Type uint8
	Body SelfSpawnBody
}

type SelfSpawnBody struct {
	Address  scale.Address
	Selector uint8
	Payload  SelfSpawnPayload
}

type SelfSpawnPayload struct {
	Template  scale.Address
	Arguments SelfSpawnArguments
	GasPrice  uint32
	Signature scale.Signature
}

type SelfSpawnArguments struct {
	Key scale.PublicKey
}
