package common_packets

import (
	"bytes"
	"encoding/binary"
)

type EntitySpecificsPacket struct {
	Entity NetworkableEntity
	ID     int32
}

func (e *EntitySpecificsPacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, e.ID); err != nil {
		panic(err)
	}
	e.Entity.WriteSpecifics(w)
}

func (e *EntitySpecificsPacket) FromReader(r *bytes.Reader) any {
	if err := binary.Read(r, binary.LittleEndian, &e.ID); err != nil {
		panic(err)
	}
	return r
}

func (e *EntitySpecificsPacket) AckRequired() bool {
	return true
}
