package handlers

import (
	"encoding/binary"
	"fmt"
	"trucker/utilities"
)

type Micodus struct {
}
type micodusHeader struct {
	device_id     []byte
	message_type  []byte
	length        []byte
	serial_number []byte
}

// ProcessMessage implements Protocol.
func (t Micodus) ProcessMessage(message []byte) interface{} {
	var message_type string = t.getMessageType(message[0:3])
	fmt.Println("Message type: ", message_type)
	switch message_type {
	case "0100":
		return t.registrationComand(message)
	case "0102":
		fmt.Println("Authenticate Command is processing")
		return t.authenticationCommand(message)
	case "0003":
		return nil
	case "0201":
		fmt.Println("Parsing location data")
		return t.parseData(message)
	}
	return nil
}

func (t *Micodus) getMessageType(message []byte) string {
	message_type := fmt.Sprintf("%02x%02x", message[1], message[2])
	return message_type
}
func (t *Micodus) getDeviceId(message []byte) string {
	var device_id string = ""
	for _, value := range message {
		device_id += fmt.Sprintf("%2x", value)
	}
	return device_id
}
func (t *Micodus) buildResponse(header micodusHeader, body []byte, result []byte) []byte {
	length := int(binary.BigEndian.Uint16(header.length[0:2]))
	response := make([]byte, 0, length)
	response = append(response, 0x7e)
	response = append(response, header.message_type...)
	response = append(response, header.length...)
	response = append(response, header.device_id...)
	response = append(response, header.serial_number...)
	if body != nil {
		response = append(response, body...)
	}
	if result != nil {
		response = append(response, result...)
	}

	checksum := utilities.XorChecksum(response[1:])
	response = append(response, checksum)
	response = append(response, 0x7e)

	return response
}
func (t *Micodus) registrationComand(message []byte) Response_GPS {
	var header micodusHeader = micodusHeader{}
	header.message_type = []byte{0x81, 0x00}
	header.device_id = message[5:11]
	header.serial_number = []byte{0x1, 0x00} // It's necessary to update the serial Number from the main entry
	header.length = []byte{0x00, 0x2}

	message_response := t.buildResponse(header, nil, []byte{0x00, 0x00})
	response_gps := Response_GPS{header.device_id, message_response}
	return response_gps
}
func (t *Micodus) authenticationCommand(message []byte) Response_GPS {
	var header micodusHeader = micodusHeader{}
	header.message_type = []byte{0x82, 0x01}
	header.device_id = message[5:11]
	header.serial_number = []byte{0x1, 0x00} // It's necessary to update the serial Number from the main entry
	header.length = []byte{0x00, 0x00}

	message_response := t.buildResponse(header, nil, nil)
	response_gps := Response_GPS{header.device_id, message_response}
	return response_gps
}
func (t *Micodus) parseData(message []byte) Location_report {
	var location Location_report
	// Extraer y convertir los valores
	latitude := float64(binary.BigEndian.Uint32(message[23:27])) / 1e6
	longitude := float64(binary.BigEndian.Uint32(message[27:31])) / 1e6 * -1
	altitude := float64(binary.BigEndian.Uint16(message[31:33]))

	location.Device_id = t.getDeviceId(message[5:11])
	location.Latitude = latitude
	location.Longitude = longitude
	location.Altitude = altitude
	return location
}
