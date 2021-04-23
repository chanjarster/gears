package dto

import "encoding/json"

// shortcut for type map[string]interface{}
type M map[string]interface{}

func UnmarshalM(jsonB []byte) (M, error) {
	var raw = make(M, 5)
	if err := json.Unmarshal(jsonB, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func MarshalM(m M) ([]byte, error) {
	return json.Marshal(m)
}
