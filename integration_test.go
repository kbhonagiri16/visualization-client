// +build integration

package grafanaclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ds = DataSource{Name: "testme",
	Type:      "influxdb",
	Access:    "proxy",
	URL:       "http://localhost:8086",
	User:      "root",
	Password:  "root",
	Database:  "test",
	IsDefault: true}

var usr = AdminCreateUser{Email: "test@me.com",
	Login:    "testme",
	Name:     "testme",
	Password: "test"}

var org = Org{Name: "testme"}

var org_list = OrgList{ID: 1, Name: "Main Org."}

var org_user = CreateOrganizationUser{Email: "test@me.com",
	Login:    "test",
	Name:     "test",
	Password: "test"}

var url = getenv("GRAFANA_URL")
var user = getenv("GRAFANA_USER")
var pass = getenv("GRAFANA_PASS")

func getenv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		fmt.Printf("Environment variable %s in not defined \n", key)
		os.Exit(1)
	}
	return value
}

func Test_DoLogon(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
}

func Test_CreateDataSource(t *testing.T) {
	t.Skip("TODO(illia) fix it later")
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	err = session.CreateDataSource(ds)
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when creating DataSource: %s", err))
}

func Test_GetDataSourceList(t *testing.T) {
	t.Skip("TODO(illia) fix it later")
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	dslist, err := session.GetDataSourceList()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting DataSource: %s", err))
	var check bool
	for _, ds := range dslist {
		if ds.Name == "testme" {
			check = true
		}
	}

	assert.Equal(t, true, check, "We didn't find the test datasource")
}

func Test_GetDataSourceListID(t *testing.T) {
	t.Skip("TODO(illia) fix it later")
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	dslist, err := session.GetDataSourceList()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting DataSource: %s", err))
	for _, ds := range dslist {
		if ds.Name == "testme" {
			resDs, _ := session.GetDataSourceListID(ds.ID)

			assert.Equal(t, "testme", resDs.Name, "We are expecting to retrieve DataSource with ID 1 and didn't get it")
		}
	}
}

func Test_GetDataSourceName(t *testing.T) {
	t.Skip("TODO(illia) fix it later")
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))

	resDs, _ := session.GetDataSourceName("testme")

	assert.Equal(t, "testme", resDs.Name, "We are expecting to retrieve testme DataSource and didn't get it")
}

func Test_DeleteDataSource(t *testing.T) {
	t.Skip("TODO(illia) fix it later")
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))

	resDs, err := session.GetDataSourceName("testme")
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Getting Datasource details: %s", err))

	err = session.DeleteDataSource(resDs.ID)
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Deleting: %s", err))
}

func Test_CreateUser(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	err = session.CreateUser(usr)
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when creating user: %s", err))
}

func Test_GetUsers(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	usrlist, err := session.GetUsers()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting DataSource: %s", err))
	var check bool
	for _, usr := range usrlist {
		if usr.Name == "testme" {
			check = true
		}
	}

	assert.Equal(t, true, check, "We didn't find the testme user")
}

func Test_GetUserID(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	usrlist, err := session.GetUsers()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting DataSource: %s", err))
	for _, usr := range usrlist {
		if usr.Name == "testme" {
			resDs, _ := session.GetUserID(usr.ID)

			assert.Equal(t, "testme", resDs.Name, "We are expecting to retrieve User with ID 1 and didn't get it")
		}
	}
}

func Test_DeleteUser(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))

	resUsers, err := session.GetUsers()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one while getting Users: %s", err))

	for _, users := range resUsers {
		if users.Name == "testme" {
			err = session.DeleteUser(users.ID)
			assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Deleting User :%s", err))
		}
	}
}

func Test_CreateOrg(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	_, err = session.CreateOrg(org)
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when creating DataSource: %s", err))
}

func Test_GetOrCreateOrgByName(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	_, err = session.GetOrCreateOrgByName("test_name")
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error: %s", err))
}

func Test_GetOrgs(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))
	orglist, err := session.GetOrganizations()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting Organization: %s", err))
	var check bool
	for _, orgs := range orglist {
		if orgs.Name == "testme" {
			check = true
		}
	}

	assert.Equal(t, true, check, "We didn't find the testme organization")
}

func Test_GetOrgID(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))

	orglist, err := session.GetOrganizations()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting Organization: %s", err))
	for _, orgs := range orglist {
		if orgs.Name == "testme" {
			resDs, _ := session.GetOrganizationID(orgs.ID)

			assert.Equal(t, "testme", resDs.Name, "We are expecting to retrieve Org with ID 1 and didn't get it")
		}
	}
}

func Test_CreateOrgUser(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))

	orglist, err := session.GetOrganizations()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting Organization: %s", err))
	for _, orgs := range orglist {
		if orgs.Name == "testme" {
			session.CreateOrganizationUser(orgs.ID, org_user)
		}
	}
	userlist, err := session.GetUsers()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting Organization: %s", err))
	var check bool
	for _, users := range userlist {
		if users.Name == "test" {
			check = true
		}
	}
	assert.Equal(t, true, check, "We didn't find the test user")
}

func Test_GetOrgUsers(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))

	orglist, err := session.GetOrganizationUsers(org_list.ID)
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one getting users in Organization: %s", err))
	var check bool
	for _, orgs := range orglist {
		if orgs.Login == "test" {
			check = true
		}
	}

	assert.Equal(t, true, check, "We didn't find the admin user in Main organization")
}

func Test_DeleteOrgUser(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))

	resUsers, err := session.GetUsers()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Getting Users: %s", err))

	resOrgs, err := session.GetOrganizations()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Getting Orgs: %s", err))

	for _, users := range resUsers {
		if users.Name == "test" {
			for _, orgs := range resOrgs {
				if orgs.Name == "testme" {
					err = session.DeleteOrganizationUser(users.ID, orgs.ID)
					assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Deleting User: %s", err))
				}
			}
		}
	}
}

func Test_DeleteOrg(t *testing.T) {
	session, _ := NewSession(user, pass, url)
	err := session.DoLogon()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Login: %s", err))

	resOrgs, err := session.GetOrganizations()
	assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Getting Orgs: %s", err))

	for _, orgs := range resOrgs {
		if orgs.Name == "testme" {
			err = session.DeleteOrganization(orgs.ID)
			assert.Nil(t, err, fmt.Sprintf("We are expecting no error and got one when Deleting Org: %s", err))
		}
	}
}
