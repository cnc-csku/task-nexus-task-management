package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type userFilter bson.M

func NewUserFilter() userFilter {
	return userFilter{}
}

func (f userFilter) WithEmail(email string) {
	f["email"] = email
}
