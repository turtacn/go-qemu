package migrate

import (
	"fmt"
	"time"
)

type State int

const (
	StateIdle State = iota
	StatePreparing
	StateNegotiating
	StateMigrating
	StateCompleting
	StateDone
	StateFailed
	StateRollingBack
	StateRollbackDone
)

type StateTransition struct {
	From      State
	To        State
	Reason    string
	Timestamp time.Time
}

type StateMachine struct {
	current State
	vmid    int
	history []StateTransition
}

func NewStateMachine(vmid int) *StateMachine {
	return &StateMachine{
		current: StateIdle,
		vmid:    vmid,
	}
}

func (sm *StateMachine) Current() State {
	return sm.current
}

func (sm *StateMachine) Transition(to State, reason string) error {
	if !sm.isValidTransition(sm.current, to) {
		return fmt.Errorf("invalid state transition: %v -> %v", sm.current, to)
	}
	sm.history = append(sm.history, StateTransition{
		From:      sm.current,
		To:        to,
		Reason:    reason,
		Timestamp: time.Now(),
	})
	sm.current = to
	return nil
}

func (sm *StateMachine) isValidTransition(from, to State) bool {
	allowed := map[State][]State{
		StateIdle:        {StatePreparing},
		StatePreparing:   {StateNegotiating, StateFailed},
		StateNegotiating: {StateMigrating, StateFailed},
		StateMigrating:   {StateCompleting, StateFailed},
		StateCompleting:  {StateDone, StateFailed},
		StateFailed:      {StateRollingBack},
		StateRollingBack: {StateRollbackDone},
	}
	for _, a := range allowed[from] {
		if a == to {
			return true
		}
	}
	return false
}
