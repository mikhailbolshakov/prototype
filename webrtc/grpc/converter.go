package grpc

import (
	"encoding/json"
	"github.com/pion/webrtc/v3"
)

func (s *Server) toSdpDomain(sdp []byte) (*webrtc.SessionDescription, error) {

	var offer *webrtc.SessionDescription
	err := json.Unmarshal(sdp, &offer)
	if err != nil {
		return nil, err
	}
	return offer, nil

}

