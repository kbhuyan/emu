package emu

import (
	"fmt"
	"io"
	"os"
	"time"
)

type EmuOptions struct {
	BaudRate  int
	TimeOut   time.Duration
	LogWriter io.Writer
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

type Emu interface {
	SendCommand(Command) error
	GetResponse() (Message, error)
	GetMessage() (Message, error)
	Start()
	Close()
}

func NewEmu(dev string, opts ...EmuOption) (Emu, error) {
	options := &EmuOptions{
		BaudRate:  115200,
		TimeOut:   15 * time.Second,
		LogWriter: os.Stdout,
	}
	for _, opt := range opts {
		opt(options)
	}

	return newEmuImpl(dev, options)
}

type Message interface {
	GetName() string
	SetAttrib(string, interface{})
	GetAttrib(string) (interface{}, bool)
}

type Command interface {
	GetName() string
	SetAttrib(string, interface{})
}

func NewCommand(name string) (Command, error) {
	_, ok := cmdRspMap[name]
	if !ok {
		return nil, fmt.Errorf("invalid command")
	}
	return &messageImpl{
		Name: name,
	}, nil
}
