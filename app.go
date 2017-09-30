package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"google.golang.org/appengine/log"
)

var schema graphql.Schema

func init() {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	var err error
	schema, err = graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Errorf(nil, "%v", err)
	}

	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	//ctx := appengine.NewContext(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	params := graphql.Params{
		Schema:        schema,
		RequestString: string(body),
	}

	resp := graphql.Do(params)
	if len(resp.Errors) > 0 {
		responseError(w, fmt.Sprintf("%+v", resp.Errors), http.StatusBadRequest)
		return
	}

	responseJSON(w, resp)
}
