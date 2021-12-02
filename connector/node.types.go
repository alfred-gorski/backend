package connector

type Node struct {
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	Status        string `json:"status" bson:"status"`
	WorkerID      string `json:"worker_id" bson:"worker_id"`
	JobID         string `json:"job_id" bson:"job_id"`
	FriendlyNames struct {
		De string `json:"de" bson:"de"`
		En string `json:"en" bson:"en"`
	} `json:"friendly_names" bson:"friendly_names"`
	VoteParticipant bool   `json:"vote_participant" bson:"vote_participant"`
	Schema          string `json:"schema" bson:"schema,omitempty"`
}

type Message struct {
	LamportTs int    `json:"lamport_ts" bson:"lamport_ts"`
	WorkerID  string `json:"worker_id" bson:"worker_id"`
	MsgType   string `json:"msg_type" bson:"msg_type"`
	Content   Node   `json:"content" bson:"content"`
	OrderID   string `json:"order_id" bson:"order_id"`
}
