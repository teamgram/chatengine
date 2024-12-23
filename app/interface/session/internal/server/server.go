// Copyright 2024 Teamgram Authors
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: teamgramio (teamgram.io@gmail.com)
//

package server

import (
	"flag"

	"github.com/teamgram/proto/v2/rpc/codec"
	"github.com/teamgram/teamgram-server/v2/app/interface/session/internal/config"
	"github.com/teamgram/teamgram-server/v2/app/interface/session/internal/server/tg/service"
	"github.com/teamgram/teamgram-server/v2/app/interface/session/internal/svc"
	"github.com/teamgram/teamgram-server/v2/app/interface/session/session/sessionservice"

	"github.com/cloudwego/kitex/server"
)

var configFile = flag.String("f", "etc/session.yaml", "the config file")

type Server struct {
	server.Server
}

func New() *Server {
	return new(Server)
}

func (s *Server) Initialize() error {
	var c config.Config
	ctx := svc.NewServiceContext(c)
	_ = ctx

	cCodec := codec.NewZRpcCodec(true)
	s.Server = sessionservice.NewServer(service.New(ctx), server.WithCodec(cCodec))
	return nil
}

func (s *Server) RunLoop() {
	if err := s.Server.Run(); err != nil {
		// log.Println("server stopped with error:", err)
	} else {
		// log.Println("server stopped")
	}
}

func (s *Server) Destroy() {
}
