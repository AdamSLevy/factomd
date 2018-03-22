// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package state

import (
	"fmt"
	"time"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
	log "github.com/sirupsen/logrus"
)

func (state *State) ValidatorLoop() {
	timeStruct := new(Timer)

	startup := false
	sendMsg := func(msg interfaces.IMsg) {
		if msg != nil {
			if _, ok := msg.(*messages.Ack); ok {
				state.LogMessage("ackQueue", "enqueue", msg)
				state.ackQueue <- msg //
			} else {
				state.LogMessage("msgQueue", "enqueue", msg)
				state.msgQueue <- msg //
			}
		}
	}

	go func() {
		for !startup {
			time.Sleep(10 * time.Millisecond)
		}
		for {
			min := <-state.tickerQueue
			timeStruct.timer(state, min)
		}
	}()
	go func() {
		for !startup {
			time.Sleep(10 * time.Millisecond)
		}
		for {
			msg := <-state.TimerMsgQueue()
			sendMsg(msg)
		}
	}()
	go func() {
		for !startup {
			time.Sleep(10 * time.Millisecond)
		}
		for {
			msg := state.apiQueue.BlockingDequeue()
			sendMsg(msg)
		}
	}()
	go func() {
		for !startup {
			time.Sleep(10 * time.Millisecond)
		}
		for {
			msg := state.InMsgQueue().BlockingDequeue()
			sendMsg(msg)
		}
	}()

	// Sort the messages.

	for {
		// Check if we should shut down.
		select {
		case <-state.ShutdownChan:
			fmt.Println("Closing the Database on", state.GetFactomNodeName())
			state.DB.Close()
			state.StateSaverStruct.StopSaving()
			fmt.Println(state.GetFactomNodeName(), "closed")
			state.IsRunning = false
			return
		default:
		}

		// Look for pending messages, and get one if there is one.
		// Process any messages we might have queued up.
		for i := 0; i < 50; i++ {
			p, b := state.Process(), state.UpdateState()
			if !p && !b {
				time.Sleep(10 * time.Millisecond)
				break
			}
			//fmt.Printf("dddd %20s %10s --- %10s %10v %10s %10v\n", "Validation", state.FactomNodeName, "Process", p, "Update", b)
		}

		startup = true

	}
}

type Timer struct {
	lastMin      int
	lastDBHeight uint32
}

func (t *Timer) timer(state *State, min int) {
	t.lastMin = min

	eom := new(messages.EOM)
	eom.Timestamp = state.GetTimestamp()
	eom.ChainID = state.GetIdentityChainID()
	eom.Sign(state)
	eom.SetLocal(true)
	consenLogger.WithFields(log.Fields{"func": "GenerateEOM", "lheight": state.GetLeaderHeight()}).WithFields(eom.LogFields()).Debug("Generate EOM")

	if state.RunLeader { // don't generate EOM if we are not a leader or are loading the DBState messages
		state.TimerMsgQueue() <- eom
	}
}
