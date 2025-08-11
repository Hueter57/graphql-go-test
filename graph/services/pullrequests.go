package services

import (
	"context"
	"log"

	"github.com/Hueter57/graphql-go-test/graph/db"
	"github.com/Hueter57/graphql-go-test/graph/model"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type pullRequestService struct {
	exec boil.ContextExecutor
}

func (p *pullRequestService) GetPullRequestByID(ctx context.Context, id string) (*model.PullRequest, error) {
	pr, err := db.Pullrequests(
		qm.Select(
			db.PullrequestColumns.ID,
			db.PullrequestColumns.BaseRefName,
			db.PullrequestColumns.Closed,
			db.PullrequestColumns.HeadRefName,
			db.PullrequestColumns.URL,
			db.PullrequestColumns.Number,
			db.PullrequestColumns.Repository,
		),
		db.PullrequestWhere.ID.EQ(id),
	).One(ctx, p.exec)
	if err != nil {
		return nil, err
	}
	return convertPullRequest(pr), nil
}

func convertPullRequest(pr *db.Pullrequest) *model.PullRequest {
	prURL, err := model.UnmarshalURI(pr.URL)
	if err != nil {
		log.Println("invalid URI", pr.URL)
	}
	return &model.PullRequest{
		ID:          pr.ID,
		BaseRefName: pr.BaseRefName,
		Closed:      (pr.Closed == 1),
		HeadRefName: pr.HeadRefName,
		URL:         prURL,
		Number:      int(pr.Number),
		Repository:  &model.Repository{ID: pr.Repository},
	}
}
