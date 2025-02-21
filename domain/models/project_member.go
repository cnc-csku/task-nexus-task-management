package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectMember struct {
	ID        bson.ObjectID     `bson:"_id" json:"id"`
	UserID    bson.ObjectID     `bson:"user_id" json:"userId"`
	ProjectID bson.ObjectID     `bson:"project_id" json:"projectId"`
	Role      ProjectMemberRole `bson:"role" json:"role"`
	Position  string            `bson:"position" json:"position"`
	JoinedAt  time.Time         `bson:"joined_at" json:"joinedAt"`
	RemovedAt *time.Time        `bson:"removed_at" json:"removedAt"`
}

type ProjectMemberRole string

const (
	ProjectMemberRoleOwner     ProjectMemberRole = "OWNER"
	ProjectMemberRoleModerator ProjectMemberRole = "MODERATOR"
	ProjectMemberRoleMember    ProjectMemberRole = "MEMBER"
)

func (p ProjectMemberRole) String() string {
	return string(p)
}

func (p ProjectMemberRole) IsValid() bool {
	switch p {
	case ProjectMemberRoleOwner, ProjectMemberRoleModerator, ProjectMemberRoleMember:
		return true
	}
	return false
}
