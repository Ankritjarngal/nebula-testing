package types

type RequestType int

const (
	RequestTask RequestType = iota
	SubmitResult
	TaskResponseMsg
)

type Message struct {
	Type    RequestType
	Payload []byte
}

type TaskRequest struct {
	WorkerID string
}

type TaskResponse struct {
	TaskID string
	Input  []byte
}

type ResultSubmission struct {
	TaskID string
	Result []byte
}
