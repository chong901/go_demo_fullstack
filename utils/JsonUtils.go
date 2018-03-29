package utils

import "encoding/json"

func ConvertToJsonMap(data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	var jsonMap map[string]interface{}

	json.Unmarshal(jsonData, &jsonMap)

	return jsonMap, nil
}
