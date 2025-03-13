package models

type KeyValuePair struct {
	Key   string           `json:"key" bson:"key"`
	Type  KeyValuePairType `json:"type" bson:"type"`
	Value interface{}      `json:"value" bson:"value"`
}

type KeyValuePairType string

const (
	KeyValuePairTypeString  KeyValuePairType = "STRING"
	KeyValuePairTypeNumber  KeyValuePairType = "NUMBER"
	KeyValuePairTypeBoolean KeyValuePairType = "BOOLEAN"
	KeyValuePairTypeDate    KeyValuePairType = "DATE"
)

func (k KeyValuePairType) String() string {
	return string(k)
}

func (k KeyValuePairType) IsValid() bool {
	switch k {
	case KeyValuePairTypeString, KeyValuePairTypeNumber, KeyValuePairTypeBoolean, KeyValuePairTypeDate:
		return true
	}
	return false
}
