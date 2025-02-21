package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

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
