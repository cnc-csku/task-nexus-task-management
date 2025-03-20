package mongo

import (
	"time"

	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type userFilter bson.M

func NewUserFilter() userFilter {
	return userFilter{}
}

func (f userFilter) WithUserID(userID bson.ObjectID) {
	f["_id"] = userID
}

func (f userFilter) WithEmail(email string) {
	f["email"] = email
}

func (f userFilter) WithUserIDs(userIDs []bson.ObjectID) {
	f["_id"] = bson.M{"$in": userIDs}
}

type userUpdate bson.M

func NewUserUpdate() userUpdate {
	return userUpdate{}
}

func (u userUpdate) UpdateProfile(in *repositories.UpdateUserProfileRequest) {
	u["$set"] = bson.M{
		"full_name":            in.FullName,
		"display_name":         in.DisplayName,
		"default_profile_url":  in.DefaultProfileUrl,
		"uploaded_profile_url": in.UploadedProfileUrl,
		"updated_at":           time.Now(),
		"upated_by":            in.UpdatedBy,
	}
}
