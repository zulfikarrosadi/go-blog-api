package group

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zulfikarrosadi/go-blog-api/auth"
	"github.com/zulfikarrosadi/go-blog-api/lib"
)

type GroupRepository interface {
	CreateGroup(context.Context, *CreateGroupRequest) (int64, error)
	GetGroup(context.Context) []Group
	FindGroupById(context.Context, int64) (*Group, error)
}

type GroupRepositoryImpl struct {
	*sql.DB
}

func (gri *GroupRepositoryImpl) CreateGroup(ctx context.Context, data *CreateGroupRequest) (int64, error) {
	accessToken := ctx.Value("accessToken").(auth.AccessToken)
	q := "INSERT INTO groups (title, description, created_by) VALUE (?,?,?)"
	r, err := gri.DB.ExecContext(ctx, q, data.Title, data.Description, accessToken.UserId)
	if err != nil {
		lib.ValidateErrorV2("create_group_repo", err)
		return 0, errors.New("failed to create group, please try again")
	}

	i, _ := r.LastInsertId()
	return i, nil
}

func (gri *GroupRepositoryImpl) GetGroup(ctx context.Context) []Group {
	q := "SELECT id, title, description, profile_picture, created_at FROM groups"
	r, err := gri.QueryContext(ctx, q)
	if err != nil {
		lib.ValidateErrorV2("get_group_repo", err)
		return nil
	}

	groups := []Group{}
	for r.Next() {
		group := Group{}
		r.Scan(&group.Id, &group.Title, &group.Description, &group.ProfilePicture, &group.CreatedAt)
		groups = append(groups, group)
	}

	return groups
}

func (gri *GroupRepositoryImpl) FindGroupById(ctx context.Context, id int64) (*Group, error) {
	q := "SELECT id, title, description, profile_picture, created_at FROM groups WHERE id = ?"
	r := gri.QueryRowContext(ctx, q, id)

	group := &Group{}
	err := r.Scan(&group.Id, &group.Title, &group.Description, &group.ProfilePicture, &group.CreatedAt)
	if err != nil {
		lib.ValidateErrorV2("find_group_by_id_repo", err)
		return nil, err
	}

	return group, nil
}
