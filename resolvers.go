package app

import (
	"errors"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"google.golang.org/appengine/datastore"
)

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

func queryUsers(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	query := datastore.NewQuery("User")
	var users []User
	if keys, err := query.GetAll(ctx, &users); err != nil {
		return users, err
	} else {
		for i, key := range keys {
			users[i].ID = strconv.FormatInt(key.IntID(), 10)
		}
	}
	return users, nil
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

func queryPostsByUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	query := datastore.NewQuery("Post")
	var posts []Post
	user, ok := params.Source.(*User)
	if ok {
		query.Filter("UserID =", user.ID)
		if keys, err := query.GetAll(ctx, &posts); err != nil {
			return posts, err
		} else {
			for i, key := range keys {
				posts[i].ID = strconv.FormatInt(key.IntID(), 10)
			}
		}
	}
	return posts, nil
}
