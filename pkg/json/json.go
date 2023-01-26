package json

import (
	"encoding/json"
	"mina.fi/devopstuni/pkg/utils"
)

func UnMarshallRequest(body []byte, format any) {
	err := json.Unmarshal(body, &format)
	if err != nil {
		utils.FailOnError(err, "failed to unmarshall")
		return
	}
}
func MarshallRequest(data any) []byte {
	jsonData, err := json.Marshal(data)
	utils.FailOnError(err, "failed to convert to json")

	return jsonData
}
