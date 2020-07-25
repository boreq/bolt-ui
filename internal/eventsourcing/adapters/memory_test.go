package adapters_test

import (
	"testing"

	"github.com/boreq/velo/internal/eventsourcing/adapters"
)

func RunTestMemory(t *testing.T, test Test) {
	adapter := adapters.NewMemoryPersistenceAdapter()
	test(t, adapter)
}
