package tests

import (
	"testing"

	"github.com/boreq/eggplant/application/tracker"
	"github.com/boreq/eggplant/internal/fixture"
	"github.com/boreq/eggplant/internal/wire"
	"github.com/stretchr/testify/require"
)

func TestAddActivity(t *testing.T) {
	tr, cleanup := NewTracker(t)
	defer cleanup()

	cmd := tracker.AddActivity{}

	err := tr.AddActivity.Execute(cmd)
	require.EqualError(t, err, "not implemented")
}

func NewTracker(t *testing.T) (*tracker.Tracker, fixture.CleanupFunc) {
	db, cleanup := fixture.Bolt(t)

	tr, err := wire.BuildTrackerForTest(db)
	if err != nil {
		t.Fatal(err)
	}

	return tr, cleanup
}
