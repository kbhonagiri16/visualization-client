package v1handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/kbhonagiri16/visualization"
	"github.com/kbhonagiri16/visualization/http_endpoint/common"
	"github.com/kbhonagiri16/visualization/logging"
	"github.com/pressly/chi"
)

var emailValid = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// V1UsersOrgs implements part of handler interface
type V1UsersOrgs struct{}

// orgUser struct for create organization user
type orgUser struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Password string `json:"password" binding:"Required"`
}

// User struct for create user
type User struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"Required"`
}

// DataSource struct for create datasource
type DataSource struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	Access            string `json:"access"`
	URL               string `json:"url"`
	Password          string `json:"password"`
	User              string `json:"user"`
	Database          string `json:"database"`
	BasicAuthUser     string `json:"basicAuthUser"`
	BasicAuthPassword string `json:"basicAuthPassword"`
	BasicAuth         bool   `json:"basicAuth"`
	IsDefault         bool   `json:"isDefault"`
}

// LoginErrorCheck handles errors for login
func LoginErrorCheck(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Logger.Error(err)
		common.WriteErrorToResponse(w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			"Internal server error occured")
		return
	}
}

func helperCreateDatasource(clients *common.ClientContainer, handler common.HandlerInterface, w http.ResponseWriter, datasource DataSource) error {
	res, err := json.Marshal(datasource)
	if err != nil {
		common.WriteErrorToResponse(w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error())
		return err

	}

	// Craete user if no errors
	err = handler.CreateDatasource(clients, res)
	if err != nil {
		switch err.(type) {
		// grafanaclient.Exists means, that user provided details
		// of user which already exists. We return 409
		case visualization.Exists:
			errMsg := fmt.Sprintf("Datasource Exists")
			common.WriteErrorToResponse(w, http.StatusConflict,
				errMsg, err.Error())
			return err
		// If any other error happened -> return 500 error
		default:
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return err
		}
	}

	return err

}

func helperOrgUser(clients *common.ClientContainer, handler common.HandlerInterface, w http.ResponseWriter, OrgID int, user orgUser) error {

	res, err := json.Marshal(user)
	if err != nil {
		common.WriteErrorToResponse(w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error())
		return err

	}
	// Check if organization ID exists
	_, err = handler.GetOrganizationID(clients, OrgID)
	if err != nil {
		switch err.(type) {
		// common.OrganizationNotFound  means, that user provided the
		// ID of non existent user. We return 404
		case visualization.NotFound:
			errMsg := fmt.Sprintf("Organization Not Found")
			common.WriteErrorToResponse(w, http.StatusNotFound,
				errMsg, err.Error())
			return err
		// If any other error happened -> return 500 error
		default:
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return err
		}
	}

	// Create org user if no error
	err = handler.CreateOrganizationUser(clients, OrgID, res)
	if err != nil {
		switch err.(type) {
		// grafanaclient.Exists means, that user provided
		// user that already exists. We return 409
		case visualization.Exists:
			errMsg := fmt.Sprintf("User Exists")
			common.WriteErrorToResponse(w, http.StatusConflict,
				errMsg, err.Error())
			return err
		// If any other error happened -> return 500 error
		default:
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return err
		}
	}
	return err
}

func helperCreateUser(clients *common.ClientContainer, handler common.HandlerInterface, w http.ResponseWriter, user User) error {
	res, err := json.Marshal(user)
	if err != nil {
		common.WriteErrorToResponse(w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error())
		return err

	}

	// Craete user if no errors
	err = handler.CreateUser(clients, res)
	if err != nil {
		switch err.(type) {
		// grafanaclient.Exists means, that user provided details
		// of user which already exists. We return 409
		case visualization.Exists:
			errMsg := fmt.Sprintf("User Exists")
			common.WriteErrorToResponse(w, http.StatusConflict,
				errMsg, err.Error())
			return err
		// If any other error happened -> return 500 error
		default:
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return err
		}
	}

	return err

}

// GetUsers get the list of users
func GetUsers(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := handler.GetUsers(clients)
		if err != nil {
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return
		}
		w.Write(users)
	}
}

// GetUsersID gets the user details by ID
func GetUsersID(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		ID, err := strconv.Atoi(userID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"ID provided is not integer")
			return
		}
		userlist, err := handler.GetUserID(clients, ID)
		if err != nil {
			switch err.(type) {
			// grafanaclient.NotFound  means, that user provided the
			// ID of non existent user. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("User Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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
		w.Write(userlist)
	}
}

// DeleteUser method deletes a user by ID
func DeleteUser(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		ID, err := strconv.Atoi(userID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"userID provided is not integer")
			return
		}

		// check if the ID exists
		_, err = handler.GetUserID(clients, ID)
		if err != nil {
			switch err.(type) {
			// grafanaclient.NotFound  means, that user provided the
			// ID of non existent user. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("User Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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

		// if ID exists then delete that user
		err = handler.DeleteUser(clients, ID)
		if err != nil {
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return
		}
	}
}

// CreateUser method creates a user
func CreateUser(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				err.Error())
			return
		}
		if len(user.Email) == 0 {
			if len(user.Name) == 0 {
				common.WriteErrorToResponse(w, http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
					"provide Either Name or Email in parameters")
				return
			}
		}
		if len(user.Email) != 0 {
			if !(emailValid.MatchString(user.Email)) {
				common.WriteErrorToResponse(w, http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
					"Email Invalid")
				return
			}
		}
		if len(user.Login) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Login in parameters")
			return
		}
		if len(user.Password) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Password in parameters")
			return
		}

		// Create user if no error
		helperCreateUser(clients, handler, w, user)
	}
}

// GetOrganization method gets the list of organizations
func GetOrganization(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orglist, err := handler.GetOrganizations(clients)
		if err != nil {
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return
		}
		w.Write(orglist)
	}
}

// GetOrganizationID method gets the organization by ID
func GetOrganizationID(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID := chi.URLParam(r, "organizationID")
		ID, err := strconv.Atoi(organizationID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"provided organizationID is not integer")
			return
		}
		orglist, err := handler.GetOrganizationID(clients, ID)
		if err != nil {
			switch err.(type) {
			// common.OrganizationNotFound  means, that user provided the
			// ID of non existent user. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("Organization Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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

		w.Write(orglist)
	}
}

// DeleteOrganization method deletes the organzation
func DeleteOrganization(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID := chi.URLParam(r, "organizationID")
		ID, err := strconv.Atoi(organizationID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"provided organizationID is not integer")
			return
		}

		// check if the ID exists
		_, err = handler.GetOrganizationID(clients, ID)
		if err != nil {
			switch err.(type) {
			// common.OrganizationNotFound  means, that user provided the
			// ID of non existent user. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("Organization Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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

		// delete the organization if ID exists
		err = handler.DeleteOrganization(clients, ID)
		if err != nil {
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return
		}
	}
}

// CreateOrganization method creates the organization
func CreateOrganization(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var org struct {
			Name string `json:"name"`
		}

		err := json.NewDecoder(r.Body).Decode(&org)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				err.Error())
			return
		}

		if len(org.Name) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Name in parameters")
			return
		}

		res, err := json.Marshal(org)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				err.Error())
			return

		}

		// Create Organization if no error
		err = handler.CreateOrganization(clients, res)
		if err != nil {
			switch err.(type) {
			// grafanaclient.Exists means, that user provided
			// organization already exists. We return 409
			case visualization.Exists:
				errMsg := fmt.Sprintf("Organization Exists")
				common.WriteErrorToResponse(w, http.StatusConflict,
					errMsg, err.Error())
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
	}
}

// CreateOrganizationUser method creates the organization
func CreateOrganizationUser(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID := chi.URLParam(r, "organizationID")
		OrgID, err := strconv.Atoi(organizationID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"provided organizationID is not integer")
			return
		}
		user := orgUser{}

		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				err.Error())
			return
		}
		if len(user.Email) == 0 {
			if len(user.Name) == 0 {
				common.WriteErrorToResponse(w, http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
					"provide Either Name or Email in parameters")
				return
			}
		}
		if len(user.Email) != 0 {
			if !(emailValid.MatchString(user.Email)) {
				common.WriteErrorToResponse(w, http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
					"Email Invalid")
				return
			}
		}
		if len(user.Login) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Login in parameters")
			return
		}
		if len(user.Password) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Password in parameters")
			return
		}
		if len(user.Role) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Role in parameters")
			return
		}

		helperOrgUser(clients, handler, w, OrgID, user)
	}
}

// DeleteOrganizationUser method deletes a user in an organization
func DeleteOrganizationUser(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		ID, err := strconv.Atoi(userID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"provided userID is not integer")
			return
		}
		orgID := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.Atoi(orgID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"provided orgID is not integer")
			return
		}

		_, err = handler.GetUserID(clients, ID)
		if err != nil {
			switch err.(type) {
			// grafanaclient.NotFound  means, that user provided the
			// ID of non existent user. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("User Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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

		_, err = handler.GetOrganizationID(clients, organizationID)
		if err != nil {
			switch err.(type) {
			// grafanaclient.NotFound  means, that user provided the
			// ID of non existent organization. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("Org Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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

		err = handler.DeleteOrganizationUser(clients, ID, organizationID)
		if err != nil {
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return
		}
	}
}

//GetOrganizationUser gets the user by organization
func GetOrganizationUser(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID := chi.URLParam(r, "organizationID")
		ID, err := strconv.Atoi(organizationID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"provided organizationID is not integer")
			return
		}
		_, err = handler.GetOrganizationID(clients, ID)
		if err != nil {
			switch err.(type) {
			// common.OrganizationNotFound  means, that user provided the
			// ID of non existent user. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("Organization Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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

		orglist, err := handler.GetOrganizationUsers(clients, ID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return
		}
		w.Write(orglist)
	}
}

// GetDatasources method gets the list of organizations
func GetDatasources(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		datasourcelist, err := handler.GetDatasources(clients)
		if err != nil {
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return
		}
		w.Write(datasourcelist)
	}
}

// GetDatasourceID gets the datasource details by ID
func GetDatasourceID(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		datasourceID := chi.URLParam(r, "datasourceID")
		ID, err := strconv.Atoi(datasourceID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"provided DatasourceID is not integer")
			return
		}
		datasource, err := handler.GetDatasourceID(clients, ID)
		if err != nil {
			switch err.(type) {
			// common.DatasourceNotFound  means, that user provided the
			// ID of non existent user. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("Datasource Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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

		w.Write(datasource)
	}
}

// DeleteDatasource method deletes the organzation
func DeleteDatasource(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		datasourceID := chi.URLParam(r, "datasourceID")
		ID, err := strconv.Atoi(datasourceID)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusUnprocessableEntity,
				http.StatusText(http.StatusUnprocessableEntity),
				"provided DatasourceID is not integer")
			return
		}

		// check if the ID exists
		_, err = handler.GetDatasourceID(clients, ID)
		if err != nil {
			switch err.(type) {
			// common.DatasourceNotFound  means, that user provided the
			// ID of non existent user. We return 404
			case visualization.NotFound:
				errMsg := fmt.Sprintf("Datasource Not Found")
				common.WriteErrorToResponse(w, http.StatusNotFound,
					errMsg, err.Error())
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

		// delete the datasource if ID exists
		err = handler.DeleteDatasource(clients, ID)
		if err != nil {
			log.Logger.Error(err)
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"Internal server error occured")
			return
		}
	}
}

// CreateDatasource method creates a user
func CreateDatasource(clients *common.ClientContainer, handler common.HandlerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		datasource := DataSource{}
		err := json.NewDecoder(r.Body).Decode(&datasource)
		if err != nil {
			common.WriteErrorToResponse(w, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				err.Error())
			return
		}
		if len(datasource.Name) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Name in parameters")
			return
		}
		if len(datasource.URL) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Url in parameters")
			return
		}
		if len(datasource.Database) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Database in parameters")
			return
		}
		if len(datasource.User) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide User in parameters")
			return
		}
		if len(datasource.Password) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Password in parameters")
			return
		}
		if len(datasource.Type) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Type in parameters")
			return
		}
		if len(datasource.Access) == 0 {
			common.WriteErrorToResponse(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				"provide Access in parameters")
			return
		}

		// Create Datasource if no error
		helperCreateDatasource(clients, handler, w, datasource)
	}
}
