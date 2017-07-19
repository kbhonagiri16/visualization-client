package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserID(t *testing.T) {
	tests := []struct {
		description      string
		users            string
		expectedData     User
		userID           string
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			expectedData:   User{UserID: "1", Email: "test@test.com", Name: "test", Login: "test", Password: ""},
			users:          "{\"UserID\":\"1\",\"Email\":\"test@test.com\",\"Name\":\"test\",\"Login\":\"test\",\"Password\":\"\"}",
			userID:         "1",
			token:          "token",
			testStatusCode: false,
		},
		{
			description:      "provided ID not found",
			expectedData:     User{UserID: "", Email: "", Name: "", Login: "", Password: ""},
			users:            "{\"UserID\":\"\",\"Email\":\"\",\"Name\":\"\",\"Login\":\"\",\"Password\":\"\"}",
			userID:           "1",
			token:            "token",
			expectedResponse: VisualizationError{code: "404", message: "ID not found", description: "Provided ID to Delete/Get was not found"},
			testStatusCode:   true,
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusNotFound)
			}
			fmt.Fprint(w, testCase.users)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.GetUserID(testCase.userID)
		assert.Equal(t, resp, testCase.expectedData, "response match")
		if testCase.testStatusCode {
			assert.Equal(t, err, testCase.expectedResponse, "response code match")
		} else {
			assert.Equal(t, err, nil, "no error")
		}
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		description  string
		users        string
		expectedData []User
		token        string
	}{
		{
			description:  "make sure handler reacts",
			expectedData: []User{User{UserID: "1", Email: "test@test.com", Name: "test", Login: "test", Password: ""}},
			users:        "[{\"UserID\":\"1\",\"Email\":\"test@test.com\",\"Name\":\"test\",\"Login\":\"test\",\"Password\":\"\"}]",
			token:        "token",
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testCase.users)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.GetUsers()
		assert.Equal(t, err, nil, "error")
		assert.Equal(t, resp, testCase.expectedData, "response match")
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		description      string
		users            string
		expectedData     User
		userID           string
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			expectedData:   User{UserID: "1", Email: "test@test.com", Name: "test", Login: "test", Password: ""},
			users:          "{\"UserID\":\"1\",\"Email\":\"test@test.com\",\"Name\":\"test\",\"Login\":\"test\",\"Password\":\"\"}",
			userID:         "1",
			token:          "token",
			testStatusCode: false,
		},
		{
			description:      "provided ID not found",
			expectedData:     User{UserID: "", Email: "", Name: "", Login: "", Password: ""},
			users:            "{\"UserID\":\"\",\"Email\":\"\",\"Name\":\"\",\"Login\":\"\",\"Password\":\"\"}",
			userID:           "1",
			token:            "token",
			expectedResponse: VisualizationError{code: "404", message: "ID not found", description: "Provided ID to Delete/Get was not found"},
			testStatusCode:   true,
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusNotFound)
			}
			fmt.Fprint(w, testCase.users)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.DeleteUser(testCase.userID)
		assert.Equal(t, resp, testCase.expectedData, "response match")
		if testCase.testStatusCode {
			assert.Equal(t, err, testCase.expectedResponse, "response code match")
		} else {
			assert.Equal(t, err, nil, "no error")
		}
	}
}

func TestGetOrganizationID(t *testing.T) {
	tests := []struct {
		description      string
		orgs             string
		expectedData     Org
		expectedResponse VisualizationError
		testStatusCode   bool
		orgID            string
		token            string
	}{
		{
			description:    "make sure handler reacts",
			expectedData:   Org{OrganizationID: "", Name: "test"},
			orgs:           "{\"OrganizationID\":\"\",\"Name\":\"test\"}",
			orgID:          "1",
			token:          "token",
			testStatusCode: false,
		},
		{
			description:      "ID not found",
			expectedData:     Org{OrganizationID: "", Name: ""},
			orgs:             "{\"OrganizationID\":\"\",\"Name\":\"\"}",
			orgID:            "1",
			token:            "token",
			testStatusCode:   true,
			expectedResponse: VisualizationError{code: "404", message: "ID not found", description: "Provided ID to Delete/Get was not found"},
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusNotFound)
			}
			fmt.Fprint(w, testCase.orgs)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.GetOrganizationID(testCase.orgID)
		assert.Equal(t, resp, testCase.expectedData, "response match")
		if testCase.testStatusCode {
			assert.Equal(t, err, testCase.expectedResponse, "response code match")
		} else {
			assert.Equal(t, err, nil, "no error")
		}
	}
}

func TestGetOrganizations(t *testing.T) {
	tests := []struct {
		description      string
		orgs             string
		expectedData     []Org
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			expectedData:   []Org{Org{OrganizationID: "", Name: "test"}},
			orgs:           "[{\"OrganizationID\":\"\",\"Name\":\"test\"}]",
			token:          "token",
			testStatusCode: false,
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testCase.orgs)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.GetOrganizations()
		assert.Equal(t, err, nil, "error")
		assert.Equal(t, resp, testCase.expectedData, "response match")
	}
}

func TestDeleteOrganization(t *testing.T) {
	tests := []struct {
		description      string
		orgs             string
		expectedData     Org
		orgID            string
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			expectedData:   Org{OrganizationID: "", Name: "test"},
			orgs:           "{\"OrganizationID\":\"\",\"Name\":\"test\"}",
			orgID:          "1",
			token:          "token",
			testStatusCode: false,
		},
		{
			description:      "provided ID not found",
			expectedData:     Org{OrganizationID: "", Name: ""},
			orgs:             "{\"OrganizationID\":\"\",\"Name\":\"test\"}",
			orgID:            "1",
			token:            "token",
			expectedResponse: VisualizationError{code: "404", message: "ID not found", description: "Provided ID to Delete/Get was not found"},
			testStatusCode:   true,
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusNotFound)
			}
			fmt.Fprint(w, testCase.orgs)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.DeleteOrganization(testCase.orgID)
		assert.Equal(t, resp, testCase.expectedData, "response match")
		if testCase.testStatusCode {
			assert.Equal(t, err, testCase.expectedResponse, "response code match")
		} else {
			assert.Equal(t, err, nil, "no error")
		}
	}
}

func TestGetOrganizationUsers(t *testing.T) {
	tests := []struct {
		description      string
		users            string
		expectedCode     int
		expectedData     []UserInOrganization
		orgID            string
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			expectedData:   []UserInOrganization{UserInOrganization{OrgID: "1", UserID: "1", Role: "Viewer", Email: "test@test.com", Login: "test", Password: ""}},
			users:          "[{\"OrgID\":\"1\",\"UserID\":\"1\",\"Role\":\"Viewer\",\"Email\":\"test@test.com\",\"Login\":\"test\",\"Password\":\"\"}]",
			orgID:          "1",
			token:          "token",
			testStatusCode: false,
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusNotFound)
			}
			fmt.Fprint(w, testCase.users)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.GetOrganizationUsers(testCase.orgID)
		assert.Equal(t, err, nil, "error")
		assert.Equal(t, resp, testCase.expectedData, "response match")
	}
}

func TestDeleteOrganizationUser(t *testing.T) {
	tests := []struct {
		description      string
		users            string
		expectedCode     int
		expectedData     UserInOrganization
		orgID            string
		userID           string
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			expectedData:   UserInOrganization{OrgID: "1", UserID: "1", Role: "Viewer", Email: "test@test.com", Login: "test", Password: ""},
			users:          "{\"OrgID\":\"1\",\"UserID\":\"1\",\"Role\":\"Viewer\",\"Email\":\"test@test.com\",\"Login\":\"test\",\"Password\":\"\"}",
			orgID:          "1",
			userID:         "1",
			token:          "token",
			testStatusCode: false,
		},
		{
			description:      "ID not found",
			expectedData:     UserInOrganization{OrgID: "", UserID: "", Role: "", Email: "", Login: "", Password: ""},
			users:            "{\"OrgID\":\"\",\"UserID\":\"\",\"Role\":\"\",\"Email\":\"\",\"Login\":\"\",\"Password\":\"\"}",
			orgID:            "1",
			userID:           "1",
			token:            "token",
			testStatusCode:   true,
			expectedResponse: VisualizationError{code: "404", message: "ID not found", description: "Provided ID to Delete/Get was not found"},
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusNotFound)
			}
			fmt.Fprint(w, testCase.users)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.DeleteOrganizationUser(testCase.orgID, testCase.userID)
		assert.Equal(t, resp, testCase.expectedData, "response match")
		if testCase.testStatusCode {
			assert.Equal(t, err, testCase.expectedResponse, "response code match")
		} else {
			assert.Equal(t, err, nil, "no error")
		}
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		description      string
		users            string
		expectedData     User
		input            User
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			input:          User{Email: "test@test.com", Name: "test", Login: "test", Password: "pass"},
			expectedData:   User{UserID: "1", Email: "test@test.com", Name: "test", Login: "test", Password: ""},
			users:          "{\"UserID\":\"1\",\"Email\":\"test@test.com\",\"Name\":\"test\",\"Login\":\"test\",\"Password\":\"\"}",
			token:          "token",
			testStatusCode: false,
		},
		{
			description:      "User exist",
			input:            User{Email: "test@test.com", Name: "test", Login: "test", Password: "pass"},
			expectedData:     User{UserID: "", Email: "", Name: "", Login: "", Password: ""},
			users:            "{\"UserID\":\"\",\"Email\":\"\",\"Name\":\"\",\"Login\":\"\",\"Password\":\"\"}",
			token:            "token",
			testStatusCode:   true,
			expectedResponse: VisualizationError{code: "409", message: "Already Exists", description: "Provided Details to create exists"},
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusConflict)
			}
			fmt.Fprint(w, testCase.users)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.CreateUser(testCase.input)
		assert.Equal(t, resp, testCase.expectedData, "response match")
		if testCase.testStatusCode {
			assert.Equal(t, err, testCase.expectedResponse, "response code match")
		} else {
			assert.Equal(t, err, nil, "no error")
		}
	}
}

func TestCreateOrganization(t *testing.T) {
	tests := []struct {
		description      string
		users            string
		expectedData     Org
		input            Org
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			input:          Org{Name: "test"},
			expectedData:   Org{OrganizationID: "1", Name: "test"},
			users:          "{\"OrganizationID\":\"1\",\"Name\":\"test\"}",
			token:          "token",
			testStatusCode: false,
		},
		{
			description:      "Org already Exist",
			input:            Org{Name: "test"},
			expectedData:     Org{OrganizationID: "", Name: ""},
			users:            "{\"OrganizationID\":\"\",\"Name\":\"\"}",
			token:            "token",
			testStatusCode:   true,
			expectedResponse: VisualizationError{code: "409", message: "Already Exists", description: "Provided Details to create exists"},
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusConflict)
			}
			fmt.Fprint(w, testCase.users)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.CreateOrganization(testCase.input)
		assert.Equal(t, resp, testCase.expectedData, "response match")
		if testCase.testStatusCode {
			assert.Equal(t, err, testCase.expectedResponse, "response code match")
		} else {
			assert.Equal(t, err, nil, "no error")
		}
	}
}

func TestCreateUserOrganization(t *testing.T) {
	tests := []struct {
		description      string
		users            string
		expectedData     UserInOrganization
		input            UserInOrganization
		OrgID            string
		token            string
		expectedResponse VisualizationError
		testStatusCode   bool
	}{
		{
			description:    "make sure handler reacts",
			input:          UserInOrganization{OrgID: "1", Email: "test@test.com", Login: "test", Password: "pass", Role: "Viewer"},
			expectedData:   UserInOrganization{OrgID: "1", Email: "test@test.com", Login: "test", Password: "pass", Role: "Viewer"},
			users:          "{\"OrgID\":\"1\",\"Email\":\"test@test.com\",\"Login\":\"test\",\"Password\":\"pass\", \"Role\":\"Viewer\"}",
			OrgID:          "1",
			token:          "token",
			testStatusCode: false,
		},
		{
			description:      "User already exists",
			input:            UserInOrganization{OrgID: "1", Email: "test@test.com", Login: "test", Password: "pass", Role: "Viewer"},
			expectedData:     UserInOrganization{OrgID: "", Email: "", Login: "", Password: "", Role: ""},
			users:            "{\"OrgID\":\"\",\"Email\":\"\",\"Login\":\"\",\"Password\":\"\", \"Role\":\"\"}",
			OrgID:            "1",
			token:            "token",
			testStatusCode:   true,
			expectedResponse: VisualizationError{code: "409", message: "Already Exists", description: "Provided Details to create exists"},
		},
	}
	for _, testCase := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if testCase.testStatusCode {
				w.WriteHeader(http.StatusConflict)
			}
			fmt.Fprint(w, testCase.users)
		}))
		defer ts.Close()
		clientHTTP := http.Client{}
		client, err := NewVisualizationClient(ts.URL, clientHTTP, testCase.token)
		assert.Equal(t, err, nil, "no error")
		resp, err := client.CreateUserOrganization(testCase.OrgID, testCase.input)
		assert.Equal(t, resp, testCase.expectedData, "response match")
		if testCase.testStatusCode {
			assert.Equal(t, err, testCase.expectedResponse, "response code match")
		} else {
			assert.Equal(t, err, nil, "no error")
		}
	}
}
