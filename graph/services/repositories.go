package services

import (
	"context"

	"github.com/Hueter57/graphql-go-test/graph/db"
	"github.com/Hueter57/graphql-go-test/graph/model"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type repositoryService struct {
	exec boil.ContextExecutor
}

func convertRepository(repo *db.Repository) *model.Repository {
	return &model.Repository{
		ID:        repo.ID,
		Owner:     &model.User{ID: repo.Owner},
		Name:      repo.Name,
		CreatedAt: repo.CreatedAt,
	}
}

func (r *repositoryService) GetRepoByFullName(ctx context.Context, owner, name string) (*model.Repository, error) {
	repo, err := db.Repositories(
		qm.Select(
			db.RepositoryColumns.ID,        // レポジトリID
			db.RepositoryColumns.Name,      // レポジトリ名
			db.RepositoryColumns.Owner,     // レポジトリを所有しているユーザーのID
			db.RepositoryColumns.CreatedAt, // 作成日時
		),
		db.RepositoryWhere.Owner.EQ(owner),
		db.RepositoryWhere.Name.EQ(name),
	).One(ctx, r.exec)
	if err != nil {
		return nil, err
	}
	return convertRepository(repo), nil
}

func (r *repositoryService) GetRepoByID(ctx context.Context, id string) (*model.Repository, error) {
	repo, err := db.Repositories(
		qm.Select(
			db.RepositoryColumns.ID,        // レポジトリID
			db.RepositoryColumns.Name,      // レポジトリ名
			db.RepositoryColumns.Owner,     // レポジトリを所有しているユーザーのID
			db.RepositoryColumns.CreatedAt, // 作成日時
		),
		db.RepositoryWhere.ID.EQ(id),
	).One(ctx, r.exec)
	if err != nil {
		return nil, err
	}
	return convertRepository(repo), nil
}
