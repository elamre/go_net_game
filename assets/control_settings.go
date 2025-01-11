package assets

import (
	"reflect"
)

type InputId int

type InputState struct {
	IsPressed     [MAX_CONTROL_ID]int
	LastIsPressed [MAX_CONTROL_ID]int
}

func (i *InputState) Changed() bool {
	same := reflect.DeepEqual(i.IsPressed, i.LastIsPressed)
	if !same {
		i.LastIsPressed = i.IsPressed
	}
	return !same
}

const (
	FORWARD_ID     = InputId(iota)
	BACKWARD_ID    = InputId(iota)
	LEFT_ID        = InputId(iota)
	RIGHT_ID       = InputId(iota)
	ACTION_ID      = InputId(iota)
	RELOAD_ID      = InputId(iota)
	NEXT_WEAPON_ID = InputId(iota)
	PREV_WEAPON_ID = InputId(iota)
	MAX_CONTROL_ID = iota
)

var (
	FORWARD_KEY     = KeyW
	BACKWARD_KEY    = KeyS
	LEFT_KEY        = KeyA
	RIGHT_KEY       = KeyD
	ACTION_KEY      = KeyE
	RELOAD_KEY      = KeyR
	NEXT_WEAPON_KEY = KeyDigit2
	PREV_WEAPON_KEY = KeyDigit1
)

type keyCombo struct {
	name string
	id   InputId
	key  Key
}

var keys = []keyCombo{{
	name: "forward",
	id:   FORWARD_ID,
	key:  FORWARD_KEY,
}, {
	name: "backward",
	id:   BACKWARD_ID,
	key:  BACKWARD_KEY,
}, {
	name: "left",
	id:   LEFT_ID,
	key:  LEFT_KEY,
}, {
	name: "right",
	id:   RIGHT_ID,
	key:  RIGHT_KEY,
}, {
	name: "action",
	id:   ACTION_ID,
	key:  ACTION_KEY,
}, {
	name: "reload",
	id:   RELOAD_ID,
	key:  RELOAD_KEY,
}, {
	name: "next",
	id:   NEXT_WEAPON_ID,
	key:  NEXT_WEAPON_KEY,
}, {
	name: "prev",
	id:   PREV_WEAPON_ID,
	key:  PREV_WEAPON_KEY,
},
}

var InputManager InputState
