package services

import (
	"context"
	"log"

	"github.com/Hueter57/graphql-go-test/graph/db"
	"github.com/Hueter57/graphql-go-test/graph/model"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type projectService struct {
	exec boil.ContextExecutor
}

func (i *projectService) GetProjectByID(ctx context.Context, id string) (*model.ProjectV2, error) {
	project, err := db.Projects(
		qm.Select(
			db.ProjectColumns.ID,
			db.ProjectColumns.Title,
			db.ProjectColumns.URL,
			db.ProjectColumns.Number,
			db.ProjectColumns.Owner,
		),
		db.ProjectWhere.ID.EQ(id),
	).One(ctx, i.exec)
	if err != nil {
		return nil, err
	}
	return convertProject(project), nil
}

func convertProject(project *db.Project) *model.ProjectV2 {
	projectURL, err := model.UnmarshalURI(project.URL)
	if err != nil {
		log.Println("invalid URI", project.URL)
	}
	return &model.ProjectV2{
		ID:     project.ID,
		Title:  project.Title,
		URL:    projectURL,
		Number: int(project.Number),
		Owner:  &model.User{ID: project.Owner},
	}
}
