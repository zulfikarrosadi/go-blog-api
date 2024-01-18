package group

import (
	"context"
	"database/sql"
)

type GroupRepository interface {
	CreateGroup(context.Context, *CreateGroupRequest) (int64, error)
	GetGroup(context.Context) []Group
	FindGroupById(context.Context, int64) (*Group, error)
	DeleteGroup(context.Context, int64) error
}

type GroupRepositoryImpl struct {
	*sql.DB
}
