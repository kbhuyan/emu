package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/kbhuyan/emu"
)

func main() {
	// Configure command-line flags
	port := flag.String("port", "/dev/ttyACM1", "Serial port device path")
	baud := flag.Int("baud", 115200, "Baud rate (115200, 9600, etc)")
	timeout := flag.Duration("timeout", 15*time.Second, "Read timeout duration")
	list := flag.Bool("list", false, "List available commands and exit")

	flag.Parse()

	// Handle --list flag
	if *list {
		printAvailableCommands()
		return
	}

	// Get command from positional arguments
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Usage: emuctl [flags] <command>\nAvailable commands: get_device_info, get_network_info, get_current_summation, get_instantaneous_demand")
	}
	command := args[0]

	device, err := emu.NewEmu(*port, emu.WithBaudRate(*baud), emu.WithTimeOut(*timeout))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer device.Close()
	device.Start()

	//fmt.Println("commanf: ", command)
	// Execute command
	cmd, err := emu.NewCommand(command)
	if err != nil {
		log.Fatalf("%v", err)
	}
	msg, err := executeCommand(device, cmd)
	if err != nil {
		log.Fatalf("%v", err)
	}
	processMessage(msg)

	for {
		msg, err := device.GetMessage()
		if err != nil {
			log.Fatalf("GetMessage failed: %v", err)
		}
		processMessage(msg)
	}
}

func processMessage(msg emu.Message) {
	if msg.GetName() == "InstantaneousDemand" {
		power, err := emu.GetInstantaneousPowerConsumption(msg)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("TimeStamp: %s Instantaneous Demand: %+v\n", time.Unix(power.TimeStamp, 0), power)
	} else if msg.GetName() == "CurrentSummationDelivered" {
		energy, err := emu.GetCumulativeEnergyConsumption(msg)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("TimeStamp: %s Current Summation Energy Delivered: %+v\n", time.Unix(energy.TimeStamp, 0), energy)
	} else {
		log.Printf("Message: %+v\n", msg)
	}

}

func executeCommand(device emu.Emu, cmd emu.Command) (emu.Message, error) {
	// Send command
	if err := device.SendCommand(cmd); err != nil {
		return nil, fmt.Errorf("command failed: %v", err)
	}

	// Get response
	//rawResp, err := device.ReadResponseSync()
	rsp, err := device.GetResponse()
	if err != nil {
		return nil, fmt.Errorf("response error: %v", err)
	}
	return rsp, nil
}

func printAvailableCommands() {
	fmt.Println(`Available EMU commands:
  get_device_info          - Device manufacturer and version info
  get_network_info         - MAC addresses and connection status
  get_current_summation    - Total energy consumption
  get_instantaneous_demand - Current power demand`)
}
