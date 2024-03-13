package files

import "encoding/json"

type JSONEpiFile struct {
	EpiFile
}

func (j *JSONEpiFile) ParseJSON() (map[string]interface{}, error) {
	contents, err := j.ReadContents()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(contents, &data); err != nil {
		return nil, err
	}

	return data, nil
}
