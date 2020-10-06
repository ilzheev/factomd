package simtest

import (
	"strings"
	"testing"

	"github.com/FactomProject/factomd/testHelper/simulation"

	"github.com/FactomProject/factomd/fnode"

	. "github.com/FactomProject/factomd/testHelper"
)

func TestFilterAPIInput(t *testing.T) {

	state0 := simulation.SetupSim("LLLAF", map[string]string{}, 25, 1, 1, t)

	simulation.RunCmd("1")
	simulation.RunCmd("w")
	simulation.RunCmd("s")

	apiRegex := "EOM.*5.*minute +1" // It has two spaces.
	SetInputFilter(apiRegex)

	simulation.WaitBlocks(state0, 5)

	// The message-filter call we did above should have caused an election and SO, Node01 should not be a leader anymore.
	if fnode.Get(1).State.Leader {
		t.Fatalf("Node01 should not be leader!")
	}

	simulation.CheckAuthoritySet(t)

	// Check Node01 Network Input logs to make sure there are no enqueued including our Regex
	out := SystemCall(`grep "enqueue" fnode01_networkinputs.txt | grep "` + apiRegex + `" | grep -v "EmbeddedMsg" | wc -l`)

	if strings.TrimSuffix(strings.Trim(string(out), " "), "\n") != string("0") {
		t.Fatalf("Filter missed let a message pass 1.")
	}

	// Check Node01 Network Input logs to make sure there are no Dropped messaged besides the ones for our Regex
	out2 := SystemCall(`grep "Drop, matched filter Regex" fnode01_networkinputs.txt | grep -Ev "` + apiRegex + `" | wc -l`)

	if strings.TrimSuffix(strings.Trim(string(out2), " "), "\n") != string("0") {
		t.Fatalf("Filter missed let a message pass 2.")
	}

	simulation.ShutDownEverything(t)
}
