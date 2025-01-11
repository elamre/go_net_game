package game_system

import (
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net_systems"
	"github.com/elamre/go_net_game/net_systems/lobby_system"
	"github.com/elamre/go_net_game/net_systems/ping_system"
)

const (
	ClientLobbyTag = "clientLobby"
	ClientPingTag  = "clientPing"
	ClientGameTag  = "clientGame"
)

func CreateClientSystem(client net.Client) *net_systems.ClientDelegator {
	clientDelegator := net_systems.NewClientDelegator(client)
	clientLobby := lobby_system.NewLobbyClientSystem(client)
	clientDelegator.RegisterSubSystem(ClientLobbyTag, clientLobby)
	clientDelegator.RegisterSubSystem(ClientPingTag, ping_system.NewPingClientSystem(client))
	return clientDelegator
}
