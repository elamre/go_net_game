package net

import (
	"fmt"
	"github.com/elamre/go_net_game/net/packet_interface"
	"log"
)

func NewServerPlayer(c ServerClient) *ServerPlayer {
	log.Printf("New server client for: %p", c)
	return &ServerPlayer{
		Client:        c,
		DataSpecifics: make(map[interface{}]interface{}),
	}
}

func (s *ServerPlayer) WritePacket(pack packet_interface.Packet) any {
	return s.Client.WritePacket(pack)
}

func (s *ServerPlayer) ReadPacket() (packet packet_interface.Packet, err error) {
	return s.Client.ReadPacket()
}

func (s ServerPlayer) String() string {
	return fmt.Sprintf("%+v [%+v]", s.Client, s.DataSpecifics)
}

type ServerPlayer struct {
	NetPlayer     NetPlayer
	Client        ServerClient
	DataSpecifics map[interface{}]interface{}
}

type ServerOptions struct {
	Port           int
	MaxConnections int
}

type ClientOptions struct {
	Target string
	Port   int
}

type Server interface {
	Start() any
	Close() any
	AddConnectionCallback(func(player *ServerPlayer))
	RemoveConnectionCallback(func(player *ServerPlayer))
	AddDisconnectionCallback(func(player *ServerPlayer))
	RemoveDisconnectionCallback(func(player *ServerPlayer))
	BroadcastPacket(packet packet_interface.Packet)
	ClientIterator(iterator func(c *ServerPlayer))
}

type ServerClient interface {
	WritePacket(packet packet_interface.Packet) any
	GotPacket() bool
	ReadPacket() (packet packet_interface.Packet, err error)
	Close() any
}

type Client interface {
	Connect() any
	IsConnected() bool
	WritePacket(packet packet_interface.Packet) any
	GotPacket() bool
	ReadPacket() (packet packet_interface.Packet, err error)
	Close() any
}
