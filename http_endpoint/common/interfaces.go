package common

import (
	"github.com/kbhonagiri16/visualization-client"
	"time"
)

/*ClientContainer represents container for storing different clients
It was created to have mockable architecture*/
type ClientContainer struct {
	Openstack       visualization.ClientInterface
	Grafana         visualization.SessionInterface
	DatabaseManager visualization.DatabaseManager
}

/*HandlerInterface represents set of handlers for api
It was created to have mockable architecture*/
type HandlerInterface interface {
	AuthOpenstack(*ClientContainer, ClockInterface, string, string) ([]byte, error)
	GetDatasources(*ClientContainer) ([]byte, error)
	GetDatasourceID(*ClientContainer, int) ([]byte, error)
	DeleteDatasource(*ClientContainer, int) error
	CreateDatasource(*ClientContainer, []byte) error
	GetUsers(*ClientContainer) ([]byte, error)
	GetUserID(*ClientContainer, int) ([]byte, error)
	DeleteUser(*ClientContainer, int) error
	CreateUser(*ClientContainer, []byte) error
	GetOrganizations(*ClientContainer) ([]byte, error)
	GetOrganizationID(*ClientContainer, int) ([]byte, error)
	DeleteOrganization(*ClientContainer, int) error
	CreateOrganization(*ClientContainer, []byte) error
	CreateOrganizationUser(*ClientContainer, int, []byte) error
	DeleteOrganizationUser(*ClientContainer, int, int) error
	GetOrganizationUsers(*ClientContainer, int) ([]byte, error)
	VisualizationsGet(*ClientContainer, string, string,
		map[string]interface{}) (*[]VisualizationWithDashboards, error)
	VisualizationsPost(*ClientContainer, VisualizationPOSTData, string) (
		*VisualizationWithDashboards, error)
	VisualizationDelete(*ClientContainer, string, string) (*VisualizationWithDashboards, error)
}

// ClockInterface serves for testing purposes of functions, that require time
type ClockInterface interface {
	Now() time.Time
}
