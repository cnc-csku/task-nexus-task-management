package mongo

import (
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type globalSettingFilter bson.M

func NewGlobalSettingFilter() globalSettingFilter {
	return globalSettingFilter{}
}

func (f globalSettingFilter) WithKey(key string) {
	f["key"] = key
}

type globalSettingUpdate bson.M

func NewGlobalSettingUpdate() globalSettingUpdate {
	return globalSettingUpdate{}
}

func (u globalSettingUpdate) WithSetting(setting *models.KeyValuePair) {
	u["$set"] = setting
}
