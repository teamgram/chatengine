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
	"context"

	"github.com/teamgram/proto/mtproto"
	"github.com/teamgram/teamgram-server/v2/app/interface/session/internal/core"
	"github.com/teamgram/teamgram-server/v2/app/interface/session/session"
)

// SessionQueryAuthKey
// session.queryAuthKey auth_key_id:long = AuthKeyInfo;
func (s *Service) SessionQueryAuthKey(ctx context.Context, request *session.TLSessionQueryAuthKey) (*mtproto.AuthKeyInfo, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.queryAuthKey - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionQueryAuthKey(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.queryAuthKey - reply: {%s}", r)
	return r, err
}

// SessionSetAuthKey
// session.setAuthKey auth_key:AuthKeyInfo future_salt:FutureSalt expires_in:int = Bool;
func (s *Service) SessionSetAuthKey(ctx context.Context, request *session.TLSessionSetAuthKey) (*mtproto.Bool, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.setAuthKey - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionSetAuthKey(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.setAuthKey - reply: {%s}", r)
	return r, err
}

// SessionCreateSession
// session.createSession client:SessionClientEvent = Bool;
func (s *Service) SessionCreateSession(ctx context.Context, request *session.TLSessionCreateSession) (*mtproto.Bool, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.createSession - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionCreateSession(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.createSession - reply: {%s}", r)
	return r, err
}

// SessionSendDataToSession
// session.sendDataToSession data:SessionClientData = Bool;
func (s *Service) SessionSendDataToSession(ctx context.Context, request *session.TLSessionSendDataToSession) (*mtproto.Bool, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.sendDataToSession - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionSendDataToSession(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.sendDataToSession - reply: {%s}", r)
	return r, err
}

// SessionSendHttpDataToSession
// session.sendHttpDataToSession client:SessionClientData = HttpSessionData;
func (s *Service) SessionSendHttpDataToSession(ctx context.Context, request *session.TLSessionSendHttpDataToSession) (*session.HttpSessionData, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.sendHttpDataToSession - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionSendHttpDataToSession(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.sendHttpDataToSession - reply: {%s}", r)
	return r, err
}

// SessionCloseSession
// session.closeSession client:SessionClientEvent = Bool;
func (s *Service) SessionCloseSession(ctx context.Context, request *session.TLSessionCloseSession) (*mtproto.Bool, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.closeSession - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionCloseSession(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.closeSession - reply: {%s}", r)
	return r, err
}

// SessionPushUpdatesData
// session.pushUpdatesData flags:# perm_auth_key_id:long notification:flags.0?true updates:Updates = Bool;
func (s *Service) SessionPushUpdatesData(ctx context.Context, request *session.TLSessionPushUpdatesData) (*mtproto.Bool, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.pushUpdatesData - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionPushUpdatesData(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.pushUpdatesData - reply: {%s}", r)
	return r, err
}

// SessionPushSessionUpdatesData
// session.pushSessionUpdatesData flags:# perm_auth_key_id:long auth_key_id:long session_id:long updates:Updates = Bool;
func (s *Service) SessionPushSessionUpdatesData(ctx context.Context, request *session.TLSessionPushSessionUpdatesData) (*mtproto.Bool, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.pushSessionUpdatesData - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionPushSessionUpdatesData(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.pushSessionUpdatesData - reply: {%s}", r)
	return r, err
}

// SessionPushRpcResultData
// session.pushRpcResultData perm_auth_key_id:long auth_key_id:long session_id:long client_req_msg_id:long rpc_result_data:bytes = Bool;
func (s *Service) SessionPushRpcResultData(ctx context.Context, request *session.TLSessionPushRpcResultData) (*mtproto.Bool, error) {
	c := core.New(ctx, s.svcCtx)
	c.Logger.Debugf("session.pushRpcResultData - metadata: {%s}, request: {%s}", c.MD, request)

	r, err := c.SessionPushRpcResultData(request)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("session.pushRpcResultData - reply: {%s}", r)
	return r, err
}
