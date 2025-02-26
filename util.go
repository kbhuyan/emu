package emu

import (
	"encoding/xml"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

type InstantaneousPowerConsumption struct {
	TimeStamp   int64
	Power       float64 //Unit is KW
	DeviceMacId string
	MeterMacId  string
}

func GetInstantaneousPowerConsumption(in Message) (*InstantaneousPowerConsumption, error) {
	//	log.Printf("message type %T\n", in)
	if msg, ok := in.(*messageImpl); ok {
		ipc := &InstantaneousPowerConsumption{}
		value, ok := msg.Attribs["TimeStamp"].(string)
		if !ok {
			return nil, fmt.Errorf("TimeStamp not found in message")
		}
		ts, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TimeStamp value: %w", err)
		}
		ipc.TimeStamp = getCorrectTimeStamp(ts)

		value, ok = msg.Attribs["Demand"].(string)
		if !ok {
			return nil, fmt.Errorf("demand not found in message")
		}
		power, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse power value: %w", err)
		}
		value, ok = msg.Attribs["Multiplier"].(string)
		if !ok {
			return nil, fmt.Errorf("multiplier not found in message")
		}
		multiplier, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse multiplier value: %w", err)
		}
		value, ok = msg.Attribs["Divisor"].(string)
		if !ok {
			return nil, fmt.Errorf("divisor not found in message")
		}
		divisor, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse divisor value: %w", err)
		}
		value, ok = msg.Attribs["DigitsRight"].(string)
		if !ok {
			return nil, fmt.Errorf("DigitsRight not found in message")
		}
		digitsRight, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse digitsRight value: %w", err)
		}
		if divisor == 0 {
			return nil, fmt.Errorf("divisor cannot be zero")
		}
		instaPower := roundToDecimal(float64(power*multiplier)/float64(divisor), int(digitsRight))
		ipc.Power = instaPower
		value, ok = msg.Attribs["DeviceMacId"].(string)
		if !ok {
			return nil, fmt.Errorf("DeviceMacId not found in message")
		}
		ipc.DeviceMacId = value
		value, ok = msg.Attribs["MeterMacId"].(string)
		if !ok {
			return nil, fmt.Errorf("MeterMacId not found in message")
		}
		ipc.MeterMacId = value
		return ipc, nil
	}
	return nil, fmt.Errorf("failed to cast message to messageImpl")
}

type CumulativeEnergyConsumption struct {
	TimeStamp   int64
	Energy      float64 //Unit is kWh
	DeviceMacId string
	MeterMacId  string
}

func GetCumulativeEnergyConsumption(in Message) (*CumulativeEnergyConsumption, error) {
	//	log.Printf("message type %T\n", in)
	if msg, ok := in.(*messageImpl); ok {
		cec := &CumulativeEnergyConsumption{}

		value, ok := msg.Attribs["TimeStamp"].(string)
		if !ok {
			return nil, fmt.Errorf("TimeStamp not found in message")
		}
		ts, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TimeStamp value: %w", err)
		}
		cec.TimeStamp = getCorrectTimeStamp(ts)

		value, ok = msg.Attribs["SummationDelivered"].(string)
		if !ok {
			return nil, fmt.Errorf("SummationDelivered not found in message")
		}
		delivered, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SummationDelivered value: %w", err)
		}

		value, ok = msg.Attribs["SummationReceived"].(string)
		if !ok {
			return nil, fmt.Errorf("SummationReceived not found in message")
		}
		received, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SummationReceived value: %w", err)
		}

		value, ok = msg.Attribs["Multiplier"].(string)
		if !ok {
			return nil, fmt.Errorf("multiplier not found in message")
		}
		multiplier, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse multiplier value: %w", err)
		}
		value, ok = msg.Attribs["Divisor"].(string)
		if !ok {
			return nil, fmt.Errorf("divisor not found in message")
		}
		divisor, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse divisor value: %w", err)
		}
		value, ok = msg.Attribs["DigitsRight"].(string)
		if !ok {
			return nil, fmt.Errorf("DigitsRight not found in message")
		}
		digitsRight, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse digitsRight value: %w", err)
		}
		if divisor == 0 {
			return nil, fmt.Errorf("divisor cannot be zero")
		}
		energy := roundToDecimal(float64((delivered-received)*multiplier)/float64(divisor), int(digitsRight))
		cec.Energy = energy

		value, ok = msg.Attribs["DeviceMacId"].(string)
		if !ok {
			return nil, fmt.Errorf("DeviceMacId not found in message")
		}
		cec.DeviceMacId = value
		value, ok = msg.Attribs["MeterMacId"].(string)
		if !ok {
			return nil, fmt.Errorf("MeterMacId not found in message")
		}
		cec.MeterMacId = value
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

func xml2kv(line string) (key string, value string, err error) {
	var element struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	}
	err = xml.Unmarshal([]byte(line), &element)
	if err != nil {
		return "", "", err
	}
	key = element.XMLName.Local
	value = element.Value
	return key, value, nil
}
