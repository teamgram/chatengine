/*
 * WARNING! All changes made in this file will be lost!
 * Created from 'scheme.tl' by 'mtprotoc'
 *
 * Copyright 2024 Teamgram Authors.
 *  All rights reserved.
 *
 * Author: teamgramio (teamgram.io@gmail.com)
 */

package service

import (
	"github.com/teamgram/teamgram-server/v2/app/interface/gnetway/gateway"
	"github.com/teamgram/teamgram-server/v2/app/interface/gnetway/internal/svc"
)

type Service struct {
	svcCtx *svc.ServiceContext
	gateway.RPCGatewayServer
}

func New(ctx *svc.ServiceContext, srv gateway.RPCGatewayServer) *Service {
	return &Service{
		svcCtx:           ctx,
		RPCGatewayServer: srv,
	}
}
