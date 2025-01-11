package webrtc

import (
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net/packet_interface"
	"github.com/elamre/go_net_game/net/webrtc/internal/webrtc_server"
	"log"
	"sync"
)

type CallbackFunction func(client *net.ServerPlayer)

func (c CallbackFunction) Equal(other CallbackFunction) bool {
	return true
}

type WebrtcHost struct {
	clientsMutex        sync.Mutex
	server              *webrtc_server.Server
	connectedClientsMap map[*webrtc_server.Connection]*net.ServerPlayer
	connectedClients    []*net.ServerPlayer
	onConnectionCb      []CallbackFunction
	onDisconnectCb      []CallbackFunction
}

func NewWebrtcHost(publicIp string, port int) *WebrtcHost {
	w := WebrtcHost{
		connectedClientsMap: make(map[*webrtc_server.Connection]*net.ServerPlayer),
		connectedClients:    make([]*net.ServerPlayer, 0),
		onConnectionCb:      make([]CallbackFunction, 0),
		onDisconnectCb:      make([]CallbackFunction, 0),
	}
	w.server = webrtc_server.New(webrtc_server.Options{
		MaxConnections: 10,
		HttpPort:       port,
		PublicIP:       publicIp,

		/*		ICEServerURLs:  []string{"stun:127.0.0.1:3478"},*/
		ICEServerURLs: []string{"stun:stun.l.google.com:19302"},
	})
	w.server.OnConnection = w.onConnect
	w.server.OnDisconnect = w.onDisconnect
	return &w
}

func (w *WebrtcHost) ClientIterator(iterator func(c *net.ServerPlayer)) {
	for _, c := range w.connectedClients {
		iterator(c)
	}
}

func (w *WebrtcHost) Start() any {
	w.server.Start()
	return nil
}

func (w *WebrtcHost) Close() any {
	// todo
	return nil
}

func (w *WebrtcHost) AddConnectionCallback(f func(client *net.ServerPlayer)) {
	w.onConnectionCb = append(w.onConnectionCb, f)
}
func (w *WebrtcHost) RemoveConnectionCallback(f func(client *net.ServerPlayer)) {
	// slice_helpers.RemoveFromList[CallbackFunction](f, w.onConnectionCb) // TODO
}
func (w *WebrtcHost) AddDisconnectionCallback(f func(client *net.ServerPlayer)) {
	w.onDisconnectCb = append(w.onDisconnectCb, f)
}
func (w *WebrtcHost) RemoveDisconnectionCallback(f func(client *net.ServerPlayer)) {
	// slice_helpers.RemoveFromList[CallbackFunction](f, w.onDisconnectCb) // TODO
}

func (w *WebrtcHost) onDisconnect(connection *webrtc_server.Connection) {
	w.clientsMutex.Lock()
	defer w.clientsMutex.Unlock()
	log.Printf("onDisconnect: %+v %p", connection, connection)

	for _, f := range w.onDisconnectCb {
		f(w.connectedClientsMap[connection])
	}
	delete(w.connectedClientsMap, connection)
}

func (w *WebrtcHost) onConnect(connection *webrtc_server.Connection) {
	w.clientsMutex.Lock()
	defer w.clientsMutex.Unlock()
	wrappedHostclient := net.NewServerPlayer(NewWebrtcHostClient(connection))
	w.connectedClients = append(w.connectedClients, wrappedHostclient)
	w.connectedClientsMap[connection] = wrappedHostclient

	for _, f := range w.onConnectionCb {
		f(wrappedHostclient)
	}
	log.Printf("onConnect: %+v, %p", connection, connection)

}

func (w *WebrtcHost) BroadcastPacket(packet packet_interface.Packet) {
	for _, c := range w.connectedClients {
		c.WritePacket(packet)
	}
}
