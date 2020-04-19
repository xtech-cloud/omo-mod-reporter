package processor

import (
	"encoding/json"

	"omo-mod-reporter/protocol"
)

type Processor struct {
	jsonHandlers map[string]func(*protocol.Request, *protocol.Response, interface{})
}

func NewProcessor() *Processor {
	return &Processor{
		jsonHandlers: make(map[string]func(*protocol.Request, *protocol.Response, interface{}), 0),
	}
}

func (this *Processor) ProcessJson(_bytes []byte, _sender interface{}) ([]byte, error) {
	req := &protocol.Request{}
	rsp := &protocol.Response{}

	err := json.Unmarshal(_bytes, req)
	if nil != err {
		rsp.Head.ErrCode = -2
		rsp.Head.ErrString = err.Error()
		rsp.Body = &protocol.EmptyBlock{}
		return jsonToBytes(rsp)
	}

	rsp.Head.Msg = req.Head.Msg
	rsp.Head.Session = req.Head.Session

	if _, ok := this.jsonHandlers[req.Head.Msg]; !ok {
		rsp.Head.ErrCode = -1
		rsp.Head.ErrString = "handler not found"
		rsp.Body = &protocol.EmptyBlock{}
		return jsonToBytes(rsp)
	}

	this.jsonHandlers[req.Head.Msg](req, rsp, _sender)
	return jsonToBytes(rsp)
}

func (this *Processor) BindJsonHandler(_msg string, _handler func(*protocol.Request, *protocol.Response, interface{})) {
	this.jsonHandlers[_msg] = _handler
}

func jsonToBytes(_json *protocol.Response) ([]byte, error) {
	return json.Marshal(_json)
}
