package endpoint

import (
	"fmt"
	"github.com/pressly/chi"
	"net/http"

	"visualization-client/http_endpoint/common"
	"visualization-client/http_endpoint/v1"
)

const v1ApiPrefix = "/v1"

// InitializeRouter would initialize all routers of our api
func InitializeRouter(clients *common.ClientContainer,
	handler common.HandlerInterface, secret string) *chi.Mux {
	rootRouter := chi.NewRouter()
	rootRouter.Mount(v1ApiPrefix, v1Api.InitializeRouter(clients, handler,
		secret))
	return rootRouter
}

// Serve is an entry point to our HTTP API
func Serve(secret string, httpPort int, clients *common.ClientContainer) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), InitializeRouter(
		clients, &v1Api.V1Handler{}, secret))
}
