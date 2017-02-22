package msg

type PushMsgReq struct {
	push_type  int
	channel_id string
	device_ids []string
	message    string
}

type PushMsgRes struct {
	errno  int
	errmsg string
}

type Message struct {
	Body string
}
