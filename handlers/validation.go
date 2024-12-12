package handlers

import (
	"bytes"
	"log"
	"trucker/validators"
)

type Procotol_Message struct {
	Protocol_name string
	Raw_message   []byte
}

func IdentifyProtocol(input chan []byte, output chan Procotol_Message) {
	device_message := <-input
	var protocol_message Procotol_Message
	if bytes.HasPrefix(device_message, []byte{0x7e}) && bytes.HasSuffix(device_message, []byte{0x7e}) {
		protocol_message.Raw_message = device_message
		protocol_message.Protocol_name = "Micodus"
		output <- protocol_message
	}
}

func ValidateMessage(input chan Procotol_Message, output chan Procotol_Message) {
	log.Println("Validating message")
	log.SetPrefix("ValidateMessage(): ")
	protocol_message := <-input
	validationResult := true
	switch protocol_message.Protocol_name {
	case "Micodus":
		var micodusValidator validators.MicodusValidator
		validationResult = micodusValidator.Validate(protocol_message.Raw_message)
	}

	if validationResult {
		output <- protocol_message
	}

}