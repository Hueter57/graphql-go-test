package main

import (
	// (一部抜粋)
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/Hueter57/graphql-go-test/graph"
	"github.com/Hueter57/graphql-go-test/graph/model"
	"github.com/Hueter57/graphql-go-test/graph/resolver"
	mockservices "github.com/Hueter57/graphql-go-test/graph/services/mock"
	"github.com/Hueter57/graphql-go-test/internal"
	"github.com/tenntenn/golden"
	"go.uber.org/mock/gomock"
)

var (
	flagUpdate bool
	goldenDir  string = "./testdata/golden/"
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func getRequestBody(t *testing.T, testdata, name string) io.Reader {
	t.Helper()

	queryBody, err := os.ReadFile(testdata + name + ".golden")
	if err != nil {
		t.Fatal(err)
	}
	query := struct{ Query string }{
		string(queryBody),
	}
	reqBody := bytes.Buffer{}
	if err := json.NewEncoder(&reqBody).Encode(&query); err != nil {
		t.Fatal("error encode", err)
	}
	return &reqBody
}

func getResponseBody(t *testing.T, res *http.Response) string {
	t.Helper()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("error read body", err)
	}
	var got bytes.Buffer
	if err := json.Indent(&got, raw, "", "\t"); err != nil {
		t.Fatal("json.Indent", err)
	}
	return got.String()
}

func TestNodeRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	repoID := "REPO_1"
	ownerID := "U_1"
	sm := mockservices.NewMockServices(ctrl)
	sm.EXPECT().GetRepoByID(gomock.Any(), repoID).Return(&model.Repository{
		ID:        repoID,
		Owner:     &model.User{ID: ownerID},
		Name:      "repo1",
		CreatedAt: time.Date(2022, 12, 30, 0, 12, 21, 0, time.UTC),
	}, nil)
	sm.EXPECT().GetUserByID(gomock.Any(), ownerID).Return(&model.User{
		ID:   ownerID,
		Name: "hsaki",
	}, nil)

	h := handler.New(internal.NewExecutableSchema(internal.Config{Resolvers: &resolver.Resolver{
		Srv:     sm,
		Loaders: graph.NewLoaders(sm),
	}}))
	h.AddTransport(transport.POST{})

	srv := httptest.NewServer(h)
	t.Cleanup(func() { srv.Close() })

	reqBody := getRequestBody(t, goldenDir, t.Name()+"In.gql")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL, reqBody)
	if err != nil {
		t.Fatal("error new request", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("error request", err)
	}
	t.Cleanup(func() { res.Body.Close() })

	got := getResponseBody(t, res)
	if diff := golden.Check(t, flagUpdate, goldenDir, t.Name()+"Out.json", got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
