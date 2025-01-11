package game_system

import (
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net_systems"
	"github.com/elamre/go_net_game/net_systems/common_system"
	"github.com/elamre/go_net_game/net_systems/ping_system"
)

const (
	ServerLobbyTag = "serverLobby"
	ServerUsersTag = "serverUsers"
	ServerPingTag  = "serverPing"
	ServerGameTag  = "serverGame"
)

func CreateServerSystem(server net.Server) *net_systems.ServerDelegator {
	serverDelegator := net_systems.NewServerDelegator(server)
	serverDelegator.RegisterSubSystem(ServerPingTag, ping_system.NewPingServerSystem())
	userManagement := common_system.NewUserManagement(server)
	serverDelegator.RegisterSubSystem(ServerUsersTag, userManagement)
	return serverDelegator
}
