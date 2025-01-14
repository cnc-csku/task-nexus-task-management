package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID                 bson.ObjectID       `bson:"_id" json:"id"`
	WorkspaceID        bson.ObjectID       `bson:"workspace_id" json:"workspaceId"`
	Name               string              `bson:"name" json:"name"`
	ProjectPrefix      string              `bson:"project_prefix" json:"projectPrefix"`
	Description        string              `bson:"description" json:"description"`
	Status             string              `bson:"status" json:"status"`
	Members            []Member            `bson:"members" json:"members"`
	Workflows          []Workflow          `bson:"workflows" json:"workflows"`
	AttributeTemplates []AttributeTemplate `bson:"attributes_templates" json:"attributesTemplates"`
	Roles              []string            `bson:"roles" json:"roles"`
	CreatedAt          time.Time           `bson:"created_at" json:"createdAt"`
	CreatedBy          bson.ObjectID       `bson:"created_by" json:"createdBy"`
	UpdatedAt          time.Time           `bson:"updated_at" json:"updatedAt"`
	UpdatedBy          bson.ObjectID       `bson:"updated_by" json:"updatedBy"`
}

type Member struct {
	UserID   bson.ObjectID `bson:"user_id" json:"userId"`
	FullName string        `bson:"full_name" json:"fullName"`
	Role     string        `bson:"role" json:"role"`
}

type Workflow struct {
	PreviousStatuses []string `bson:"previous_statuses" json:"previousStatuses"`
	Status           string   `bson:"status" json:"status"`
}

type AttributeTemplate struct {
	Name  string      `bson:"name" json:"name"`
	Type  string      `bson:"type" json:"type"`
	Value interface{} `bson:"value" json:"value"`
}

const (
	ProjectStatus_ACTIVE   = "ACTIVE"
	ProjectStatus_ARCHIVED = "ARCHIVED"
)
