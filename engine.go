package visualization

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/kbhonagiri16/visualization-client/logging"
	//_ "github.com/kbhonagiri16/visualization-client/models"
	"github.com/satori/go.uuid"
)

var (
	engine *xorm.Engine
)

// DatabaseManager represents what functionality we are expecting from db layer
type DatabaseManager interface {
	QueryVisualizationsDashboards(string, string, string, map[string]interface{}) (
		*map[Visualization][]*Dashboard, error)
	CreateVisualizationsWithDashboards(string, string, map[string]interface{},
		[]string, []string) (*Visualization, []*Dashboard, error)
	DeleteVisualization(*Visualization) error
	BulkUpdateDashboard([]*Dashboard) error
	BulkDeleteDashboard([]*Dashboard) error
	GetVisualizationWithDashboardsBySlug(string, string) (*Visualization, []*Dashboard, error)
}

// InitializeEngine initializes connection to db
func InitializeEngine(mysqlUsername, mysqlPassword, mysqlHost,
	mysqlDatabaseName string, mysqlPort int) error {
	var err error
	// connection parameter has to look like
	// user:password@tcp(host:port)/dbname))
	engine, err = xorm.NewEngine("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		mysqlUsername,
		mysqlPassword,
		mysqlHost,
		mysqlPort,
		mysqlDatabaseName,
	))
	return err
}

// GetEngine returns database connection
func GetEngine() *xorm.Engine {
	return engine
}

// XORMManager is an implementation of DatabaseManager
type XORMManager struct {
	engine *xorm.Engine
}

// NewXORMManager is XORMManager constructor
func NewXORMManager() *XORMManager {
	return &XORMManager{engine}
}

// CreateVisualizationFromParam takes provided arguments and returns created model
// without storing to db
func (m *XORMManager) CreateVisualizationFromParam(name, organizationID string,
	tags map[string]interface{}) (*Visualization, error) {

	log.Logger.Debugf("Creating new Visualization entry named '%s'", name)

	// XORM framework does not support generic JSON fields
	// That's why XORM processes data as serialized json string
	// That means that map[string]interface{} tags must be serialized to string
	encodedTags, err := json.Marshal(tags)
	if err != nil {
		log.Logger.Errorf("Error on storing not serializable map[string]interface{}"+
			" to json field : '%s'", err)
		return nil, err
	}

	visualization := &Visualization{
		Slug:           uuid.NewV4().String(),
		Name:           name,
		OrganizationID: organizationID,
		Tags:           string(encodedTags),
	}

	return visualization, nil
}

func getVisualizationLookupQuery(slug, name, organizationID string,
	tags map[string]interface{}) (string, []interface{}) {

	// create Query, with ? placeholders for queries. This would protect from
	// sql injection attacks. Function returns query and parameters to be passed to it
	queryChunks := []string{}
	queryParams := []interface{}{}

	if slug != "" {
		queryChunks = append(queryChunks, fmt.Sprintf("%s.%s = ?",
			VisualizationTableName, VisualizationSlugColumn))
		queryParams = append(queryParams, slug)
	}

	if name != "" {
		queryChunks = append(queryChunks, fmt.Sprintf("%s.%s = ?",
			VisualizationTableName, VisualizationNameColumn))
		queryParams = append(queryParams, name)
	}

	if organizationID != "" {
		queryChunks = append(queryChunks, fmt.Sprintf("%s.%s = ?",
			VisualizationTableName, VisualizationOrgColumn))
		queryParams = append(queryParams, organizationID)
	}

	// use JSON_EXTRACT to query json_field
	for tagName, tagValue := range tags {
		queryChunks = append(queryChunks, fmt.Sprintf(
			"JSON_EXTRACT(%s, '$.%s') = ?", VisualizationTagsColumn, tagName))
		queryParams = append(queryParams, tagValue)
	}

	query := strings.Join(queryChunks, " AND ")

	log.Logger.Debugf("Got lookup query '%s'", query)
	return query, queryParams
}

// QueryVisualizationsDashboards takes name, tags and organizationID and returns matched entries
func (m *XORMManager) QueryVisualizationsDashboards(slug, name, organizationID string,
	tags map[string]interface{}) (*map[Visualization][]*Dashboard, error) {

	query, queryParams := getVisualizationLookupQuery(slug, name, organizationID, tags)

	var queryResult []struct {
		Visualization Visualization `xorm:"extends"`
		Dashboard     Dashboard     `xorm:"extends"`
	}
	err := m.engine.Table(VisualizationTableName).Join(
		"INNER", DashboardTableName,
		fmt.Sprintf("%s.%s = %s.%s", DashboardTableName,
			DashboardVisualizationColumn,
			VisualizationTableName,
			VisualizationIDColumn)).Where(
		query, queryParams...).Find(&queryResult)
	if err != nil {
		log.Logger.Errorf("Error on getting visualizations from db: '%s'", err)
		return nil, err
	}

	result := map[Visualization][]*Dashboard{}
	for index, queryEntry := range queryResult {
		result[queryEntry.Visualization] = append(result[queryEntry.Visualization],
			&queryResult[index].Dashboard)
	}

	return &result, nil
}

// GetVisualizationWithDashboardsBySlug returs visualization with all related dashboards
func (m *XORMManager) GetVisualizationWithDashboardsBySlug(
	slug, organizationID string) (*Visualization, []*Dashboard, error) {
	// TODO(oshyman) fix lookup query
	noNameProvided := ""
	noTagsProvided := map[string]interface{}{}
	query, queryParams := getVisualizationLookupQuery(slug, noNameProvided, organizationID, noTagsProvided)

	var queryResult []struct {
		Visualization Visualization `xorm:"extends"`
		Dashboard     Dashboard     `xorm:"extends"`
	}
	err := m.engine.Table("visualization").Join("INNER", DashboardTableName,
		fmt.Sprintf("%s.%s = %s.%s", DashboardTableName,
			DashboardVisualizationColumn,
			VisualizationTableName,
			VisualizationIDColumn)).Where(
		query, queryParams...).Find(&queryResult)
	if err != nil {
		log.Logger.Errorf("Error on getting visualizations from db: '%s'", err)
		return nil, nil, err
	}

	var visualizationDatabase *Visualization
	dashboardsDatabase := []*Dashboard{}
	for index := range queryResult {
		visualizationDatabase = &queryResult[index].Visualization
		dashboardsDatabase = append(dashboardsDatabase, &queryResult[index].Dashboard)
	}
	return visualizationDatabase, dashboardsDatabase, nil
}

func getBulkDeleteDashboardQuery(dashboards []*Dashboard) (string, []interface{}) {
	if len(dashboards) > 0 {
		// create Query, with ? placeholders for queries. This would protect from
		// sql injection attacks. Function returns query and parameters to be passed to it
		ids := strings.Repeat("?, ", len(dashboards)-1)
		ids = ids + "?"
		queryParams := []interface{}{}
		for _, dashboard := range dashboards {
			queryParams = append(queryParams, dashboard.ID)
		}
		query := fmt.Sprintf("DELETE FROM %s WHERE id IN (%s);",
			DashboardTableName, ids)
		log.Logger.Debugf("Bulk Delete dashboard query is '%s'", query)
		return query, queryParams
	}
	return "", nil
}

// BulkDeleteDashboard removes multiple Dashboards at once
func (m *XORMManager) BulkDeleteDashboard(dashboards []*Dashboard) error {
	// Xorm does not support bulk delete
	if len(dashboards) > 0 {
		query, queryParams := getBulkDeleteDashboardQuery(dashboards)
		_, err := engine.Exec(query, queryParams...)
		return err
	}
	return nil
}

//BulkUpdateDashboard updates multiple records at once
func (m *XORMManager) BulkUpdateDashboard(dashboards []*Dashboard) error {
	// Xorm does not provide bulk update
	if len(dashboards) > 0 {
		table := m.engine.TableInfo(Dashboard{})
		queryParams := []interface{}{}
		columnDBNames := []string{}
		columnDBValues := []string{}
		for _, colName := range table.Columns() {
			columnDBNames = append(columnDBNames, colName.Name)
			columnDBValues = append(columnDBValues,
				fmt.Sprintf("%s=VALUES(%[1]s)", colName.Name))
		}
		columnDBNamesString := strings.Join(columnDBNames, ", ")
		columnUpdateString := strings.Join(columnDBValues, ", ")

		amountOfParameters := reflect.ValueOf(
			&Dashboard{}).Elem().NumField()
		singleDashboardParameter := "(" + strings.Repeat(
			"?, ", amountOfParameters-1) + "?)"
		renderedParameters := strings.Repeat(
			singleDashboardParameter+", ", len(dashboards)-1) + singleDashboardParameter

		for _, dashboard := range dashboards {
			reflection := reflect.ValueOf(dashboard).Elem()
			for i := 0; i < reflection.NumField(); i++ {
				field := reflection.Field(i)
				queryParams = append(queryParams, field.Interface())
			}
		}
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON DUPLICATE KEY UPDATE %s;",
			DashboardTableName, columnDBNamesString,
			renderedParameters, columnUpdateString)
		log.Logger.Debugf("Bulk Update dashboard query is '%s'", query)

		_, err := engine.Exec(query, queryParams...)
		return err
	}
	return nil
}

// DeleteVisualization removes visualization model from db
func (m *XORMManager) DeleteVisualization(visualization *Visualization) error {
	if visualization != nil {
		_, err := m.engine.Id(visualization.ID).Delete(&Visualization{})
		return err
	}
	return nil
}

// CreateVisualizationsWithDashboards creates all data for single visualization
// in one transaction
func (m *XORMManager) CreateVisualizationsWithDashboards(name, organizationID string,
	tags map[string]interface{}, dashboardNames, renderedTemplates []string) (
	*Visualization, []*Dashboard, error) {

	// validate data for visualization
	visualization, err := m.CreateVisualizationFromParam(name, organizationID, tags)
	if err != nil {
		return nil, nil, err
	}

	session := m.engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return nil, nil, err
	}
	_, err = session.Insert(visualization)
	if err != nil {
		session.Rollback()
		return nil, nil, err
	}

	var dashboards []*Dashboard
	for index, name := range dashboardNames {
		dashboards = append(dashboards, &Dashboard{
			ID:               uuid.NewV4().String(),
			Visualization:    visualization.ID,
			Name:             name,
			RenderedTemplate: renderedTemplates[index],
		})
	}

	_, err = session.Insert(dashboards)
	if err != nil {
		session.Rollback()
		return nil, nil, err
	}

	err = session.Commit()
	if err != nil {
		return nil, nil, err
	}
	return visualization, dashboards, nil
}
