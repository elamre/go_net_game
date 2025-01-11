package game_system_packets

import (
	"bytes"
	"encoding/binary"
	"github.com/elamre/go_net_game/assets"
)

type InputStatePacket struct {
	encodedInputs     uint16
	decodedInputState assets.InputState
}

func NewInputStatePacket(inputState *assets.InputState) InputStatePacket {
	inp := InputStatePacket{}
	for i := 0; i < assets.MAX_CONTROL_ID; i++ {
		inp.encodedInputs |= uint16(inputState.IsPressed[i]) << i
	}
	return inp
}

func (e InputStatePacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, e.encodedInputs); err != nil {
		panic(err)
	}

}
func (e InputStatePacket) FromReader(r *bytes.Reader) any {
	if err := binary.Read(r, binary.LittleEndian, &e.encodedInputs); err != nil {
		panic(err)
	}
	return e
}

func (e InputStatePacket) GetDecodedInputState() assets.InputState {
	for i := 0; i < assets.MAX_CONTROL_ID; i++ {
		e.decodedInputState.IsPressed[i] = int((e.encodedInputs >> i) & 0b1)
	}
	return e.decodedInputState
}

func (e InputStatePacket) AckRequired() bool {
	return true
}
