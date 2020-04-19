package protocol

import "encoding/json"

type EmptyBlock struct {
}

type RequestHeadBlock struct {
	Msg     string `json:"msg"`
	Session string `json:"session"`
}

type ResponseHeadBlock struct {
	Msg       string `json:"msg"`
	Session   string `json:"session"`
	ErrCode   int    `json:"errcode"`
	ErrString string `json:"errstring"`
}

type NotifyHeadBlock struct {
	Msg string `json:"msg"`
}

type Request struct {
	Head RequestHeadBlock `json:"head"`
	Body interface{}      `json:"body"`
}

type Response struct {
	Head ResponseHeadBlock `json:"head"`
	Body interface{}       `json:"body"`
}

type Notify struct {
	Head NotifyHeadBlock `json:"head"`
	Body interface{}     `json:"body"`
}

func ToJSON(_data interface{}) ([]byte, error) {
	return json.Marshal(_data)
}

func FromJSON(_data []byte, _target interface{}) error {
	return json.Unmarshal(_data, _target)
}
