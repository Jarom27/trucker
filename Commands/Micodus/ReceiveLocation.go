package micodus

import (
	"encoding/binary"
	"fmt"
	"trucker/commands"
)

type ReceiveLocation struct {
	Device_id    []byte
	SerialNumber int32
}

func (t *ReceiveLocation) Execute(message []byte) (commands.CommandResponse, error) {
	t.Device_id = message[5:11]
	br := &commands.BaseResponse{Data: make(map[string]interface{})}
	br.Data["gps"] = t.response_GPS()
	br.Data["location_report"] = t.response_Location(message)
	return br, nil
}
func (t *ReceiveLocation) response_GPS() []byte {
	messageBuilder := NewMicodusBuilder()

	message_type := []byte{0x82, 0x01}
	message_length := []byte{0x00, 0x00}
	message_serial := make([]byte, 2)
	binary.BigEndian.PutUint16(message_serial, uint16(t.SerialNumber))
	header := make([]byte, 0)

	header = append(header, message_type...)
	header = append(header, message_length...)
	header = append(header, t.Device_id...)
	header = append(header, message_serial...)

	messageBuilder.BuildHeaders(header)
	messageBuilder.BuildFlags([]byte{0x7e})
	messageBuilder.BuildChecksum(MicodusChecksum)

	return messageBuilder.GetResult()
}
func (t *ReceiveLocation) response_Location(message []byte) commands.LocationReport {
	latitude := float64(binary.BigEndian.Uint32(message[23:27])) / 1e6
	longitude := float64(binary.BigEndian.Uint32(message[27:31])) / 1e6 * -1
	altitude := float64(binary.BigEndian.Uint16(message[31:33]))

	location := commands.LocationReport{
		Device_id: fmt.Sprintf("%x", t.Device_id),
		Latitude:  latitude,
		Longitude: longitude,
		Altitude:  altitude,
	}
	return location
}
