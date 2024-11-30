// Server side api

///////////////////////////////
///////// Nebula server API
///////////////////////////////

---types---
sessionClientEvent server_id:string conn_type:int auth_key_id:long key_type:int perm_auth_key_id:long session_id:long client_ip:string = SessionClientEvent;
sessionClientData  server_id:string conn_type:int auth_key_id:long key_type:int perm_auth_key_id:long session_id:long client_ip:string quick_ack:int salt:long payload:bytes = SessionClientData;
httpSessionData payload:bytes = HttpSessionData;

---functions---
session.queryAuthKey auth_key_id:long = AuthKeyInfo;
session.setAuthKey auth_key:AuthKeyInfo future_salt:FutureSalt expires_in:int = Bool;

session.createSession client:SessionClientEvent = Bool;
session.sendDataToSession data:SessionClientData = Bool;
session.sendHttpDataToSession client:SessionClientData = HttpSessionData;
session.closeSession client:SessionClientEvent = Bool;

session.pushUpdatesData flags:# perm_auth_key_id:long notification:flags.0?true updates:Updates = Bool;
session.pushSessionUpdatesData flags:# perm_auth_key_id:long auth_key_id:long session_id:long updates:Updates = Bool;
session.pushRpcResultData perm_auth_key_id:long auth_key_id:long session_id:long client_req_msg_id:long rpc_result_data:bytes = Bool;

// LAYER 0
