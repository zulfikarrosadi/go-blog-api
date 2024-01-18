package group

import "time"

type Group struct {
	Id             int64     `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateGroupRequest struct {
	Title          string    `json:"title" validate:"required"`
	Description    string    `json:"description"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}
