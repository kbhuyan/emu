package emu

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

type messageImpl struct {
	Name    string
	Attribs map[string]interface{}
}

func (m *messageImpl) GetName() string {
	return m.Name
}

func (m *messageImpl) SetAttrib(key string, value interface{}) {
	m.Attribs[key] = value
}

func (m *messageImpl) GetAttrib(key string) (interface{}, bool) {
	value, ok := m.Attribs[key]
	return value, ok
}

type emuImpl struct {
	conn io.ReadWriteCloser
	//	config    EmuConfig
	responses chan Message
	ctx       context.Context
	cancel    context.CancelFunc
	cmdState  *commandState
	opt       *EmuOptions
	sendMsg   bool
}

func newEmuImpl(dev string, opt *EmuOptions) (Emu, error) {
	initLog(opt.LogWriter)
	// Configure serial port
	mode := &serial.Mode{
		BaudRate: opt.BaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(dev, mode)
	if err != nil {
		return nil, fmt.Errorf("serial open failed: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &emuImpl{
		conn: port,
		//		config:    cfg,
		responses: make(chan Message, 10),
		ctx:       ctx,
		cancel:    cancel,
		cmdState:  nil,
		opt:       opt,
		sendMsg:   false,
	}, nil
}

func (e *emuImpl) Start() {
	go e.reader()
}

func (e *emuImpl) SendCommand(command Command) error {
	time.Sleep(100 * time.Millisecond)
	cmd := newCommandState()
	cmd.command = command
	cmd.Status = CmdPending
	value, exists := cmdRspMap[command.GetName()]
	if !exists {
		return fmt.Errorf("unknown command")
	} else {
		cmd.rspName = value
	}
	e.cmdState = cmd
	return nil
}

func (e *emuImpl) GetResponse() (Message, error) {
	select {
	case resp := <-e.responses:
		return resp, nil
	case <-time.After(e.opt.TimeOut):
		return nil, fmt.Errorf("response timeout")
	}
}

func (e *emuImpl) GetMessage() (Message, error) {
	e.sendMsg = true
	select {
	case resp := <-e.responses:
		return resp, nil
	case <-time.After(e.opt.TimeOut):
		return nil, fmt.Errorf("response timeout")
	}
}
func (e *emuImpl) Close() {
	e.cancel()
	e.conn.Close()
}

func (e *emuImpl) reader() {
	rp := newResponseProcessor()
	reader := bufio.NewReader(e.conn)
	for {
		select {
		case <-e.ctx.Done():
			log.Printf("context done")
			return
		default:
			if reader.Buffered() == 0 {
				if rp.state != RspReceiving && e.cmdState != nil {
					e.sendCommand()
				} else if rp.state == RspReceiving {
					//					log.Printf("processing response %s", rp.resp.GetName())
				} else {
					//					log.Printf("nothing to do")
					//	time.Sleep(100 * time.Millisecond)
				}
			}
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					WarningLogger.Printf("EOF: nothing to read")
					break
				}
				ErrorLogger.Printf("Read error: %v", err)
				break
			}
			rp.process(line)
			line = ""
			if rp.state == RspReceived {
				//	log.Printf("message received:%+v", rp.resp)
				responseCache[rp.resp.Name] = rp.resp
				if e.cmdState != nil && e.cmdState.Status == CmdSent {
					//check if response is for the command
					msg, ok := responseCache[e.cmdState.rspName]
					if ok {
						//	log.Printf("response received: %s", msg.GetName())
						e.responses <- msg
						e.cmdState = nil
					}
				}
				if e.sendMsg {
					e.responses <- rp.resp
					e.sendMsg = false
				}
				rp = newResponseProcessor()
			} else if rp.state == RspError {
				ErrorLogger.Println("error in processing response")
				rp = newResponseProcessor()
			}
		}
	}
}

func (e *emuImpl) sendCommand() error {

	if e.cmdState.Status == CmdPending {
		xmlCmd := "<Command><Name>" + e.cmdState.command.GetName() + "</Name></Command>"
		InfoLogger.Printf("sending command: %s", string(xmlCmd))
		_, err := e.conn.Write([]byte(xmlCmd))
		if err != nil {
			e.cmdState.Status = CmdError
			return fmt.Errorf("write failed: %w", err)
		}
		e.cmdState.Status = CmdSent
	}
	return nil
}

type rspState int

const (
	RspPending rspState = iota
	RspReceiving
	RspReceived
	RspError
	RspTimeout
	RspUnknown
)

type responseProcessor struct {
	state rspState
	resp  *messageImpl
}

type cmdStatus int

const (
	CmdPending cmdStatus = iota
	CmdSent
	CmdReceived
	CmdError
	CmdTimeout
	CmdUnknown
)

type commandState struct {
	Status  cmdStatus
	rspName string
	command Command
}

func newCommandState() *commandState {
	return &commandState{Status: CmdUnknown}
}

func newResponseProcessor() *responseProcessor {
	return &responseProcessor{state: RspPending, resp: &messageImpl{}}
}

func startResponseTag(line string) (string, bool) {
	for k := range responseCache {
		if strings.Contains(line, "<"+k+">") {
			return k, true
		}
	}
	return "", false
}

func stopResponseTag(line string) (string, bool) {
	for k := range responseCache {
		if strings.Contains(line, "</"+k+">") {
			return k, true
		}
	}
	return "", false
}
func (rp *responseProcessor) process(line string) {
	tag, ok := startResponseTag(line)
	if ok {
		if rp.state == RspPending {
			rp.state = RspReceiving
			rp.resp.Name = tag
			rp.resp.Attribs = make(map[string]interface{})
			//			log.Printf("start:%s", rp.resp.GetName())
		} else {
			WarningLogger.Printf("response already in progress, abandoning the previous response %s", rp.resp.Name)
			rp.state = RspReceiving
			rp.resp.Name = tag
			rp.resp.Attribs = make(map[string]interface{})
			//			log.Printf("start:%s", rp.resp.Name)
		}
		return
	}
	tag, ok = stopResponseTag(line)
	if ok {
		if rp.state == RspReceiving && rp.resp.GetName() == tag {
			rp.state = RspReceived
		} else {
			WarningLogger.Printf("invalid end of response %s received. expecting[%s, %d]", tag, rp.resp.GetName(), rp.state)
			rp.state = RspError
		}
		return

	}
	if rp.state == RspReceiving {
		//parse xml element from the line with <key>vale</key>
		//add element and value to the response map
		key, value, err := xml2kv(line)
		if err != nil {
			ErrorLogger.Printf("xml parse error:[%s] %v", line, err)
			return
		}
		rp.resp.Attribs[key] = value
	} else {
		ErrorLogger.Printf("invalid response status %d to receive line: %s", rp.state, line)
	}
}

// type Element struct {
// 	XMLName xml.Name
// 	Value   string `xml:",chardata"`
// }
