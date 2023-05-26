package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	GRPCGATEWAYUSERAGENTHEADER = "grpcgateway-user-agent"
	USERAGENTHEADER            = "user-agent"
	XFORWARDEDFORHEADER        = "x-forwarded-host"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) (data *Metadata) {
	data = &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(GRPCGATEWAYUSERAGENTHEADER); len(userAgents) > 0 {
			data.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(USERAGENTHEADER); len(userAgents) > 0 {
			data.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(XFORWARDEDFORHEADER); len(clientIPs) > 0 {
			data.ClientIP = clientIPs[0]
		}
	}

	if ip, ok := peer.FromContext(ctx); ok {
		data.ClientIP = ip.Addr.String()
	}

	return
}
