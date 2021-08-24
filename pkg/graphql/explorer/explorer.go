package explorer

import (
	"fmt"
	"net/http"

	"github.com/rafaelrubbioli/fileapi/pkg/config"
)

func Handler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <link rel="shortcut icon" href="data:image/png;base64,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" />
    <style>
      body {
        height: 100%%;
        margin: 0;
        width: 100%%;
        overflow: hidden;
      }
      #graphiql {
        height: 100vh;
      }
    </style>
	<link href="https://unpkg.com/graphiql@1.0.6/graphiql.min.css" rel="stylesheet" />
    <script crossorigin src="https://unpkg.com/react/umd/react.production.min.js" ></script>
    <script crossorigin src="https://unpkg.com/react-dom/umd/react-dom.production.min.js"></script>
    <script crossorigin src="https://unpkg.com/graphiql@1.0.6/graphiql.min.js" ></script>
  </head>
  <body>
    <div id="graphiql">Loading...</div>
    <script>
      const search = window.location.search;
      const parameters = {};
      search.substr(1).split('&').forEach(function (entry) {
        const eq = entry.indexOf('=');
        if (eq >= 0) {
          parameters[decodeURIComponent(entry.slice(0, eq))] =
            decodeURIComponent(entry.slice(eq + 1));
        }
      });
      
      if (parameters.variables) {
        try {
          parameters.variables = JSON.stringify(JSON.parse(parameters.variables), null, 2);
        } catch (e) {}
      }
      
      function onEditQuery(newQuery) {
        parameters.query = newQuery;
        updateURL();
      }
      
      function onEditVariables(newVariables) {
        parameters.variables = newVariables;
        updateURL();
      }
      
      function onEditOperationName(newOperationName) {
        parameters.operationName = newOperationName;
        updateURL();
      }
      
      function updateURL() {
        const newSearch = '?' + Object.keys(parameters).filter(function (key) {
          return Boolean(parameters[key]);
        }).map(function (key) {
          return encodeURIComponent(key) + '=' +
            encodeURIComponent(parameters[key]);
        }).join('&');
        history.replaceState(null, null, newSearch);
      }
      
      function graphQLFetcher(graphQLParams, graphiQLOptions) {
        const { headers } = graphiQLOptions;
        return fetch('%s' + (window.bpc === 1 ? '?bpc=1' : ''), {
          method: 'post',
          headers: Object.assign(headers || {}, {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
		  }),
          body: JSON.stringify(graphQLParams),
          credentials: 'include',
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
        React.createElement(GraphiQL, {
          fetcher: graphQLFetcher,
          headerEditorEnabled: true,
          shouldPersistHeaders: true,
          query: parameters.query,
          variables: parameters.variables,
          operationName: parameters.operationName,
          onEditQuery: onEditQuery,
          onEditVariables: onEditVariables,
          onEditOperationName: onEditOperationName
        }),
        document.getElementById('graphiql')
      );
    </script>
  </body>
</html>
`, config.BaseURL())
}
