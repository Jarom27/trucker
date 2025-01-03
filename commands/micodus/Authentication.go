package micodus

import (
	"encoding/binary"
	"trucker/commands"
)

type Authentication struct {
	Device_id    []byte
	SerialNumber int32
}

func (t *Authentication) Execute(message []byte) (commands.CommandResponse, error) {
	t.Device_id = message[5:11]
	br := &commands.BaseResponse{Data: make(map[string]interface{})}
	br.Data["gps"] = t.response_GPS()
	return br, nil
}
func (t *Authentication) response_GPS() []byte {
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
