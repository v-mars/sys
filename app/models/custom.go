package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type IntArray []int

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *IntArray) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %s", value)
	}
	var err  error
	result := IntArray{}
	if len(bytes) != 0{
		err = json.Unmarshal(bytes, &result)
	}
	*j = result
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j IntArray) Value() (driver.Value, error) {
	if &j == nil{
		return nil, nil
	}
	bytes, err := json.Marshal(j)
	return string(bytes), err
}


type IntNestArray [][]int

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *IntNestArray) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %s", value)
	}
	var err  error
	result := IntNestArray{}
	if len(bytes) != 0{
		err = json.Unmarshal(bytes, &result)
	}
	*j = result
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j IntNestArray) Value() (driver.Value, error) {
	if &j == nil{
		return nil, nil
	}
	bytes, err := json.Marshal(j)
	return string(bytes), err
}