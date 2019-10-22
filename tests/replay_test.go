package test

import (
	"github.com/SRI-CSL/gllvm/shared"
	"testing"
	"fmt"
)


func  Test_replay_thttpd(t *testing.T) {
	arg := "../data/thttpd_replay.log"

	ok := shared.Replay(arg)

	if !ok {
		t.Errorf("Replay of %v returned %v\n", arg, ok)
	} else {
		fmt.Println("Replay OK")
	}
}

func  Test_replay_nodejs(t *testing.T) {
	arg := "../data/nodejs_replay.log"

	ok := shared.Replay(arg)

	if !ok {
		t.Errorf("Replay of %v returned %v\n", arg, ok)
	} else {
		fmt.Println("Replay OK")
	}
}
