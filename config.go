package emu

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
		"Warning":                   nil,
		"Error":                     nil,
	}

	cmdRspMap = map[string]string{
		"restart":               "Warning",
		"get_device_info":       "DeviceInfo",
		"get_network_info":      "NetworkInfo",
		"get_time":              "TimeCluster",
		"get_connection_status": "ConnectionStatus",
		"get_message":           "MessageCluster",
		"get_fast_poll_status":  "FastPollStatus",
		// Simple Metering Commands
		"get_current_summation_delivered": "CurrentSummationDelivered",
		"get_instantaneous_demand":        "InstantaneousDemand",
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
)
