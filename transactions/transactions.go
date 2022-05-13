package transactions

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
	Signature [64]byte // it can't be a part of payload
}

type SelfSpawnArguments struct {
	Key scale.PublicKey
}
