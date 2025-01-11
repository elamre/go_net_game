package game_system_packets

import (
	"bytes"
	"encoding/binary"
)

type PhysicsPacket struct {
	EntityId       int32
	PosX, PosY     float64
	DeltaX, DeltaY float64
}

func NewPhysicsPacket(entityId int32, posX, posY, deltaX, deltaY float64) PhysicsPacket {
	return PhysicsPacket{
		EntityId: entityId,
		PosX:     posX,
		PosY:     posY,
		DeltaX:   deltaX,
		DeltaY:   deltaY,
	}
}

func (e PhysicsPacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, e.EntityId); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.PosX); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.PosY); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.DeltaX); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.DeltaY); err != nil {
		panic(err)
	}

}
func (e PhysicsPacket) FromReader(r *bytes.Reader) any {
	if err := binary.Read(r, binary.LittleEndian, &e.EntityId); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &e.PosX); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &e.PosY); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &e.DeltaX); err != nil {
		panic(err)
	}

	if err := binary.Read(r, binary.LittleEndian, &e.DeltaY); err != nil {
		panic(err)
	}

	return e
}

func (e PhysicsPacket) AckRequired() bool {
	return false
}
