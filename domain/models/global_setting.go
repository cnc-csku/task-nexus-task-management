package models

type GlobalSetting struct {
	Key   string            `json:"key" bson:"key"`
	Type  GlobalSettingType `json:"type" bson:"type"`
	Value interface{}       `json:"value" bson:"value"`
}

type GlobalSettingType string

const (
	GlobalSettingTypeString GlobalSettingType = "STRING"
	GlobalSettingTypeInt    GlobalSettingType = "NUMBER"
	GlobalSettingTypeBool   GlobalSettingType = "BOOLEAN"
)

