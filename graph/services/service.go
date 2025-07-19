package services

import (
	"context"

	"github.com/Hueter57/graphql-go-test/graph/model"
	"github.com/aarondl/sqlboiler/v4/boil"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByName(ctx context.Context, name string) (*model.User, error)
}

type RepositoryService interface {
	GetRepoByFullName(ctx context.Context, owner, name string) (*model.Repository, error)
	GetRepoByID(ctx context.Context, id string) (*model.Repository, error)
}

type IssueService interface {
	GetIssueByRepoAndNumber(ctx context.Context, repoID string, number int) (*model.Issue, error)
	ListIssueInRepository(ctx context.Context, repoID string, after *string, before *string, first *int, last *int) (*model.IssueConnection, error)
}

type Services interface {
	UserService
	RepositoryService
	IssueService
}

type services struct {
	*userService
	*repositoryService
	*issueService
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		userService:       &userService{exec: exec},
		repositoryService: &repositoryService{exec: exec},
		issueService:      &issueService{exec: exec},
	}
}
