package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// VisualizationError errors for Visualization client
type VisualizationError struct {
	code        string
	message     string
	description string
}

// Error generate a error message.
// If Code is zero, we know it's not a http error.
func (e VisualizationError) Error() string {
	return fmt.Sprintf("ERROR: %s", e.description)
}

// VisualizationClient client for Visualization
type VisualizationClient struct {
	url            string
	client         *http.Client
	token          AuthToken
	JWT            string
	openstackToken string
}

// NewVisualizationClient returns client with token
func NewVisualizationClient(url string, client http.Client, openstackToken string) (*VisualizationClient, error) {
	return &VisualizationClient{client: &client, url: url, openstackToken: openstackToken}, nil

}

// reIssue this method reissues the token
func (v *VisualizationClient) reIssue() error {
	token, err := v.Authenticate()
	if err != nil {
		return err
	}
	v.token = token
	v.JWT = token.JWT
	return err
}

// authorizeToken for token checking and reauth
func (v *VisualizationClient) authorizeToken(withAuth bool) {
	if !withAuth {
		// validate token
		tokenExpires := v.token.Token.ExpiresAt.UnixNano() / 1000000
		now := time.Now().UnixNano() / 1000000
		if tokenExpires < now {
			v.reIssue()
		}
	}
	return
}

// doRequest does the authorized request
func (v *VisualizationClient) doRequest(withAuth bool) {
	v.authorizeToken(withAuth)
	return
}

// headerRequest adds header to Request
func (v *VisualizationClient) headerRequest(request *http.Request, withAuth bool) *http.Request {
	if withAuth {
		request.Header.Add("X-OpenStack-Auth-Token", v.openstackToken)
	} else {
		if v.token == (AuthToken{}) {
			v.reIssue()
		}
		bearer := fmt.Sprintf("Bearer %v", v.JWT)
		request.Header.Add("Authorization", bearer)
	}
	return request
}

// httpRequest handles the request to server.
//It returns the response body and a error if something went wrong
func (v *VisualizationClient) httpRequest(method string, url string, body io.Reader, withAuth bool) (result io.Reader, err error) {
	v.doRequest(withAuth)
	request, err := http.NewRequest(method, url, body)
	request.Header.Set("Content-Type", "application/json")
	request = v.headerRequest(request, withAuth)

	response, err := v.client.Do(request)
	var Message VisualizationError
	if err != nil {
		Message.description = err.Error()
		return result, Message
	}

	if response.StatusCode != 200 {
		dec := json.NewDecoder(response.Body)
		err = dec.Decode(&Message)
		if err != nil {
			return
		}
		switch response.StatusCode {
		case 409:
			Message.code = "409"
			Message.message = "Already Exists"
			Message.description = "Provided Details to create exists"
			return result, Message
		case 404:
			Message.code = "404"
			Message.message = "ID not found"
			Message.description = "Provided ID to Delete/Get was not found"
			return result, Message
		}

		return result, Message
	}
	result = response.Body
	return
}

// AuthToken for requests
type AuthToken struct {
	JWT   string `json:"jwt"`
	Token Token  `json:"token"`
}

// Token for AuthToken
type Token struct {
	OrganizationID string    `json:"organizationId"`
	ExpiresAt      time.Time `json:"expiresAt"`
	IsAdmin        bool      `json:"isAdmin"`
}

// UserInOrganization Get Users in organization
type UserInOrganization struct {
	OrgID    string `json:"orgID"`
	UserID   string `json:"userID"`
	Login    string `json:"login"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Org Get organization list
type Org struct {
	OrganizationID string `json:"organizationID"`
	Name           string `json:"name"`
}

// User gets Users List
type User struct {
	UserID   string `json:"userID"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Authenticate gets a openstack token
func (v *VisualizationClient) Authenticate() (token AuthToken, err error) {
	reqURL := v.url + "/auth/openstack"
	response, err := v.httpRequest("POST", reqURL, nil, true)
	if err != nil {
		return
	}

	authToken := json.NewDecoder(response)
	err = authToken.Decode(&token)
	if err != nil {
		return AuthToken{}, err
	}

	return
}

// GetUsers returns list of users
func (v *VisualizationClient) GetUsers() (user []User, err error) {
	reqURL := v.url + "/admin/users"
	response, err := v.httpRequest("GET", reqURL, nil, false)
	if err != nil {
		return []User{}, err
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&user)

	return
}

// GetUserName returns user by Name
func (v *VisualizationClient) GetUserName(name string) (user User, err error) {
	users, err := v.GetUsers()
	if err != nil {
		return
	}

	for _, elem := range users {
		if elem.Name == name {
			user = elem
		}
	}
	return
}

// GetUserID Get User by ID
func (v *VisualizationClient) GetUserID(ID string) (user User, err error) {
	reqURL := fmt.Sprintf("%s/admin/users/%s", v.url, ID)
	response, err := v.httpRequest("GET", reqURL, nil, false)
	if err != nil {
		return
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&user)

	return
}

// CreateUser creates a user
func (v *VisualizationClient) CreateUser(user User) (userDetails User, err error) {
	reqURL := v.url + "/admin/users"
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return
	}

	response, err := v.httpRequest("POST", reqURL, bytes.NewBuffer(jsonStr), false)
	if err != nil {
		return
	}
	dec := json.NewDecoder(response)
	err = dec.Decode(&userDetails)

	return
}

// DeleteUser Delete the user with given id
func (v *VisualizationClient) DeleteUser(ID string) (user User, err error) {
	reqURL := fmt.Sprintf("%s/users/%s", v.url, ID)

	response, err := v.httpRequest("DELETE", reqURL, nil, false)
	if err != nil {
		return
	}
	dec := json.NewDecoder(response)
	err = dec.Decode(&user)

	return
}

// GetOrganizations returns list of organizations
func (v *VisualizationClient) GetOrganizations() (org []Org, err error) {
	reqURL := v.url + "/admin/organizations"
	response, err := v.httpRequest("GET", reqURL, nil, false)
	if err != nil {
		return []Org{}, err
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&org)

	return
}

// GetOrganizationID Get Org by ID
func (v *VisualizationClient) GetOrganizationID(OrgID string) (orgID Org, err error) {
	reqURL := fmt.Sprintf("%s/admin/organizations/%s", v.url, OrgID)
	response, err := v.httpRequest("GET", reqURL, nil, false)
	if err != nil {
		return
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&orgID)

	return
}

// DeleteOrganization Delete the organization with given id
func (v *VisualizationClient) DeleteOrganization(ID string) (org Org, err error) {
	reqURL := fmt.Sprintf("%s/admin/organizations/%s", v.url, ID)

	response, err := v.httpRequest("DELETE", reqURL, nil, false)
	if err != nil {
		return
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&org)

	return
}

// CreateOrganization creates a organization
func (v *VisualizationClient) CreateOrganization(org Org) (orgs Org, err error) {
	reqURL := v.url + "/admin/organizations"
	jsonStr, err := json.Marshal(org)
	if err != nil {
		return
	}

	response, err := v.httpRequest("POST", reqURL, bytes.NewBuffer(jsonStr), false)
	if err != nil {
		return
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&orgs)

	return
}

// GetOrganizationUsers gets Users in Organisation
func (v *VisualizationClient) GetOrganizationUsers(ID string) (org []UserInOrganization, err error) {
	reqURL := fmt.Sprintf("%s/admin/organizations/%s/users", v.url, ID)
	response, err := v.httpRequest("GET", reqURL, nil, false)
	if err != nil {
		return []UserInOrganization{}, err
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&org)
	return
}

// DeleteOrganizationUser Delete User in Organisation
func (v *VisualizationClient) DeleteOrganizationUser(userID string, orgID string) (org UserInOrganization, err error) {
	reqURL := fmt.Sprintf("%s/admin/organizations/%s/users/%s", v.url, orgID, userID)

	response, err := v.httpRequest("DELETE", reqURL, nil, false)
	if err != nil {
		return
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&org)

	return
}

// CreateUserOrganization Add User in Organisation
func (v *VisualizationClient) CreateUserOrganization(OrgID string, user UserInOrganization) (org UserInOrganization, err error) {
	reqURL := fmt.Sprintf("%s/admin/organizations/%s/users", v.url, OrgID)
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return UserInOrganization{}, err
	}

	response, err := v.httpRequest("POST", reqURL, bytes.NewBuffer(jsonStr), false)
	if err != nil {
		return UserInOrganization{}, err
	}

	dec := json.NewDecoder(response)
	err = dec.Decode(&org)

	return
}
