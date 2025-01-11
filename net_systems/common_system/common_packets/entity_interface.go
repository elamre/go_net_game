package common_packets

import (
	"bytes"
	"github.com/elamre/go_net_game/net/packet_interface"
	"github.com/elamre/go_net_game/net_systems/game_system_packets"
	"math"
)

type MinimumEntityInterface interface {
	GetOwner() int32
	GetId() int32
	GetType() int32
}

type NetworkableEntity interface {
	MinimumEntityInterface
	AddPackets(packets *[]packet_interface.Packet)
	WriteSpecifics(w *bytes.Buffer)
	ReadSpecifics(r *bytes.Reader)
}

type UpdatableEntity interface {
	MinimumEntityInterface
	Update()
	ClientSimulatePhysics()
	UpdatePhysics(packet game_system_packets.PhysicsPacket)
}

type DefaultMinimumEntity struct {
	ClientCreated bool
	Id            int32
	Owner         int32
	Type          int32
}

func (d *DefaultMinimumEntity) GetOwner() int32 {
	return d.Owner
}

func (d *DefaultMinimumEntity) GetId() int32 {
	return d.Id
}

func (d *DefaultMinimumEntity) GetType() int32 {
	return d.Type
}

type DefaultUpdatableEntity struct {
	DefaultMinimumEntity
	X, Y           float64
	DeltaX, DeltaY float64

	targetPosX, targetPosY float64

	previousX, previousY           float64
	previousDeltaX, previousDeltaY float64

	Changed       bool
	PacketsToSend []packet_interface.Packet

	UpdateClientSpecifics func()
}

func (d *DefaultUpdatableEntity) AddPackets(packets *[]packet_interface.Packet) {
	if d.Changed {
		*packets = append(*packets, game_system_packets.NewPhysicsPacket(d.GetId(), d.X, d.Y, d.DeltaX, d.DeltaY))
		/*		d.Changed = false
		 */
	}
}

func (d *DefaultUpdatableEntity) ClientSimulatePhysics() {
	if d.UpdateClientSpecifics != nil {
		d.UpdateClientSpecifics()
	}
	if !d.ClientCreated {
		if math.Abs(math.Floor(d.Y)-math.Floor(d.targetPosY)) > 2 || math.Abs(math.Floor(d.X)-math.Floor(d.targetPosX)) > 2 {
			d.X += d.DeltaX
			d.Y += d.DeltaY
		}
	}
}

func (d *DefaultUpdatableEntity) Update() {
	d.X += d.DeltaX
	d.Y += d.DeltaY
	if d.X != d.previousX || d.Y != d.previousY || d.DeltaY != d.previousDeltaY || d.DeltaX != d.previousDeltaX {
		d.Changed = true
	}
	d.previousX = d.X
	d.previousY = d.Y
	d.previousDeltaX = d.DeltaX
	d.previousDeltaY = d.DeltaY
}

func (d *DefaultUpdatableEntity) UpdatePhysics(packet game_system_packets.PhysicsPacket) {
	deltaX := packet.PosX - d.X
	deltaY := packet.PosY - d.Y
	if math.Abs(deltaX) > 30 || math.Abs(deltaY) > 30 {
		d.X = packet.PosX
		d.Y = packet.PosY
		d.DeltaX = packet.DeltaX
		d.DeltaY = packet.DeltaY
	} else if !d.ClientCreated {
		angle := math.Atan2(deltaY, deltaX)
		d.DeltaX = math.Cos(angle)
		d.DeltaY = math.Sin(angle)
	}
	d.targetPosX = packet.PosX
	d.targetPosY = packet.PosY
}
