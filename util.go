package emu

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

func GetInstantaneousPowerConsumption(in Message) (*InstantaneousPowerConsumption, error) {
	DebugLogger.Printf("%+v", in)
	if msg, ok := in.(*messageImpl); ok {
		ipc := &InstantaneousPowerConsumption{}
		if ipc.TimeStamp, ok = msg.Attribs["TimeStamp"].(int64); !ok {
			return nil, fmt.Errorf("TimeStamp not found in message")
		}
		var power, multiplier, divisor, digitsRight int64
		if power, ok = msg.Attribs["Demand"].(int64); !ok {
			return nil, fmt.Errorf("demand not found in message")
		}
		if multiplier, ok = msg.Attribs["Multiplier"].(int64); !ok {
			return nil, fmt.Errorf("multiplier not found in message")
		}
		if divisor, ok = msg.Attribs["Divisor"].(int64); !ok {
			return nil, fmt.Errorf("divisor not found in message")
		}
		if divisor == 0 {
			return nil, fmt.Errorf("divisor cannot be zero")
		}
		if digitsRight, ok = msg.Attribs["DigitsRight"].(int64); !ok {
			return nil, fmt.Errorf("DigitsRight not found in message")
		}
		instaPower := roundToDecimal(float64(power*multiplier)/float64(divisor), int(digitsRight))
		ipc.Power = instaPower
		if ipc.DeviceMacId, ok = msg.Attribs["DeviceMacId"].(string); !ok {
			return nil, fmt.Errorf("DeviceMacId not found in message")
		}

		if ipc.MeterMacId, ok = msg.Attribs["MeterMacId"].(string); !ok {
			return nil, fmt.Errorf("MeterMacId not found in message")
		}
		return ipc, nil
	}
	return nil, fmt.Errorf("failed to cast message to messageImpl")
}

func GetCumulativeEnergyConsumption(in Message) (*CumulativeEnergyConsumption, error) {
	DebugLogger.Printf("%+v", in)
	if msg, ok := in.(*messageImpl); ok {
		cec := &CumulativeEnergyConsumption{}

		if cec.TimeStamp, ok = msg.Attribs["TimeStamp"].(int64); !ok {
			return nil, fmt.Errorf("TimeStamp not found in message")
		}
		var delivered, received, multiplier, divisor, digitsRight int64

		if delivered, ok = msg.Attribs["SummationDelivered"].(int64); !ok {
			return nil, fmt.Errorf("SummationDelivered not found in message")
		}
		if received, ok = msg.Attribs["SummationReceived"].(int64); !ok {
			return nil, fmt.Errorf("SummationReceived not found in message")
		}

		if multiplier, ok = msg.Attribs["Multiplier"].(int64); !ok {
			return nil, fmt.Errorf("multiplier not found in message")
		}
		if divisor, ok = msg.Attribs["Divisor"].(int64); !ok {
			return nil, fmt.Errorf("divisor not found in message")
		}
		if divisor == 0 {
			return nil, fmt.Errorf("divisor cannot be zero")
		}
		if digitsRight, ok = msg.Attribs["DigitsRight"].(int64); !ok {
			return nil, fmt.Errorf("DigitsRight not found in message")
		}
		energy := roundToDecimal(float64((delivered-received)*multiplier)/float64(divisor), int(digitsRight))
		cec.Energy = energy

		if cec.DeviceMacId, ok = msg.Attribs["DeviceMacId"].(string); !ok {
			return nil, fmt.Errorf("DeviceMacId not found in message")
		}

		if cec.MeterMacId, ok = msg.Attribs["MeterMacId"].(string); !ok {
			return nil, fmt.Errorf("MeterMacId not found in message")
		}
		return cec, nil
	}
	return nil, fmt.Errorf("failed to cast message to messageImpl")
}

func roundToDecimal(num float64, decimals int) float64 {
	pow10 := math.Pow(10, float64(decimals))
	return math.Round(num*pow10) / pow10
}

func getCorrectTimeStamp(ts int64) int64 {
	return time.Unix(ts, 0).AddDate(30, 0, -1).Unix()
}

func structToMap(obj interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i).Interface()
		out[field.Name] = fieldValue
	}
	return out
}
