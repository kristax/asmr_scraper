package model

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type DataCache struct {
	Id     uint           `gorm:"primaryKey"`
	Target string         `gorm:"uniqueIndex:idx_code"`
	Code   string         `gorm:"uniqueIndex:idx_code"`
	Data   datatypes.JSON `gorm:"type:json"`
	Client string
}

func NewDataCache(client, target, code string, data any) *DataCache {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return &DataCache{
		Id:     0,
		Client: client,
		Target: target,
		Code:   code,
		Data:   bytes,
	}
}
