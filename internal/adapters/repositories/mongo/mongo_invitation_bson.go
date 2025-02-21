package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type invitationFilter bson.M

func NewInvitationFilter() invitationFilter {
	return invitationFilter{}
}

func (f invitationFilter) WithID(id bson.ObjectID) {
	f["_id"] = id
}

func (f invitationFilter) WithWorkspaceID(workspaceID bson.ObjectID) {
	f["workspace_id"] = workspaceID
}

func (f invitationFilter) WithInviteeUserID(inviteeUserID bson.ObjectID) {
	f["invitee_user_id"] = inviteeUserID
}

func (f invitationFilter) WithNotExpired() {
	f["expired_at"] = bson.M{"$gt": time.Now()}
}

func (f invitationFilter) WithNotResponded() {
	f["responded_at"] = nil
}
