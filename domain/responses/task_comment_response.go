package responses

import (
	"time"
)

type ListTaskCommentResponse struct {
	ID              string    `bson:"_id" json:"id"`
	Content         string    `bson:"content" json:"content"`
	UserID          string    `bson:"user_id" json:"userId"`
	UserDisplayName string    `bson:"user_display_name" json:"userDisplayName"`
	UserProfileUrl  string    `bson:"user_profile_url" json:"userProfileUrl"`
	TaskID          string    `bson:"task_id" json:"taskId"`
	CreatedAt       time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updatedAt"`
}
