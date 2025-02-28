package emu

var commandIdString = map[CommandId]string{
	RESTART:         "RESTART",
	GET_DEVICE_INFO: "GET_DEVICE_INFO",
	GET_TIME:        "GET_TIME",
	GET_CONN_STATUS: "GET_CONN_STATUS",
}

var stringCommandId = map[string]CommandId{
	"RESTART":         RESTART,
	"GET_DEVICE_INFO": GET_DEVICE_INFO,
	"GET_TIME":        GET_TIME,
	"GET_CONN_STATUS": GET_CONN_STATUS,
}
