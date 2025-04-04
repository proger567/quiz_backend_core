package dto

import (
	"encoding/json"
	"fmt"
)

//TODO Move this type to pkg?

type Int64Array []int64

func (ia Int64Array) MarshalJSON() ([]byte, error) {
	stringArray := make([]string, len(ia))
	for i, v := range ia {
		stringArray[i] = fmt.Sprintf("%d", v) // Преобразуем int64 в string
	}
	return json.Marshal(stringArray)
}

func (ia *Int64Array) UnmarshalJSON(data []byte) error {
	var stringArray []string
	if err := json.Unmarshal(data, &stringArray); err != nil {
		return err
	}

	intArray := make([]int64, len(stringArray))
	for i, s := range stringArray {
		var v int64
		if _, err := fmt.Sscan(s, &v); err != nil { // Преобразуем string в int64
			return err
		}
		intArray[i] = v
	}
	*ia = intArray
	return nil
}
