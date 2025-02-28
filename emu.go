package emu

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

type messageImpl struct {
	Id      CommandId
	Name    string
	Attribs map[string]interface{}
}

func (m *messageImpl) GetName() string {
	return m.Name
}

func (m *messageImpl) CommandId() CommandId {
	return m.Id
}

func (m *messageImpl) SetAttrib(key string, value interface{}) {
	m.Attribs[key] = value
}

func (m *messageImpl) GetAttrib(key string) (interface{}, bool) {
	value, ok := m.Attribs[key]
	return value, ok
}

type emuImpl struct {
	conn         io.ReadWriteCloser
	responses    chan Message
	ctx          context.Context
	cancel       context.CancelFunc
	cmdState     *commandState
	opt          *EmuOptions
	sendMsg      bool
	sendMsgNames []MessageName
}

func newEmuImpl(dev string, opt *EmuOptions) (Emu, error) {
	initLog(opt.LogWriter, opt.LogLevel)
	// Configure serial port
	mode := &serial.Mode{
		BaudRate: opt.BaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(dev, mode)
	if err != nil {
		return nil, ErrDeviceIO.Errorf("serial open failed: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &emuImpl{
		conn:      port,
		responses: make(chan Message, 1),
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

func (e *emuImpl) SendCommand(c Command) error {
	if command, ok := c.(*messageImpl); ok {
		if value, exists := cmdRspMap[command.GetName()]; exists {
			time.Sleep(100 * time.Millisecond)
			e.cmdState = &commandState{command: command, status: CmdPending, rspName: value}
			return nil
		}
	}
	return fmt.Errorf("invalid command type %T or %+v", c, c)
}

func (e *emuImpl) GetResponse() (Message, error) {
	select {
	case resp := <-e.responses:
		return resp, nil
	case <-time.After(e.opt.TimeOut):
		return nil, ErrTimeOut
	case <-e.ctx.Done():
		return nil, ErrChannelClosed.Errorf("channel closed %+v", e.ctx.Err())
	}
}

func (e *emuImpl) GetMessage(names []MessageName) (Message, error) {
	e.sendMsg = true
	e.sendMsgNames = names
	select {
	case resp := <-e.responses:
		return resp, nil
	case <-time.After(e.opt.TimeOut):
		return nil, ErrTimeOut
	case <-e.ctx.Done():
		return nil, ErrChannelClosed.Errorf("channel closed %+v", e.ctx.Err())
	}
}
func (e *emuImpl) Close() {
	InfoLogger.Println("closing the emu session.")
	e.cancel()
	time.Sleep(closingGracePeriord)
	e.conn.Close()
}

func (e *emuImpl) GetCumulativeEnergyConsumption() (*CumulativeEnergyConsumption, error) {

	cmd := &messageImpl{Name: "get_current_summation_delivered"}
	if err := e.SendCommand(cmd); err != nil {
		return nil, err
	}
	if rsp, err := e.GetResponse(); err != nil {
		return nil, err
	} else {
		return GetCumulativeEnergyConsumption(rsp)
	}

}

func (e *emuImpl) GetInstantaneousPowerConsumption() (*InstantaneousPowerConsumption, error) {
	cmd := &messageImpl{Name: "get_instantaneous_demand"}
	if err := e.SendCommand(cmd); err != nil {
		return nil, err
	}
	if rsp, err := e.GetResponse(); err != nil {
		return nil, err
	} else {
		return GetInstantaneousPowerConsumption(rsp)
	}
}

func (e *emuImpl) reader() {
	rp := newResponseProcessor()
	reader := bufio.NewReader(e.conn)
	for {
		select {
		case <-e.ctx.Done():
			InfoLogger.Printf("context done. %+v", e.ctx.Err())
			return
		default:
			if reader.Buffered() == 0 {
				if rp.state != RspReceiving && e.cmdState != nil {
					e.sendCommand()
				}
			}
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					WarningLogger.Printf("EOF: nothing to read")
					break
				} else {
					ErrorLogger.Printf("Read error: %v", err)
				}
				break
			}
			rp.process(line)
			if rp.state == RspReceived {
				responseCache[rp.resp.Name] = rp.resp
				if e.cmdState != nil && e.cmdState.status == CmdSent {
					//check if response is for the command
					msg, ok := responseCache[e.cmdState.rspName]
					if ok && msg != nil {
						e.responses <- msg
						e.cmdState = nil
					}
				}
				if e.sendMsg && e.isSendMsgName(rp.resp.Name) {
					e.responses <- rp.resp
					e.sendMsg = false
				}
				rp = newResponseProcessor()
			} else if rp.state == RspError {
				WarningLogger.Printf("Abandoning processing response: [%s, %+v]\n", rp.state, rp.resp)
				rp = newResponseProcessor()
			}
		}
	}
}

func (e *emuImpl) isSendMsgName(n string) bool {
	nm := MessageName(n)
	if len(e.sendMsgNames) > 0 {
		for _, name := range e.sendMsgNames {
			if name == nm {
				return true
			}
		}
		return false
	} else {
		return true
	}
}

func (e *emuImpl) sendCommand() error {

	if e.cmdState.status == CmdPending {
		xmlCmd := "<Command><Name>" + e.cmdState.command.GetName() + "</Name></Command>"
		DebugLogger.Printf("sending command: %s", string(xmlCmd))
		if _, err := e.conn.Write([]byte(xmlCmd)); err != nil {
			e.cmdState.status = CmdError
			return ErrDeviceWrite.Errorf("error while writing to devive %+v", err)
		}
		e.cmdState.status = CmdSent
	}
	return nil
}

type rspState int

const (
	RspPending rspState = iota + 1
	RspReceiving
	RspReceived
	RspError
	RspTimeout
	RspUnknown
)

func (s rspState) String() string {
	switch s {
	case RspPending:
		return "RspPending"
	case RspReceiving:
		return "RspReceiving"
	case RspReceived:
		return "RspReceived"
	case RspError:
		return "RspError"
	case RspTimeout:
		return "RspTimeout"
	case RspUnknown:
		return "RspUnknown"
	default:
		return "Invalid"
	}
}

type responseProcessor struct {
	state rspState
	resp  *messageImpl
}

type cmdStatus int

const (
	CmdPending cmdStatus = iota + 1
	CmdSent
	CmdReceived
	CmdError
	CmdTimeout
	CmdUnknown
)

func (c cmdStatus) String() string {
	switch c {
	case CmdPending:
		return "CmdPending"
	case CmdSent:
		return "CmdSent"
	case CmdReceived:
		return "CmdReceived"
	case CmdError:
		return "CmdError"
	case CmdTimeout:
		return "CmdTimeout"
	case CmdUnknown:
		return "CmdUnknown"
	default:
		return "Invalid"
	}
}

type commandState struct {
	status  cmdStatus
	rspName string
	command *messageImpl
}

func newCommandState() *commandState {
	return &commandState{status: CmdUnknown}
}

func newResponseProcessor() *responseProcessor {
	return &responseProcessor{state: RspPending, resp: &messageImpl{Attribs: make(map[string]interface{})}}
}

func (rp *responseProcessor) startResponseTag(line string) (string, bool) {
	for k := range responseCache {
		if strings.HasPrefix(line, "<"+k+">") {
			return k, true
		}
	}
	return "", false
}

func (rp *responseProcessor) stopResponseTag(line string) (string, bool) {
	for k := range responseCache {
		if strings.HasPrefix(line, "</"+k+">") {
			return k, true
		}
	}
	return "", false
}

func (rp *responseProcessor) getAttrib(line string) (key string, value interface{}, err error) {
	var element struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	}
	if err = xml.Unmarshal([]byte(line), &element); err != nil {
		return "", nil, ErrMsgProc.Errorf("unable to parse xml: %s, %+v", line, err)
	}
	key = element.XMLName.Local
	strValue := element.Value
	if at, ok := attribTypeMap[key]; ok {
		var err error = nil
		switch at {
		case INT64:
			value, err = strconv.ParseInt(strValue, 0, 64)
		case UINT64:
			value, err = strconv.ParseInt(strValue, 0, 64)
		case UINT32:
			value, err = strconv.ParseInt(strValue, 0, 32)
		case UINT16:
			value, err = strconv.ParseInt(strValue, 0, 16)
		case UINT8:
			value, err = strconv.ParseInt(strValue, 0, 16)
		case BOOLEAN:
			if strValue == "Y" {
				value = true
			} else {
				value = false
			}
		case EPOCH:
			if tv, err := strconv.ParseInt(strValue, 0, 64); err == nil {
				value = getCorrectTimeStamp(tv)
			}
		case STRING:
			value = strValue
		default:
			err = fmt.Errorf("invalid attrib type %s", at)
		}
		if err != nil {
			return "", nil, ErrMsgProc.Errorf("unable to convert %s's value %s to type %s. %+v", key, strValue, at, err)
		}
	}
	return key, value, nil
}

func (rp *responseProcessor) process(line string) {
	line = strings.TrimLeft(line, " \t")
	DebugLogger.Println("Processing:", line)
	tag, ok := rp.startResponseTag(line)
	if ok {
		if rp.state == RspPending {
			rp.state = RspReceiving
			rp.resp.Name = tag
		} else {
			WarningLogger.Printf("response for %s was in progress, abandoning. now starting for %s", rp.resp.Name, tag)
			rp.state = RspReceiving
			rp.resp.Name = tag
		}
		return
	}
	tag, ok = rp.stopResponseTag(line)
	if ok {
		if rp.state == RspReceiving && rp.resp.GetName() == tag {
			rp.state = RspReceived
		} else {
			WarningLogger.Printf("invalid end of response %s received. expecting[%s, %s]. line: %s", tag, rp.resp.GetName(), rp.state, line)
			rp.state = RspError
		}
		return

	}
	if rp.state == RspReceiving {
		//parse xml element from the line with <key>vale</key>
		//add element and value to the response map
		key, value, err := rp.getAttrib(line)
		if err != nil {
			WarningLogger.Printf("abandoning message %s as xml parse error:%v while processing. line: %s", rp.resp.GetName(), err, line)
			rp.state = RspError
		} else {
			rp.resp.Attribs[key] = value
		}
	} else {
		WarningLogger.Printf("ignoring as invalid response state %s to receive line: %s", rp.state, line)
		if rp.state != RspPending {
			rp.state = RspError
		}
	}
}

func (rp *responseProcessor) processv2(line string) {
	//if state is RspReceiving then look for stopResponseTag and attributes
	//else ignore line as it start to receive in the middle of an response
	switch rp.state {
	case RspPending:
		if tag, ok := rp.startResponseTag(line); ok {
			rp.state = RspReceiving
			rp.resp.Name = tag
		} else {
			WarningLogger.Printf("starting to receive in the middle of the message, ignoring. line: %s", line)
		}
	case RspReceiving:
		if tag, ok := rp.stopResponseTag(line); ok {
			if rp.resp.GetName() == tag {
				rp.state = RspReceived
			} else {
				WarningLogger.Printf("invalid end of response %s received. expecting[%s, %+v]. line: %s", tag, rp.resp.GetName(), rp.state, line)
				rp.state = RspError
			}
		} else {
			//parse xml element from the line with <key>vale</key>
			//add key and value to the response Attribs[key] = value
			key, value, err := rp.getAttrib(line)
			if err != nil {
				WarningLogger.Printf("abandoning message %s as xml parse error:%v while processing. line: %s", rp.resp.GetName(), err, line)
				rp.state = RspError
			} else {
				rp.resp.Attribs[key] = value
			}
		}
	default:
		ErrorLogger.Printf("invalid response state %+v to receive. line: %s", rp.state, line)
	}
}
