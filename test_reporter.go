package reporter

import (
	"fmt"
	"net"
	"testing"

	"github.com/xtech-cloud/omo-mod-reporter/processor"
	"github.com/xtech-cloud/omo-mod-reporter/protocol"
)

func handleReporterPing(_req *protocol.Request, _rsp *protocol.Response, _sender interface{}) {
	_rsp.Head.Msg = "pong"
	_rsp.Body = &protocol.EmptyBlock{}
	sender := _sender.(*net.UDPAddr)
	fmt.Println(sender.String())
}

func handleWorkerPing(_req *protocol.Request, _rsp *protocol.Response, _sender interface{}) {
	_rsp.Head.Msg = "pong"
	_rsp.Body = &protocol.EmptyBlock{}
}

func Test_RunReporter(_t *testing.T) {
	reporterProcessor := processor.NewProcessor()
	reporterProcessor.BindJsonHandler("ping", handleReporterPing)

	reporter, _ := NewReporter()
	reporter.Run(":18999", reporterProcessor.ProcessJson)
}
