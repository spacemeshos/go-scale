package transactions

import (
	"testing"

	"github.com/spacemeshos/go-scale/gen"
	"github.com/stretchr/testify/require"
)

func TestGenTransactions(t *testing.T) {
	require.NoError(t, gen.Generate("transactions", "transactions/transactions_scale.go",
		SelfSpawn{}, SelfSpawnArguments{}, SelfSpawnBody{}, SelfSpawnPayload{}))
}
