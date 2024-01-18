package group

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zulfikarrosadi/go-blog-api/lib"
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

func (gri *GroupRepositoryImpl) CreateGroup(ctx context.Context, data *CreateGroupRequest) (int64, error) {
	q := "INSERT INTO groups (title, description) VALUE (?,?)"
	r, err := gri.DB.ExecContext(ctx, q, data.Title, data.Description)
	if err != nil {
		lib.ValidateErrorV2("create_group_repo", err)
		return 0, errors.New("failed to create group, please try again")
	}

	i, _ := r.LastInsertId()
	return i, nil
}
