package emu

import "time"

const (
	closingGracePeriord time.Duration = time.Second * 5
)

type atrribType uint8

const (
	EPOCH atrribType = iota + 1
	INT64
	UINT64
	UINT32
	UINT16
	UINT8
	BOOLEAN
	STRING
)

func (a atrribType) String() string {
	switch a {
	case EPOCH:
		return "EPOCH"
	case INT64:
		return "INT64"
	case UINT64:
		return "UINT64"
	case UINT32:
		return "UINT32"
	case UINT16:
		return "UINT16"
	case UINT8:
		return "UINT8"
	case BOOLEAN:
		return "BOOLEAN"
	case STRING:
		return "STRING"
	default:
		return "UNKNOWN"
	}
}

var (
	responseCache = map[string](Message){
		"NetworkInfo":               nil,
		"ApsTable":                  nil,
		"Information":               nil,
		"TimeCluster":               nil,
		"NwkTable":                  nil,
		"PriceCluster":              nil,
		"DeviceInfo":                nil,
		"Google":                    nil,
		"SimpleMeteringCluster":     nil,
		"InstantaneousDemand":       nil,
		"BlockPriceDetail":          nil,
		"ConnectionStatus":          nil,
		"BillingPeriodList":         nil,
		"MessageCluster":            nil,
		"FastPollStatus":            nil,
		"CurrentSummationDelivered": nil,
		"ScheduleInfo":              nil,
		"Warning":                   nil,
		"Error":                     nil,
		"Ack":                       &messageImpl{Id: -1, Name: "Ack", Attribs: map[string]interface{}{"Status": "Success"}},
	}

	cmdIdcmdMap = map[CommandId]string{
		RESTART:         "restart",
		GET_DEVICE_INFO: "get_device_info",
		GET_TIME:        "get_time",
		GET_CONN_STATUS: "get_connection_status",
	}

	cmdRspMap = map[string]string{
		"restart":               "Ack",
		"get_device_info":       "DeviceInfo",
		"get_network_info":      "NetworkInfo",
		"get_time":              "TimeCluster",
		"get_connection_status": "ConnectionStatus",
		"get_message":           "MessageCluster",
		"get_fast_poll_status":  "FastPollStatus",
		// Simple Metering Commands
		"get_current_summation_delivered": "CurrentSummationDelivered",
		"get_instantaneous_demand":        "InstantaneousDemand",
		//Experimental
		"get_local_attributes": "Ack",
		"get_price_blocks":     "Ack",
		"get_schedule":         "Ack",
		"get_profile_data":     "Ack",

		//"get_last_period_usage":           "Warning",
		//"get_price":                "BlockPriceDetail",
		//"get_billing_period":   "BillingPeriodList",
		//	"get_aps_table":        "ApsTable",
		//"get_information":   "Information",
		//"print_network_tables": "NwkTable",
		//"get_price_cluster":    "PriceCluster",
		//	"get_google":               "Google",
		//"get_simple_metering": "SimpleMeteringCluster",
		//	"get_restart_info":         "Warning",
	}

	attribTypeMap = map[string]atrribType{
		"DeviceMacId": STRING,
		"MeterMacId":  STRING,
		"TimeStamp":   EPOCH,
		"Message":     STRING,
		"FastPoll":    BOOLEAN,
		// Name:InstantaneousDemand Attribs:map[Demand:0x000d3d DeviceMacId:0xd8d5b90000011821 DigitsLeft:0x05
		// DigitsRight:0x03 Divisor:0x000003e8 MeterMacId:0x001c640010ea31ad Multiplier:0x00000003 SuppressLeadingZero:Y
		// TimeStamp:0x2f520b33]
		"DigitsRight":         UINT8,
		"DigitsLeft":          UINT8,
		"SuppressLeadingZero": BOOLEAN,
		"Divisor":             INT64,
		"Multiplier":          UINT32,
		"Demand":              UINT32,
		// Name:CurrentSummationDelivered Attribs:map[DeviceMacId:0xd8d5b90000011821 DigitsLeft:0x05 DigitsRight:0x04
		// Divisor:0x000003e8 MeterMacId:0x001c640010ea31ad Multiplier:0x00000003 SummationDelivered:0x0000000001cde3d8
		// SummationReceived:0x0000000000000000 SuppressLeadingZero:Y TimeStamp:0x2f520bb8]
		"SuppressTrailingZero": BOOLEAN,
		"SummationDelivered":   UINT64,
		"SummationReceived":    UINT64,
		// Name:ConnectionStatus Attribs:map[Channel:20 Description:Successfully Joined DeviceMacId:0xd8d5b90000011821
		// ExtPanId:0x001c640010ea31ad LinkStrength:0x46 MeterMacId:0x001c640010ea31ad ShortAddr:0xbd2c Status:Connected]
		"Channel":      STRING,
		"ExtPanId":     STRING,
		"LinkStrength": UINT8,
		"Description":  STRING,
		"ShortAddr":    STRING,
		"Status":       STRING,
		// TimeCluster Attribs:map[DeviceMacId:0xd8d5b90000011821 LocalTime:0x2f519517 MeterMacId:0x001c640010ea31ad
		// UTCTime:0x2f520597]
		"LocalTime": EPOCH,
		"UTCTime":   EPOCH,
		// Name:DeviceInfo Attribs:map[DateCode:20211020355a0605 DeviceMacId:0xd8d5b90000011821
		// FWVersion:2.0.0 (7400) HWVersion:2.7.3 ImageType:0x2201 InstallCode:0x200a2c7d6b50ff8a
		// LinkKey:0xa09c9c4ad3a61e87e5beb4c5d186a377 Manufacturer:Rainforest Automation, Inc.
		// ModelId:Z105-2-EMU2-LEDD_JM]
		"DateCode":     STRING,
		"FWVersion":    STRING,
		"HWVersion":    STRING,
		"ImageType":    UINT16,
		"InstallCode":  UINT64,
		"LinkKey":      STRING,
		"Manufacturer": STRING,
		"ModelId":      STRING,
		// Name:NetworkInfo Attribs:map[Channel:20 CoordMacId:0x001c640010ea31ad Description:Successfully Joined
		// DeviceMacId:0xd8d5b90000011821 ExtPanId:0x001c640010ea31ad LinkStrength:0x64 ShortAddr:0xbd2c Status:Connected]
		"CoordMacId": STRING,
		// Name:MessageCluster Attribs:map[ConfirmationRequired:N Confirmed:N DeviceMacId:0xd8d5b90000011821
		// Duration: Id: MeterMacId:0x001c640010ea31ad Priority: Queue:Active StartTime: Text: TimeStamp:]
		"ConfirmationRequired": BOOLEAN,
		"Confirmed":            BOOLEAN,
		"Duration":             STRING, // @TODO: Need to verified with real data
		"Id":                   STRING, // @TODO: Need to verified with real data
		"Priority":             STRING, // @TODO: Need to verified with real data
		"Queue":                STRING,
		"StartTime":            EPOCH,  // @TODO: Need to verified with real data
		"Text":                 STRING, // @TODO: Need to verified with real data
		// FastPollStatus Attribs:map[DeviceMacId:0xd8d5b90000011821 EndTime:0x00000000 Frequency:0x00
		// MeterMacId:0x001c640010ea31ad]
		"EndTime":   UINT32, // @TODO: Need to verified with real data
		"Frequency": UINT32,
		// Name:ScheduleInfo Attribs:map[DeviceMacId:0xd8d5b90000011821 Enabled:Y Event:summation
		// Frequency:0x000000f0 MeterMacId:0x001c640010ea31ad Mode:rest]
		"Enabled": BOOLEAN,
		"Event":   STRING,
		"Mode":    STRING,
	}
)
