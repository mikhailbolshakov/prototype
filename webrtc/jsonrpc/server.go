package jsonrpc

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	meta "gitlab.medzdrav.ru/prototype/webrtc/meta"
	"net/http"
	"net/url"

	"github.com/sourcegraph/jsonrpc2"
	websocketjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"

	// pprof
	_ "net/http/pprof"
)

// Signal is the grpc/http/websocket signaling server
type Server struct {
	sessionService  domain.SessionsService
	webrtcService   domain.WebrtcService
	roomCoordinator domain.RoomCoordinator
}

// newSignal creates a signaling server
func New(http *kitHttp.Server, sessionService domain.SessionsService, webrtcService domain.WebrtcService) *Server {

	s := &Server{
		sessionService:  sessionService,
		webrtcService:   webrtcService,
	}
	http.SetWsUpgrader(s)

	return s
}

// Set implements WsUpgrader interface
func (s *Server) Set(noAuthRouter *mux.Router, upgrader *websocket.Upgrader) {

	noAuthRouter.HandleFunc("/webrtc", func(w http.ResponseWriter, r *http.Request) {

		sid := r.URL.Query().Get("session")
		room := r.URL.Query().Get("room")

		l := log.L().Cmp("ion").Mth("ws-handler").F(log.FF{"sid": sid, "room": room}).Dbg()

		if sid == "" {
			l.Err("sid empty")
			http.Error(w, "sid empty", http.StatusInternalServerError)
			return
		}

		if room == "" {
			l.Err("room empty")
			http.Error(w, "room empty", http.StatusInternalServerError)
			return
		}

		ctxRq := kitContext.NewRequestCtx().
			WithNewRequestId().
			Webrtc()
		ctx := ctxRq.ToContext(r.Context())

		l.C(ctx)

		session, err := s.sessionService.AuthSession(ctx, sid)
		if err != nil {
			l.E(err).Err("auth session")
			http.Error(w, "auth session error", http.StatusInternalServerError)
			return
		}

		ctx = ctxRq.WithUser(session.UserId, session.Username).WithChatUserId(session.ChatUserId).ToContext(ctx)

		roomMeta, err := s.webrtcService.GetOrCreateRoom(ctx, room)
		if err != nil {
			l.E(err).St().Err()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if roomMeta.NodeId != meta.NodeId {
			endpoint := fmt.Sprintf("%s/webrtc?session=%s&room=%s", roomMeta.Endpoint, sid, roomMeta.Id)
			l.DbgF("redirecting to %s", endpoint)
			u, err := url.Parse(endpoint)
			if err != nil {
				l.E(err).St().Err("parsing backend url to proxy websocket")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			proxy := websocketproxy.NewProxy(u)
			proxy.Upgrader = upgrader
			l.DbgF("starting proxy for room -> node %v @ %v", roomMeta.NodeId, endpoint)
			proxy.ServeHTTP(w, r)
			l.DbgF("closed proxy for session -> node %v @ %v", roomMeta.NodeId, endpoint)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.E(err).St().Err()
			return
		}
		defer c.Close()

		peer := s.webrtcService.NewPeer(ctx)
		signal := newSignal(session.UserId, session.Username, peer, s.webrtcService)

		defer peer.Close(ctx)

		jc := jsonrpc2.NewConn(r.Context(), websocketjsonrpc2.NewObjectStream(c), signal)
		<-jc.DisconnectNotify()

	})
}
