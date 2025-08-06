package services

import (
	"context"

	"github.com/Hueter57/graphql-go-test/graph/model"
	"github.com/aarondl/sqlboiler/v4/boil"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	ListUsersByID(ctx context.Context, ids []string) ([]*model.User, error)
}

type RepositoryService interface {
	GetRepoByFullName(ctx context.Context, owner, name string) (*model.Repository, error)
	GetRepoByID(ctx context.Context, id string) (*model.Repository, error)
}

type IssueService interface {
	GetIssueByID(ctx context.Context, id string) (*model.Issue, error)
	GetIssueByRepoAndNumber(ctx context.Context, repoID string, number int) (*model.Issue, error)
	ListIssueInRepository(ctx context.Context, repoID string, after *string, before *string, first *int, last *int) (*model.IssueConnection, error)
}

type ProjectService interface {
	GetProjectByID(ctx context.Context, id string) (*model.ProjectV2, error)
}

type PullRequestService interface {
	GetPullRequestByID(ctx context.Context, id string) (*model.PullRequest, error)
}

type Services interface {
	UserService
	RepositoryService
	IssueService
	PullRequestService
	ProjectService
}

type services struct {
	*userService
	*repositoryService
	*issueService
	*projectService
	*pullRequestService
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		userService:       &userService{exec: exec},
		repositoryService: &repositoryService{exec: exec},
		issueService:      &issueService{exec: exec},
		projectService:    &projectService{exec: exec},
		pullRequestService: &pullRequestService{exec: exec},
	}
}
