package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/graph-gophers/graphql-go"
)

// Service trying to fix something ??????
type Service interface {
	Hello() string
}

func main() {

	// Open project graphql schema
	f, err := os.Open("resource/schema.graphql")
	if err != nil {
		panic(err)
	}

	// Read stream
	b, err := ioutil.ReadAll(f)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	// Mock from?
	// s, err := mock.NewService()
	// if err != nil {
	// 	panic(err)
	// }

	// Parsed schema
	schema := graphql.MustParseSchema(
		string(b),
		func() interface{} {
			return "Hello World!"
		})

	routes := []Route{
		{
			Name:          "ServiceGraphQL",
			Method:        "POST",
			Pattern:       "/query",
			GzipMandatory: false,
			HandlerFunc:   SchemaHandler(schema),
		}, {
			Name:          "ServiceGraphiQL",
			Method:        "GET",
			Pattern:       "/",
			GzipMandatory: false,
			HandlerFunc:   GraphiQLHandler(),
		},
	}

	handler := NewRouter(routes)
	panic(http.ListenAndServe(":8080", handler))
}

// SchemaHandler dunno what this code do ???????
func SchemaHandler(schema *graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
		json.NewEncoder(w).Encode(response)
	}
}

// GraphiQLHandler dunno what this code do, I guess that its for return home HTML ??????
func GraphiQLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(graphiql)
	}
}

var graphiql = []byte(`
	<!DOCTYPE html>
	<html>
		<head>
			<link href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" rel="stylesheet" />
			<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
		</head>
		<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
			<div id="graphiql" style="height: 100vh;">Loading...</div>
			<script>
				function graphQLFetcher(graphQLParams) {
					return fetch("/query", {
						method: "post",
						body: JSON.stringify(graphQLParams),
						credentials: "include",
					}).then(function (response) {
						return response.text();
					}).then(function (responseBody) {
						try {
							return JSON.parse(responseBody);
						} catch (error) {
							return responseBody;
						}
					});
				}
				ReactDOM.render(
					React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
					document.getElementById("graphiql")
				);
			</script>
		</body>
	</html>
	`)
