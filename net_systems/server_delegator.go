package net_systems

import (
	"github.com/elamre/go_helpers/pkg/misc"
	"github.com/elamre/go_helpers/pkg/slice_helpers"
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net/packet_interface"
	. "github.com/elamre/go_net_game/net_systems/common_system"
	"github.com/elamre/logger/pkg/logger"
	"log"
	"reflect"
	"sync"
)

type ServerDelegator struct {
	clientMutex sync.Mutex
	server      net.Server

	netPlayers []*net.ServerPlayer

	subSystems         map[string]ServerSubSystem
	subSystemsCallback map[reflect.Type][]PacketServerCallback

	l *logger.Logger
}

func NewServerDelegator(server net.Server) *ServerDelegator {
	s := &ServerDelegator{
		subSystems:         make(map[string]ServerSubSystem),
		server:             server,
		netPlayers:         make([]*net.ServerPlayer, 0),
		subSystemsCallback: make(map[reflect.Type][]PacketServerCallback),
		l:                  logger.NewLogger(),
	}
	server.AddConnectionCallback(s.clientConnect)
	server.AddDisconnectionCallback(s.clientDisconnect)
	return s
}

func (s *ServerDelegator) clientConnect(player *net.ServerPlayer) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	s.netPlayers = append(s.netPlayers, player)
	log.Printf("Client connected! %+v [%+v]", player, s.netPlayers)
}

func (s *ServerDelegator) clientDisconnect(player *net.ServerPlayer) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	s.l.LogInfof("Client disconnected! %+v [%+v]", player, s.netPlayers)

	result := slice_helpers.RemoveFromListEquals[*net.ServerPlayer](s.netPlayers, func(s *net.ServerPlayer) bool {
		return s.Client == player.Client
	})
	if result == nil {
		s.l.LogWarning("Could not find the client")
	} else {
		s.netPlayers = result
	}

}

func (s *ServerDelegator) Update() {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	for {
		noPackets := true

		for _, p := range s.netPlayers {
			pack := misc.CheckErrorRetVal[packet_interface.Packet](p.ReadPacket())
			if pack != nil {
				noPackets = false
				for _, cb := range s.subSystemsCallback[reflect.TypeOf(pack)] {
					cb(p, s, pack)
				}
			}
		}
		if noPackets {
			break
		}
	}

	for _, sub := range s.subSystems {
		sub.Update(s)
	}
}

func (s *ServerDelegator) RegisterPacketCallback(cb PacketServerCallback, packet packet_interface.Packet) {
	t := reflect.TypeOf(packet)
	if _, ok := s.subSystemsCallback[t]; !ok {
		s.subSystemsCallback[t] = make([]PacketServerCallback, 0)
	}
	s.subSystemsCallback[t] = append(s.subSystemsCallback[t], cb)
}

func (s *ServerDelegator) RemovePacketCallback(cb PacketServerCallback, packetType packet_interface.Packet) {
}
func (s *ServerDelegator) RemoveSubSystem(name string) {
	delete(s.subSystems, name)
}

func (s *ServerDelegator) GetSubsystem(name string) interface{} {
	return s.subSystems[name]
}

func (s *ServerDelegator) RegisterSubSystem(name string, system ServerSubSystem) {
	s.subSystems[name] = system
	system.RegisterCallbacks(s)
}
