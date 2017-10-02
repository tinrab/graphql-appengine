package app

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"google.golang.org/appengine/datastore"
)

type PostListResult struct {
	Nodes      []Post `json:"nodes"`
	TotalCount int    `json:"totalCount"`
}

func createUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	name, _ := params.Args["name"].(string)
	user := &User{Name: name}

	key := datastore.NewIncompleteKey(ctx, "User", nil)
	if generatedKey, err := datastore.Put(ctx, key, user); err != nil {
		return User{}, err
	} else {
		user.ID = strconv.FormatInt(generatedKey.IntID(), 10)
	}
	return user, nil
}

func queryUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	if strID, ok := params.Args["id"].(string); ok {
		id, err := strconv.ParseInt(strID, 10, 64)
		if err != nil {
			return nil, errors.New("Invalid id")
		}
		user := &User{ID: strID}
		key := datastore.NewKey(ctx, "User", "", id, nil)
		if err := datastore.Get(ctx, key, user); err != nil {
			return nil, errors.New("User not found")
		}
		return user, nil
	}
	return User{}, nil
}

func createPost(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	content, _ := params.Args["content"].(string)
	userID, _ := params.Args["userId"].(string)
	post := &Post{UserID: userID, Content: content, CreatedAt: time.Now().UTC()}

	key := datastore.NewIncompleteKey(ctx, "Post", nil)
	if generatedKey, err := datastore.Put(ctx, key, post); err != nil {
		return Post{}, err
	} else {
		post.ID = strconv.FormatInt(generatedKey.IntID(), 10)
	}
	return post, nil
}

func queryPosts(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	query := datastore.NewQuery("Post")
	if limit, ok := params.Args["limit"].(int); ok {
		query = query.Limit(limit)
	}
	if offset, ok := params.Args["offset"].(int); ok {
		query = query.Offset(offset)
	}
	return queryPostList(ctx, query)
}

func queryPostsByUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	query := datastore.NewQuery("Post")
	if limit, ok := params.Args["limit"].(int); ok {
		query = query.Limit(limit)
	}
	if offset, ok := params.Args["offset"].(int); ok {
		query = query.Offset(offset)
	}
	if user, ok := params.Source.(*User); ok {
		query = query.Filter("UserID =", user.ID)
	}
	return queryPostList(ctx, query)
}

func queryPostList(ctx context.Context, query *datastore.Query) (PostListResult, error) {
	var result PostListResult
	if keys, err := query.GetAll(ctx, &result.Nodes); err != nil {
		return result, err
	} else {
		for i, key := range keys {
			result.Nodes[i].ID = strconv.FormatInt(key.IntID(), 10)
		}
		result.TotalCount = len(result.Nodes)
	}
	return result, nil
}
