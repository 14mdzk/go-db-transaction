package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserActivityName string

const (
	UserActivityNameCreate UserActivityName = "create"
	UserActivityNameUpdate UserActivityName = "update"
	UserActivityNameDelete UserActivityName = "delete"
)

type UserActivity struct {
	ID          uuid.UUID
	Object      string
	ObjectID    *uuid.UUID
	Name        UserActivityName
	Description string
	CreatedAt   time.Time
}

func NewUserActivity(objectID *uuid.UUID, activityName UserActivityName, description string) *UserActivity {
	return &UserActivity{
		Object:      "user",
		ObjectID:    objectID,
		Name:        activityName,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
