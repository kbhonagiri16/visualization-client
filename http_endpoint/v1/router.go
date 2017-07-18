package v1Api

import (
	"fmt"
	"github.com/pressly/chi"
	"net/http"

	"github.com/kbhonagiri16/visualization-client/http_endpoint/authentication"
	"github.com/kbhonagiri16/visualization-client/http_endpoint/common"
	v1handlers "github.com/kbhonagiri16/visualization-client/http_endpoint/v1/handlers"
	"github.com/kbhonagiri16/visualization-client/logging"
)

const adminAPIPrefix = "/admin"
const authPrefix = "/auth"

func authRouter(clients *common.ClientContainer,
	handler common.HandlerInterface, secret string) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/openstack", func(w http.ResponseWriter, r *http.Request) {
		// expected HEADER token name
		const openstackHeaderName = "X-OpenStack-Auth-Token"
		// response to be returned with 401 code to user
		const errorMsg = "Unauthorized. Token is invalid or expired."

		// check if openstack token is in request HEADERs
		openstackToken := r.Header.Get(openstackHeaderName)
		if openstackToken == "" {
			common.WriteErrorToResponse(w, http.StatusUnauthorized, errorMsg,
				fmt.Sprintf("header %s is not specified", openstackHeaderName))
			return
		}

		// try to authenticate with provided token
		token, err := handler.AuthOpenstack(clients, &common.RealClock{},
			openstackToken, secret)
		if err != nil {
			switch err.(type) {
			// common.InvalidOpenstackToken means, that user provided invalid
			// token. We have to return 401 error then
			case common.InvalidOpenstackToken:
				common.WriteErrorToResponse(w, http.StatusUnauthorized,
					errorMsg, err.Error())
				return
			// If any other error happened -> return 500 error
			default:
				log.Logger.Error(err)
				common.WriteErrorToResponse(w, http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
					"Internal server error occured")
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write(token)
	})
	return router
}

func adminRouter(clients *common.ClientContainer,
	authMiddleware func(http.Handler) http.Handler,
	handler common.HandlerInterface, secret string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Use(httpAuth.AdminAuthenticationMiddleware(secret))

	// routes for users
	r.Route("/users", func(r chi.Router) {
		// Get users list
		r.Get("/", v1handlers.GetUsers(clients, handler))

		// Get users  by id
		r.Get("/{userID}", v1handlers.GetUsersID(clients, handler))

		// Delete users by id
		r.Delete("/{userID}", v1handlers.DeleteUser(clients, handler))

		// Create user
		r.Post("/", v1handlers.CreateUser(clients, handler))
	})

	// routes for organizations
	r.Route("/organizations", func(r chi.Router) {
		// Get organizations list
		r.Get("/", v1handlers.GetOrganization(clients, handler))

		// Get organization by id
		r.Get("/{organizationID}", v1handlers.GetOrganizationID(clients, handler))

		// Delete organizations by id
		r.Delete("/{organizationID}", v1handlers.DeleteOrganization(clients, handler))

		// Create organization
		r.Post("/", v1handlers.CreateOrganization(clients, handler))

		// Delete user in organizations by id
		r.Delete("/{organizationID}/users/{userID}", v1handlers.DeleteOrganizationUser(clients, handler))

		// Get users in organization
		r.Get("/{organizationID}/users", v1handlers.GetOrganizationUser(clients, handler))

		// Post create user in organization
		r.Post("/{organizationID}/users", v1handlers.CreateOrganizationUser(clients, handler))
	})
	return r
}

func visualizationRouter(clients *common.ClientContainer,
	handler common.HandlerInterface,
	authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	// temporary commented to simplify testing
	router.Use(authMiddleware)
	router.Get("/visualizations", v1handlers.VisualizationsGet(
		clients, handler))
	router.Post("/visualizations", v1handlers.VisualizationsPost(
		clients, handler))
	router.Delete("/visualization/{visualizationID}", v1handlers.VisualizationDelete(
		clients, handler))
	return router
}

func datasourceRouter(clients *common.ClientContainer,
	handler common.HandlerInterface,
	authMiddleware func(http.Handler) http.Handler) *chi.Mux {

	router := chi.NewRouter()
	router.Use(authMiddleware)

	// route for datasources
	router.Get("/", v1handlers.GetDatasources(
		clients, handler))
	router.Get("/{datasourceID}", v1handlers.GetDatasourceID(
		clients, handler))
	router.Post("/", v1handlers.CreateDatasource(
		clients, handler))
	router.Delete("/{datasourceID}", v1handlers.DeleteDatasource(
		clients, handler))

	return router
}

// InitializeRouter initializes /v1 routers
func InitializeRouter(clients *common.ClientContainer,
	handler common.HandlerInterface, secret string) *chi.Mux {
	router := chi.NewRouter()
	authMiddleware := httpAuth.AuthenticationMiddleware(secret)
	router.Mount(adminAPIPrefix, adminRouter(clients, authMiddleware, handler, secret))
	router.Mount(authPrefix, authRouter(clients, handler, secret))
	router.Mount("/", visualizationRouter(clients, handler, authMiddleware))
	router.Mount("/datasources", datasourceRouter(clients, handler, authMiddleware))
	return router
}
