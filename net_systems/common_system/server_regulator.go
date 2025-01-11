package common_system

import (
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net/packet_interface"
)

type PacketServerCallback = func(c *net.ServerPlayer, d ServerRegulator, pack packet_interface.Packet)
type PacketClientCallback = func(c net.Client, d ClientRegulator, pack packet_interface.Packet)

type ServerSubSystem interface {
	RegisterCallbacks(r ServerRegulator)
	Update(r ServerRegulator)
}

type ClientSubSystem interface {
	RegisterCallbacks(r ClientRegulator)
	Update()
}

type ServerRegulator interface {
	RemoveSubSystem(name string)
	GetSubsystem(name string) interface{}
	RegisterSubSystem(name string, system ServerSubSystem)

	RegisterPacketCallback(cb PacketServerCallback, packet packet_interface.Packet)
	RemovePacketCallback(cb PacketServerCallback, packetType packet_interface.Packet)
}

type ClientRegulator interface {
	RemoveSubSystem(name string)
	GetSubsystem(name string) interface{}
	RegisterSubSystem(name string, system ClientSubSystem)

	RegisterPacketCallback(cb PacketClientCallback, packet packet_interface.Packet)
	RemovePacketCallback(cb PacketClientCallback, packetType packet_interface.Packet)
}
