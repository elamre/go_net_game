package main

import (
	"flag"
	"github.com/elamre/go_net_game/game_system"
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net/packet_interface"
	"github.com/elamre/go_net_game/net/webrtc"
	"github.com/elamre/go_net_game/net_systems/common_system/common_packets"
	"github.com/elamre/go_net_game/net_systems/game_system_packets"
	"github.com/elamre/go_net_game/net_systems/lobby_system/lobby_system_packets"
	"github.com/elamre/go_net_game/net_systems/ping_system/ping_packets"
	"github.com/elamre/logger/pkg/logger"
	"log"
	"sync"
)

var (
	argPort = flag.Int("port", 50051, "network port")
	argIp   = flag.String("ip", "0.0.0.0", "ip to connect to, 0.0.0.0 to host")
	argName = flag.String("name", "elmar", "set player name")
)

var once sync.Once
var settingsLogger = logger.GetSettings().GetLogger("settings")

func init() {
	once.Do(func() {
		// Just to make sure, in came init gets called accidentally
		common_packets.Register()
		lobby_system_packets.Register()
		ping_packets.Register()
		game_system_packets.Register()
		settingsLogger.LogInfof("Packages registered: %+v", packet_interface.GetRegisteredPackets())
	})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	if *argIp == "0.0.0.0" {
		var serverObject net.Server = webrtc.NewWebrtcHost(*argIp, *argPort)
		serverDelegator := game_system.CreateServerSystem(serverObject)

		_ = serverDelegator
		serverObject.Start()
		for {
			serverDelegator.Update()
		}
	}
}
