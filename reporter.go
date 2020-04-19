package reporter

import (
	"net"
)

type Reporter struct {
	Conn    *net.UDPConn
	inChan  chan *Report
	outChan chan *Report
	buff    []byte
}

type Report struct {
	Conn      *net.UDPConn
	Addr      *net.UDPAddr
	Processor func([]byte, interface{}) ([]byte, error)
	Content   []byte //接收的内容
	Comment   []byte //发送的内容
}

func NewReporter() (*Reporter, error) {

	reporter := &Reporter{
		Conn:    nil,
		inChan:  make(chan *Report),
		outChan: make(chan *Report),
		buff:    make([]byte, 1024*1024),
	}
	return reporter, nil
}

func (this *Reporter) Run(_proc string, _processor func([]byte, interface{}) ([]byte, error)) {
	udp_addr, err := net.ResolveUDPAddr("udp", _proc)

	if nil != err {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", udp_addr)
	if nil != err {
		panic(err)
	}
	this.Conn = conn
	//异步处理
	go this.asyncProcess()
	//异步回复
	go this.asyncReply()

	for {
		rlen, addr, err := this.Conn.ReadFromUDP(this.buff)
		if err != nil {
			continue
		}
		if rlen > 0 {
			data := make([]byte, rlen)
			copy(data, this.buff[:rlen])
			report := &Report{
				Conn:      this.Conn,
				Addr:      addr,
				Content:   data,
				Comment:   make([]byte, 0),
				Processor: _processor,
			}
			//将报告放入读取管道
			this.inChan <- report
		}
	}
}

func (this *Reporter) asyncProcess() {
	for {
		select {
		//从读取管道中取出一个报告用于处理
		case report := <-this.inChan:
			comment, err := report.Processor(report.Content, report.Addr)
			//出现异常时，丢弃报告，不回复
			if nil == err {
				//将处理完的报告放入发送管道
				report.Comment = comment
				this.outChan <- report
			}
		}
	}
}

func (this *Reporter) asyncReply() {
	for {
		select {
		//从发送管道取出一个报告用于回复
		case report := <-this.outChan:
			report.Conn.WriteToUDP(report.Comment, report.Addr)
		}
	}
}
