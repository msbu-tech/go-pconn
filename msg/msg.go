package msg

type PushMsgReq struct {
	Push_type  int
	Channel_id string
	Device_ids []string
	Content    string
}

type PushMsgRes struct {
	Errno  int
	Errmsg string
}

type Message struct {
	Body string
}
