package commands

import "encoding/json"

type BaseResponse struct {
	Device_id string
	Data      map[string]interface{}
}

func (br *BaseResponse) ToMap() (map[string]interface{}, error) {
	return br.Data, nil
}
func (br *BaseResponse) ToJSON() ([]byte, error) {
	return json.Marshal(br)
}
