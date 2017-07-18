package visualization

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/kbhonagiri16/visualization-client/logging"
	"github.com/mitchellh/mapstructure"
	"time"
)

const tokenAuthHeaderName = "X-Subject-Token"

/*Client :

realization of openstack.ClientInterface.
All openstack client and credentials are encapsulated inside struct and
are not exported. To use any openstack functionality - define method of this
struct, and update openstack.ClientInteface.
This abstraction is implemented, because we don't want to have real calls
to openstack in unit tests. Using function to proxy all openstack calls
gives us oportunity to mock proxy functions during unit testing stage
*/
type Client struct {
	credentialsProvider *gophercloud.ProviderClient
	keystoneClient      *gophercloud.ServiceClient
}

// ValidateToken - ask keystone if token is still valid
func (cli *Client) ValidateToken(token string) (bool, error) {
	valid, err := tokens.Validate(cli.keystoneClient, token)
	if err != nil {
		log.Logger.Errorf("Error validating token %s", err)
	}
	return valid, err
}

// GetTokenInfo - get data from token
func (cli *Client) GetTokenInfo(token string) (*TokenInfo, error) {
	response := tokens.Get(cli.keystoneClient, token)

	if response.Err != nil {
		log.Logger.Errorf("Error retrieving token %s", response.Err)
		return nil, response.Err
	}

	// this body parsing code is implemented, because gophercloud parses only
	// ExpiresAt value, everything else is not parsed
	var parsedResponse struct {
		Token struct {
			ExpiresAt string              `mapstructure:"expires_at"`
			Roles     []map[string]string `mapstructure:"roles"`
			Project   struct {
				ID   string `mapstructure:"id"`
				Name string `mapstructure:"name"`
			} `mapstructure:"project"`
		} `mapstructure:"token"`
	}

	err := mapstructure.Decode(response.Body, &parsedResponse)
	if err != nil {
		log.Logger.Errorf("Error parsing token %s", err)
		return nil, err
	}

	var resultToken TokenInfo
	resultToken.ID = response.Header.Get(tokenAuthHeaderName)

	resultToken.ExpiresAt, err = time.Parse(gophercloud.RFC3339Milli,
		parsedResponse.Token.ExpiresAt)
	if err != nil {
		log.Logger.Errorf("Error parsing time in token %s", err)
		return nil, err
	}
	resultToken.Roles = parsedResponse.Token.Roles
	resultToken.ProjectID = parsedResponse.Token.Project.ID
	resultToken.ProjectName = parsedResponse.Token.Project.Name

	log.Logger.Debugf("Successfully retrieved openstack data for token %s",
		token)
	return &resultToken, nil
}

// NewOpenstackClient creates new Client structure from provided arguments
func NewOpenstackClient(identityEndpoint, username, password,
	tenantName, domainName string) (*Client, error) {

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: identityEndpoint,
		Username:         username,
		Password:         password,
		TenantName:       tenantName,
		DomainName:       domainName,
		AllowReauth:      true,
	}

	credentialsProvider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

	keystoneClient, err := openstack.NewIdentityV3(
		credentialsProvider, gophercloud.EndpointOpts{},
	)

	if err != nil {
		return nil, err
	}

	return &Client{credentialsProvider, keystoneClient}, nil
}

/*ClientInterface :

Interface for openstack functionality usage
*/
type ClientInterface interface {
	ValidateToken(string) (bool, error)
	GetTokenInfo(string) (*TokenInfo, error)
}
