package game_system_packets

import (
	"bytes"
	"encoding/binary"
)

const (
	EntityCreated   = iota
	EntityDestroyed = iota
	EntityRequest   = iota
)

type EntityStatePacket struct {
	OwnerId     int32
	EntityId    int32
	EntityState int32
	EntityType  int32
}

func NewEntityStatePacket(ownerId, entityId, typeId, entityState int32) EntityStatePacket {
	return EntityStatePacket{
		OwnerId:     int32(ownerId),
		EntityId:    int32(entityId),
		EntityState: int32(entityState),
		EntityType:  int32(typeId),
	}

}

func (e EntityStatePacket) ToWriter(w *bytes.Buffer) {
	if err := binary.Write(w, binary.LittleEndian, e.OwnerId); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.EntityId); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.EntityState); err != nil {
		panic(err)
	}
	if err := binary.Write(w, binary.LittleEndian, e.EntityType); err != nil {
		panic(err)
	}

}
func (e EntityStatePacket) FromReader(r *bytes.Reader) any {
	if err := binary.Read(r, binary.LittleEndian, &e.OwnerId); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &e.EntityId); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &e.EntityState); err != nil {
		panic(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &e.EntityType); err != nil {
		panic(err)
	}
	return e
}

func (e EntityStatePacket) AckRequired() bool {
	return true
}
