package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"google.golang.org/appengine"
)

var schema graphql.Schema
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.String},
		"name": &graphql.Field{Type: graphql.String},
		"posts": &graphql.Field{
			Type:    postListType,
			Resolve: queryPostsByUser,
			Args: graphql.FieldConfigArgument{
				"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
				"offset": &graphql.ArgumentConfig{Type: graphql.Int},
			},
		},
	},
})
var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"userId":    &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.DateTime},
		"content":   &graphql.Field{Type: graphql.String},
	},
})
var postListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PostList",
	Fields: graphql.Fields{
		"nodes":      &graphql.Field{Type: graphql.NewList(postType)},
		"totalCount": &graphql.Field{Type: graphql.Int},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: createUser,
		},
		"createPost": &graphql.Field{
			Type: postType,
			Args: graphql.FieldConfigArgument{
				"userId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: createPost,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: queryUser,
		},
	},
})

func init() {
	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resp := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: string(body),
		Context:       ctx,
	})
	if len(resp.Errors) > 0 {
		responseError(w, fmt.Sprintf("%+v", resp.Errors), http.StatusBadRequest)
		return
	}
	responseJSON(w, resp)
}
