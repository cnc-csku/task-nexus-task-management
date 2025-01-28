package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Workspace struct {
	ID        bson.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name      string            `bson:"name" json:"name"`
	Members   []WorkspaceMember `bson:"members" json:"members"`
	CreatedBy bson.ObjectID     `bson:"created_by" json:"createdBy"`
	CreatedAt time.Time         `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updatedAt"`
}

type WorkspaceMember struct {
	UserID    bson.ObjectID       `bson:"user_id" json:"userId"`
	Role      WorkspaceMemberRole `bson:"role" json:"role"`
	JoinedAt  time.Time           `bson:"joined_at" json:"joinedAt"`
	RemovedAt *time.Time          `bson:"removed_at,omitempty" json:"removedAt,omitempty"`
}

type WorkspaceMemberRole string

const (
	WorkspaceMemberRoleOwner     WorkspaceMemberRole = "OWNER"
	WorkspaceMemberRoleModerator WorkspaceMemberRole = "MODERATOR"
	WorkspaceMemberRoleMember    WorkspaceMemberRole = "MEMBER"
)

func (w WorkspaceMemberRole) String() string {
	return string(w)
}

func (w WorkspaceMemberRole) IsValid() bool {
	switch w {
	case WorkspaceMemberRoleOwner, WorkspaceMemberRoleModerator, WorkspaceMemberRoleMember:
		return true
	}
	return false
}
