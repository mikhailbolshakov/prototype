package ion

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/config"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"

	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/sourcegraph/jsonrpc2"
	websocketjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"

	// pprof
	_ "net/http/pprof"
)

// Signal is the grpc/http/websocket signaling server
type Signal struct {
	c       coordinator
	errChan chan error
	config  *config.SignalConfig
}

// newSignal creates a signaling server
func newSignal(s *sfu.SFU, c coordinator, conf *config.SignalConfig) (*Signal, chan error) {
	e := make(chan error)
	w := &Signal{
		c:       c,
		errChan: e,
		config:  conf,
	}
	return w, e
}

func (s *Signal) wsHandler(w http.ResponseWriter, r *http.Request) {

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

	// TODO: centralized auth
	// we need to get a session as a url parameter
	// then we check session in auth service and get session attrs (userid .. etc)
	ctx := kitContext.NewRequestCtx().
		WithNewRequestId().
		Webrtc().
		ToContext(r.Context())

	l.C(ctx)

	meta, err := s.c.getOrCreateRoom(ctx, room)
	if err != nil {
		l.E(err).St().Err()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	if meta.Redirect {
		endpoint := fmt.Sprintf("%s/webrtc?session=%s&room=%s", meta.NodeEndpoint, sid, meta.RoomID)
		l.DbgF("redirecting to %s", endpoint)
		u, err := url.Parse(endpoint)
		if err != nil {
			l.E(err).St().Err("parsing backend url to proxy websocket")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		proxy := websocketproxy.NewProxy(u)
		proxy.Upgrader = &upgrader
		l.DbgF("starting proxy for room -> node %v @ %v", meta.NodeID, endpoint)
		prometheusGaugeProxyClients.Inc()
		proxy.ServeHTTP(w, r)
		prometheusGaugeProxyClients.Dec()
		l.DbgF("closed proxy for session -> node %v @ %v", meta.NodeID, endpoint)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.E(err).St().Err()
		return
	}
	defer c.Close()

	prometheusGaugeClients.Inc()
	p := JSONSignal{
		sync.Mutex{},
		s.c,
		sfu.NewPeer(s.c),
	}
	defer p.Close()

	jc := jsonrpc2.NewConn(r.Context(), websocketjsonrpc2.NewObjectStream(c), &p)
	<-jc.DisconnectNotify()
	prometheusGaugeClients.Dec()
}

// ServeWebsocket listens for incoming websocket signaling requests
func (s *Signal) ServeWebsocket() {

	r := mux.NewRouter()
	r.Handle("/webrtc", http.HandlerFunc(s.wsHandler))
	r.Handle("/metrics", metricsHandler())
	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	http.Handle("/", r)

	var err error
	l := log.L().Pr("ws").Cmp("ion")
	if s.config.Key != "" && s.config.Cert != "" {
		l.InfF("listening at https://%s", s.config.HTTPAddr)
		err = http.ListenAndServeTLS(s.config.HTTPAddr, s.config.Cert, s.config.Key, nil)
	} else {
		l.InfF("listening at http://%s", s.config.HTTPAddr)
		err = http.ListenAndServe(s.config.HTTPAddr, nil)
	}

	if err != nil {
		s.errChan <- err
	}
}

