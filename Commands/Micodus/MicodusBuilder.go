package micodus

import (
	commands "trucker/Commands"
)

type MicodusBuilder struct {
	flags    []byte
	headers  []byte
	body     []byte
	checksum commands.Checksum
}

func NewMicodusBuilder() *MicodusBuilder {
	return &MicodusBuilder{}
}

func (t *MicodusBuilder) BuildBody(body []byte) {
	t.body = body
}
func (t *MicodusBuilder) BuildFlags(flag []byte) {
	t.flags = flag
}
func (t *MicodusBuilder) BuildHeaders(headers []byte) {
	t.headers = headers
}
func (t *MicodusBuilder) BuildChecksum(calculation commands.Checksum) {
	t.checksum = calculation
}
func (t *MicodusBuilder) GetResult() []byte {
	result := make([]byte, 0)
	result = append(result, t.flags...)
	result = append(result, t.headers...)
	result = append(result, t.body...)
	result = append(result, t.checksum(result[1:])...)
	result = append(result, t.flags...)
	return result
}
