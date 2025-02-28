package emu

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	ErrTimeOut       = NewError("TIME_OUT")
	ErrChannelClosed = NewError("CHAN_CLOSED")
	ErrDeviceWrite   = NewError("DEVICE_WRITE")
	ErrDeviceRead    = NewError("DEVICE_READ")
	ErrDeviceIO      = NewError("DEVICE_IO")
	ErrMsgProc       = NewError("MESSAGE_PROCESSING")
)

type LogLevel int

const (
	LOG_ALL LogLevel = iota + 1
	LOG_INFO
	LOG_WARNING
	LOG_ERROR
	LOG_OFF
)

func StringToLogLevel(l string) (LogLevel, error) {
	switch l {
	case "LOG_ALL":
		return LOG_ALL, nil
	case "LOG_INFO":
		return LOG_INFO, nil
	case "LOG_WARNING":
		return LOG_WARNING, nil
	case "	LOG_ERROR":
		return LOG_ERROR, nil
	case "LOG_OFF":
		return LOG_OFF, nil
	default:
		return 0, fmt.Errorf("invalid log level: %s", l)
	}
}

type EmuOptions struct {
	BaudRate  int
	TimeOut   time.Duration
	LogWriter io.Writer
	LogLevel  LogLevel
}

type EmuOption func(*EmuOptions)

func WithBaudRate(baudRate int) EmuOption {
	return func(o *EmuOptions) {
		o.BaudRate = baudRate
	}
}

func WithTimeOut(timeout time.Duration) EmuOption {
	return func(o *EmuOptions) {
		o.TimeOut = timeout
	}
}

func WithLogWriter(w io.Writer) EmuOption {
	return func(o *EmuOptions) {
		o.LogWriter = w
	}
}
func WithLoggingLevel(l LogLevel) EmuOption {
	return func(o *EmuOptions) {
		o.LogLevel = l
	}
}

type Emu interface {
	SendCommand(Command) error
	GetResponse() (Message, error)
	GetMessage(names []MessageName) (Message, error)
	Start()
	Close()
	GetCumulativeEnergyConsumption() (*CumulativeEnergyConsumption, error)
	GetInstantaneousPowerConsumption() (*InstantaneousPowerConsumption, error)
}

func NewEmu(dev string, opts ...EmuOption) (Emu, error) {
	options := &EmuOptions{
		BaudRate:  115200,
		TimeOut:   15 * time.Second,
		LogWriter: os.Stdout,
		LogLevel:  LOG_ERROR,
	}
	for _, opt := range opts {
		opt(options)
	}

	return newEmuImpl(dev, options)
}

type CumulativeEnergyConsumption struct {
	TimeStamp   int64
	Energy      float64 //Unit is kWh
	DeviceMacId string
	MeterMacId  string
}

type InstantaneousPowerConsumption struct {
	TimeStamp   int64
	Power       float64 //Unit is KW
	DeviceMacId string
	MeterMacId  string
}

type MessageName string

const (
	DeviceInfo                MessageName = "DeviceInfo"
	NetworkInfo               MessageName = "NetworkInfo"
	TimeCluster               MessageName = "TimeCluster"
	InstantaneousDemand       MessageName = "InstantaneousDemand"
	CurrentSummationDelivered MessageName = "CurrentSummationDelivered"
)

type Message interface {
	GetName() string
	SetAttrib(string, interface{})
	GetAttrib(string) (interface{}, bool)
}

type CommandId int8

const (
	RESTART         CommandId = iota + 1 // restarts the emu-2 device
	GET_DEVICE_INFO                      // gets the basic emu-2 device info HW/SW version, make/model etc.
	GET_TIME                             // gets the time (local and UTC) on the emu-2 as sync with the smart energy meter
	GET_CONN_STATUS                      // gets the current connection status with the smart energy meter
)

func (c CommandId) String() string {
	if str, ok := commandIdString[c]; ok {
		return str
	}
	return "unknown"
}

func StrToCommandId(s string) (CommandId, error) {
	if c, ok := stringCommandId[s]; ok {
		return c, nil
	}
	return -1, fmt.Errorf("invalid string %s", s)
}

type Command interface {
	CommandId() CommandId
	SetAttrib(string, interface{})
}

func NewCommand(id CommandId) (Command, error) {
	if name, ok := cmdIdcmdMap[id]; ok {
		return &messageImpl{
			Id:   id,
			Name: name,
		}, nil
	}
	return nil, fmt.Errorf("invalid command id %+v", id)
}
