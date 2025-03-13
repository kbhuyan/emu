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

type emuMessageName string

const (
	emuNetworkInfo               emuMessageName = "NetworkInfo"
	emuApsTable                  emuMessageName = "ApsTable"
	emuInformation               emuMessageName = "Information"
	emuTimeCluster               emuMessageName = "TimeCluster"
	emuNwkTable                  emuMessageName = "NwkTable"
	emuPriceCluster              emuMessageName = "PriceCluster"
	emuDeviceInfo                emuMessageName = "DeviceInfo"
	emuGoogle                    emuMessageName = "Google"
	emuSimpleMeteringCluster     emuMessageName = "SimpleMeteringCluster"
	emuInstantaneousDemand       emuMessageName = "InstantaneousDemand"
	emuBlockPriceDetail          emuMessageName = "BlockPriceDetail"
	emuConnectionStatus          emuMessageName = "ConnectionStatus"
	emuBillingPeriodList         emuMessageName = "BillingPeriodList"
	emuMessageCluster            emuMessageName = "MessageCluster"
	emuFastPollStatus            emuMessageName = "FastPollStatus"
	emuCurrentSummationDelivered emuMessageName = "CurrentSummationDelivered"
	emuScheduleInfo              emuMessageName = "ScheduleInfo"
	emuWarning                   emuMessageName = "Warning"
	emuAck                       emuMessageName = "Ack"
)

type emuCommandName string

const (
	emuRestart                      emuCommandName = "restart"
	emuGetDeviceInfo                emuCommandName = "get_device_info"
	emuGetNetworkInfo               emuCommandName = "get_network_info"
	emuGetTime                      emuCommandName = "get_time"
	emuGetConnStatus                emuCommandName = "get_connection_status"
	emuGetMessage                   emuCommandName = "get_message"
	emuGetFastPollStatus            emuCommandName = "get_fast_poll_status"
	emuGetCurrentSummationDelivered emuCommandName = "get_current_summation_delivered"
	emuGetInstantaneousDemand       emuCommandName = "get_instantaneous_demand"
	emuGetLocalAttributes           emuCommandName = "get_local_attributes"
	emuGetPriceBlocks               emuCommandName = "get_price_blocks"
	emuGetSchedule                  emuCommandName = "get_schedule"
	emuGetProfileData               emuCommandName = "get_profile_data"
)

type emuMessageAttribute string

const (
	emuDeviceMacId          emuMessageAttribute = "DeviceMacId"
	emuMeterMacId           emuMessageAttribute = "MeterMacId"
	emuTimeStamp            emuMessageAttribute = "TimeStamp"
	emuMessage              emuMessageAttribute = "Message"
	emuFastPoll             emuMessageAttribute = "FastPoll"
	emuDigitsRight          emuMessageAttribute = "DigitsRight"
	emuDigitsLeft           emuMessageAttribute = "DigitsLeft"
	emuSuppressLeadingZero  emuMessageAttribute = "SuppressLeadingZero"
	emuDivisor              emuMessageAttribute = "Divisor"
	emuMultiplier           emuMessageAttribute = "Multiplier"
	emuDemand               emuMessageAttribute = "Demand"
	emuSuppressTrailingZero emuMessageAttribute = "SuppressTrailingZero"
	emuSummationDelivered   emuMessageAttribute = "SummationDelivered"
	emuSummationReceived    emuMessageAttribute = "SummationReceived"
	emuChannel              emuMessageAttribute = "Channel"
	emuExtPanId             emuMessageAttribute = "ExtPanId"
	emuLinkStrength         emuMessageAttribute = "LinkStrength"
	emuDescription          emuMessageAttribute = "Description"
	emuShortAddr            emuMessageAttribute = "ShortAddr"
	emuStatus               emuMessageAttribute = "Status"
	emuLocalTime            emuMessageAttribute = "LocalTime"
	emuUTCTime              emuMessageAttribute = "UTCTime"
	emuDateCode             emuMessageAttribute = "DateCode"
	emuFWVersion            emuMessageAttribute = "FWVersion"
	emuHWVersion            emuMessageAttribute = "HWVersion"
	emuImageType            emuMessageAttribute = "ImageType"
	emuInstallCode          emuMessageAttribute = "InstallCode"
	emuLinkKey              emuMessageAttribute = "LinkKey"
	emuManufacturer         emuMessageAttribute = "Manufacturer"
	emuModelId              emuMessageAttribute = "ModelId"
	emuCoordMacId           emuMessageAttribute = "CoordMacId"
	emuConfirmationRequired emuMessageAttribute = "ConfirmationRequired"
	emuConfirmed            emuMessageAttribute = "Confirmed"
	emuDuration             emuMessageAttribute = "Duration"
	emuId                   emuMessageAttribute = "Id"
	emuPriority             emuMessageAttribute = "Priority"
	emuQueue                emuMessageAttribute = "Queue"
	emuStartTime            emuMessageAttribute = "StartTime"
	emuText                 emuMessageAttribute = "Text"
	emuEndTime              emuMessageAttribute = "EndTime"
	emuFrequency            emuMessageAttribute = "Frequency"
	emuEnabled              emuMessageAttribute = "Enabled"
	emuEvent                emuMessageAttribute = "Event"
	emuMode                 emuMessageAttribute = "Mode"
)

type emMessage2ApiMessage func(*messageImpl) (Message, error)

var (
	messageProcessorMap = map[emuMessageName]emMessage2ApiMessage{
		emuCurrentSummationDelivered: emuCurrentSummationDelivered2CumulativeEnergy,
		emuInstantaneousDemand:       emuInstantaneousDemand2InstantaneousPower,
	}

	apiMessageNames = []MessageName{
		DeviceInfo, NetworkInfo, TimeCluster, InstantaneousPower, CumulativeEnergy,
	}
	emuResponses = []emuMessageName{
		emuNetworkInfo,
		emuApsTable,
		emuInformation,
		emuTimeCluster,
		emuNwkTable,
		emuPriceCluster,
		emuDeviceInfo,
		emuGoogle,
		emuSimpleMeteringCluster,
		emuInstantaneousDemand,
		emuBlockPriceDetail,
		emuConnectionStatus,
		emuBillingPeriodList,
		emuMessageCluster,
		emuFastPollStatus,
		emuCurrentSummationDelivered,
		emuScheduleInfo,
		emuWarning,
		emuAck,
	}

	cmdIdcmdMap = map[CommandId]emuCommandName{
		RESTART:         emuRestart,
		GET_DEVICE_INFO: emuGetDeviceInfo,
		GET_TIME:        emuGetTime,
		GET_CONN_STATUS: emuGetConnStatus,
	}

	cmdRspMap = map[emuCommandName]emuMessageName{
		emuRestart:                      emuAck,
		emuGetDeviceInfo:                emuDeviceInfo,
		emuGetNetworkInfo:               emuNetworkInfo,
		emuGetTime:                      emuTimeCluster,
		emuGetConnStatus:                emuConnectionStatus,
		emuGetMessage:                   emuMessageCluster,
		emuGetFastPollStatus:            emuFastPollStatus,
		emuGetCurrentSummationDelivered: emuCurrentSummationDelivered,
		emuGetInstantaneousDemand:       emuInstantaneousDemand,
		emuGetLocalAttributes:           emuAck,
		emuGetPriceBlocks:               emuAck,
		emuGetSchedule:                  emuAck,
		emuGetProfileData:               emuAck,
	}

	attribTypeMap = map[emuMessageAttribute]atrribType{
		emuDeviceMacId:          STRING,
		emuMeterMacId:           STRING,
		emuTimeStamp:            EPOCH,
		emuMessage:              STRING,
		emuFastPoll:             BOOLEAN,
		emuDigitsRight:          UINT8,
		emuDigitsLeft:           UINT8,
		emuSuppressLeadingZero:  BOOLEAN,
		emuDivisor:              INT64,
		emuMultiplier:           UINT32,
		emuDemand:               UINT32,
		emuSuppressTrailingZero: BOOLEAN,
		emuSummationDelivered:   UINT64,
		emuSummationReceived:    UINT64,
		emuChannel:              STRING,
		emuExtPanId:             STRING,
		emuLinkStrength:         UINT8,
		emuDescription:          STRING,
		emuShortAddr:            STRING,
		emuStatus:               STRING,
		emuLocalTime:            EPOCH,
		emuUTCTime:              EPOCH,
		emuDateCode:             STRING,
		emuFWVersion:            STRING,
		emuHWVersion:            STRING,
		emuImageType:            UINT16,
		emuInstallCode:          UINT64,
		emuLinkKey:              STRING,
		emuManufacturer:         STRING,
		emuModelId:              STRING,
		emuCoordMacId:           STRING,
		emuConfirmationRequired: BOOLEAN,
		emuConfirmed:            BOOLEAN,
		emuDuration:             STRING,
		emuId:                   STRING,
		emuPriority:             STRING,
		emuQueue:                STRING,
		emuStartTime:            EPOCH,
		emuText:                 STRING,
		emuEndTime:              UINT32,
		emuFrequency:            UINT32,
		emuEnabled:              BOOLEAN,
		emuEvent:                STRING,
		emuMode:                 STRING,
	}
)
