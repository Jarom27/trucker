package handlers

import "fmt"

type Protocol interface {
	ProcessMessage(message []byte) interface{}
}
type Response_GPS struct {
	Device_id []byte
	Message   []byte
}
type Location_report struct {
	Device_id string
	Latitude  float64
	Longitude float64
	Altitude  float64
}

func Process(input chan Procotol_Message, gps chan Response_GPS, location chan Location_report, leavings chan interface{}) {
	protocol_message := <-input
	var (
		protocol Protocol
	)
	switch protocol_message.Protocol_name {
	case "Micodus":
		protocol = Micodus{}
	}
	fmt.Printf("Processing this message %x\n", protocol_message.Raw_message)
	response := protocol.ProcessMessage(protocol_message.Raw_message)
	if value, ok := response.(Response_GPS); ok {
		fmt.Printf("Responding to GPS with this message %x\n", value.Message)
		gps <- value
	}
	if value, ok := response.(Location_report); ok {
		location <- value
	}
	if response == nil {
		leavings <- response
	}

}
