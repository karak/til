//go:generate gqlgen -schema ../schema.graphql -typemap types.json

package graph

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	"github.com/vvakame/til/graphql/try-go-gqlgen/models"
)

type MyApp struct {
	todos   []models.Todo
	UserMap map[string]models.UserImpl
}

func NewMyApp() *MyApp {
	return &MyApp{
		UserMap: make(map[string]models.UserImpl),
	}
}

func (a *MyApp) Query_node(ctx context.Context, id string) (Node, error) {
	switch id[0:1] {
	case "U":
		user, ok := a.UserMap[id]
		if !ok {
			return nil, nil
		}
		return user, nil
	case "T":
		for _, todo := range a.todos {
			if id == todo.ID {
				return todo, nil
			}
		}
		return nil, nil
	}

	return nil, errors.Errorf("unknown ID format: %s", id)
}

func (a *MyApp) Query_nodes(ctx context.Context, ids []string) ([]Node, error) {
	respList := make([]Node, 0, len(ids))
	for _, id := range ids {
		node, err := a.Query_node(ctx, id)
		if err != nil {
			return nil, err
		}
		respList = append(respList, node)
	}
	return respList, nil
}

func (a *MyApp) Query_todos(ctx context.Context) ([]models.Todo, error) {
	return a.todos, nil
}

func (a *MyApp) Query_searchTodo(ctx context.Context, id *string) ([]models.Todo, error) {
	if id != nil {
		for _, todo := range a.todos {
			if todo.ID == *id {
				return []models.Todo{todo}, nil
			}
		}

		return nil, errors.Errorf("id: %s is not exists", *id)
	}

	return nil, errors.New("query parameter is not specified")
}

func (a *MyApp) Mutation_createTodo(ctx context.Context, text string) (models.Todo, error) {
	user := models.UserImpl{
		ID:   fmt.Sprintf("U%d", rand.Int()),
		Name: fmt.Sprintf("Name of U%d", rand.Int()),
	}
	todo := models.Todo{
		Text:   text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		UserID: user.ID,
	}
	a.UserMap[user.ID] = user
	a.todos = append(a.todos, todo)
	return todo, nil
}

func (a *MyApp) Todo_user(ctx context.Context, obj *models.Todo) (models.UserImpl, error) {
	user, err := ctx.Value(models.UserLoaderKey).(*models.UserImplLoader).Load(obj.UserID)
	if err != nil {
		return models.UserImpl{}, err
	}
	return *user, nil
}
