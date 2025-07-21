package graph

import (
	"context"
	"errors"

	"github.com/Hueter57/graphql-go-test/graph/model"
	"github.com/Hueter57/graphql-go-test/graph/services"
	"github.com/graph-gophers/dataloader/v7"
)

type Loaders struct {
	UserLoader dataloader.Interface[string, *model.User]
}

func NewLoaders(Srv services.Services) *Loaders {
	userBatcher := &userBatcher{Srv: Srv}

	return &Loaders{
		// dataloader.Loader[string, *model.User]構造体型をセットするために、
		// dataloader.NewBatchedLoader関数を呼び出す
		UserLoader: dataloader.NewBatchedLoader(userBatcher.BatchGetUsers),
	}
}

type userBatcher struct {
	Srv services.Services
}

func (u *userBatcher) BatchGetUsers(ctx context.Context, ids []string) []*dataloader.Result[*model.User] {
	// 引数と戻り値のスライスlenは等しくする
	results := make([]*dataloader.Result[*model.User], len(ids))
	for i := range results {
		results[i] = &dataloader.Result[*model.User]{
			Error: errors.New("not found"),
		}
	}

	// 検索条件であるIDが、引数でもらったIDsスライスの何番目のインデックスに格納されていたのか検索できるようにmap化する
	indexs := make(map[string]int, len(ids))
	for i, id := range ids {
		indexs[id] = i
	}

	// サービス層のメソッドを使い、指定されたIDを持つユーザーを全て取得する
	// (ListUsersByIDメソッド内では、IN句を用いたselect文が実行されている)
	users, err := u.Srv.ListUsersByID(ctx, ids)

	// 取得結果を、戻り値resultの中の適切な場所に格納する
	for _, user := range users {
		var rsl *dataloader.Result[*model.User]
		if err != nil {
			rsl = &dataloader.Result[*model.User]{
				Error: err,
			}
		} else {
			rsl = &dataloader.Result[*model.User]{
				Data: user,
			}
		}
		results[indexs[user.ID]] = rsl
	}
	return results
}
