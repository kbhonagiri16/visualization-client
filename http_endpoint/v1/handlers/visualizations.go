package v1handlers

import (
	"bytes"
	"fmt"
	"github.com/kbhonagiri16/visualization"
	"github.com/kbhonagiri16/visualization/http_endpoint/common"
	"github.com/kbhonagiri16/visualization/logging"
	"github.com/kbhonagiri16/visualization/models"
	"github.com/ulule/deepcopier"
	"text/template"
)

// V1Visualizations implements part of handler interface
type V1Visualizations struct{}

// VisualizationDashboardToResponse transforms models to response format
func VisualizationDashboardToResponse(visualization *models.Visualization,
	dashboards []*models.Dashboard) *common.VisualizationWithDashboards {
	// This function is used, when we have to return visualization with
	// limited number of dashboards (for example in post method)
	log.Logger.Debug("rendering data to user")
	visualizationResponse := &common.VisualizationResponseEntry{}
	dashboardResponse := []*common.DashboardResponseEntry{}
	deepcopier.Copy(visualization).To(visualizationResponse)
	for index := range dashboards {
		dashboardRes := &common.DashboardResponseEntry{}
		deepcopier.Copy(dashboards[index]).To(dashboardRes)
		dashboardResponse = append(dashboardResponse, dashboardRes)
	}
	return &common.VisualizationWithDashboards{
		visualizationResponse, dashboardResponse}
}

// GroupedVisualizationDashboardToResponse transforms map of visualizations to response format
func GroupedVisualizationDashboardToResponse(
	data *map[models.Visualization][]*models.Dashboard) *[]common.VisualizationWithDashboards {
	// This function is used, when

	log.Logger.Debug("rendering data to user")
	response := []common.VisualizationWithDashboards{}
	for visualizationPtr, dashboards := range *data {
		renderedVisualization := VisualizationDashboardToResponse(
			&visualizationPtr, dashboards)
		response = append(response, *renderedVisualization)
	}
	return &response
}

// VisualizationsGet handler queries visualizations
func (h *V1Visualizations) VisualizationsGet(clients *common.ClientContainer,
	organizationID, name string, tags map[string]interface{}) (
	*[]common.VisualizationWithDashboards, error) {
	log.Logger.Debug("Querying data to user according to name and tags")

	data, err := clients.DatabaseManager.QueryVisualizationsDashboards(
		"", name, organizationID, tags)
	if err != nil {
		log.Logger.Errorf("Error getting data from db: '%s'", err)
		return nil, err
	}

	return GroupedVisualizationDashboardToResponse(data), nil
}

func renderTemplates(templates []string, templateParamaters []interface{}) (
	[]string, error) {
	// this function takes Visualization data and returns rendered templates
	log.Logger.Debug("Rendering golang templates")
	renderedTemplates := []string{}
	for index := range templates {
		// validate that golang template is valid
		// "missingkey=error" would return error, if user did not provide
		// all parameters for his own template
		tmpl, err := template.New("").Option(
			"missingkey=error").Parse(templates[index])
		if err != nil {
			// something is wrong with structure of user provided template
			return nil, common.NewUserDataError(
				fmt.Sprintf("ErrorMsg: '%s', TemplateIndex: '%d'",
					err.Error(), index))
		}

		// render golang template with user provided arguments to buffer
		templateBuffer := new(bytes.Buffer)
		err = tmpl.Execute(templateBuffer, templateParamaters[index])
		if err != nil {
			// something is wrong with rendering of user provided template
			return nil, common.NewUserDataError(err.Error())
		}
		renderedTemplates = append(renderedTemplates, templateBuffer.String())
	}
	return renderedTemplates, nil
}

// VisualizationsPost handler creates new visualizations
func (h *V1Visualizations) VisualizationsPost(clients *common.ClientContainer,
	data common.VisualizationPOSTData, organizationID string) (
	*common.VisualizationWithDashboards, error) {

	/*
		1 - validate and render  all golang templates provided by user,
		    if there are any errors, then immediately return error to user
		2 - validate that rendered templates matches grafana json structure
			if there are any mismatch - return error to user
		3 - create db entry for visualization and every dashboard.
		4 - for each validated template - upload it to grafana, store received
			slug for future update of dashboard db entry
		5 - return data to user
	*/

	log.Logger.Debug("Extracting names, templates, data from provided user data")
	templates := []string{}
	templateParamaters := []interface{}{}
	dashboardNames := []string{}
	for _, dashboardData := range data.Dashboards {
		templates = append(templates, dashboardData.TemplateBody)
		templateParamaters = append(templateParamaters, dashboardData.TemplateParameters)
		dashboardNames = append(dashboardNames, dashboardData.Name)
	}
	log.Logger.Debug("Extracted names, templates, data from provided user data")

	renderedTemplates, err := renderTemplates(templates, templateParamaters)
	if err != nil {
		return nil, err
	}

	// create db entries for visualizations and dashboards
	log.Logger.Debug("Creating database entries for visualizations and dashboards")
	visualizationDB, dashboardsDB, err := clients.DatabaseManager.CreateVisualizationsWithDashboards(
		data.Name, organizationID, data.Tags, dashboardNames, renderedTemplates)
	log.Logger.Debug("Created database entries for visualizations and dashboards")
	if err != nil {
		return nil, err
	}

	/*
		Here concistency problem is faced. We can not guarantee, that data,
		stored in database would successfully be updated in grafana, due to
		possible errors on grafana side (service down, etc.). At the same time
		we can not guarantee, that data created in grafana would successfully
		stored into db.

		To resolve such kind of issue - following approach is taken. The highest
		priority is given to database data.
		That means, that creation of visualization happens in 3 steps
		1 - create database entry for visualizations and all dashboards.
			Grafana slug field is left empty
		2 - create grafana entries via grafana api, get slugs as the result
		3 - update database entries with grafana slugs
	*/

	uploadedGrafanaSlugs := []string{}

	log.Logger.Debug("Uploading dashboard data to grafana")
	for _, renderedTemplate := range renderedTemplates {
		slug, grafanaUploadErr := clients.Grafana.UploadDashboard(
			[]byte(renderedTemplate), organizationID, false)
		if grafanaUploadErr != nil {
			// We can not create grafana dashboard using user-provided template
			log.Logger.Errorf("Error during performing grafana call "+
				" for dashboard upload %s", grafanaUploadErr)
			log.Logger.Debugf("Due to error '%s' - already created grafana "+
				" dashboards, matching the same visualization, would be deleted",
				grafanaUploadErr)

			updateDashboardsDB := []*models.Dashboard{}
			deleteDashboardsDB := []*models.Dashboard{}
			for index, slugToDelete := range uploadedGrafanaSlugs {
				grafanaDeletionErr := clients.Grafana.DeleteDashboard(slugToDelete, organizationID)
				// if already created dashboard was failed to delete -
				// corresponding db entry has to be updated with grafanaSlug
				// to guarantee consistency
				if grafanaDeletionErr != nil {
					log.Logger.Errorf("Error during performing grafana call "+
						" for dashboard deletion %s", grafanaDeletionErr)
					dashboard := dashboardsDB[index]
					dashboard.Slug = uploadedGrafanaSlugs[index]
					updateDashboardsDB = append(
						updateDashboardsDB, dashboard)
				} else {
					log.Logger.Debug("deleted dashboard from grafana")
					deleteDashboardsDB = append(deleteDashboardsDB,
						dashboardsDB[index])
				}
			}

			// Delete dashboards, that were not uploaded to grafana
			deleteDashboardsDB = append(deleteDashboardsDB,
				dashboardsDB[len(uploadedGrafanaSlugs):]...)
			if len(updateDashboardsDB) > 0 {
				dashboardsToReturn := []*models.Dashboard{}
				dashboardsToReturn = append(dashboardsToReturn, updateDashboardsDB...)
				log.Logger.Debug("Updating db dashboards with grafana slugs")
				updateErrorDB := clients.DatabaseManager.BulkUpdateDashboard(
					updateDashboardsDB)
				if updateErrorDB != nil {
					log.Logger.Errorf("Error during cleanup on grafana upload"+
						" error '%s'. Unable to update db entities of dashboards"+
						" with slugs of corresponding grafana dashboards for"+
						"dashboards not deleted from grafana '%s'",
						grafanaUploadErr, updateErrorDB)
				}
				log.Logger.Debug("Deleting db dashboards that are not uploaded" +
					" to grafana")
				deletionErrorDB := clients.DatabaseManager.BulkDeleteDashboard(
					deleteDashboardsDB)
				if deletionErrorDB != nil {
					log.Logger.Debug("due to failed deletion operation - extend" +
						" the slice of returned dashboards to user")
					dashboardsToReturn = append(dashboardsToReturn, deleteDashboardsDB...)
					log.Logger.Errorf("Error during cleanup on grafana upload"+
						" error '%s'. Unable to delete entities of grafana "+
						"dashboards deleted from grafana '%s'",
						grafanaUploadErr, updateErrorDB)
				}
				result := VisualizationDashboardToResponse(
					visualizationDB, dashboardsToReturn)
				return result, common.NewClientError(
					"Unable to create new grafana dashboards, and remove old ones")
			}
			log.Logger.Debug("trying to delete visualization with " +
				"corresponding dashboards from database. dashboards have no " +
				"matching grafana uploads")
			visualizationDeletionErr := clients.DatabaseManager.DeleteVisualization(
				visualizationDB)
			if visualizationDeletionErr != nil {
				log.Logger.Error("Unable to delete visualization entry " +
					"from db with corresponding dashboards entries. " +
					"all entries are returned to user")
				result := VisualizationDashboardToResponse(
					visualizationDB, updateDashboardsDB)
				return result, common.NewClientError(
					"Unable to create new grafana dashboards, and remove old ones")
			}
			log.Logger.Debug("All created data was deleted both from grafana " +
				"and from database without errors. original grafana error is returned")
			return nil, grafanaUploadErr
		}
		log.Logger.Infof("Created dashboard named '%s'", slug)
		uploadedGrafanaSlugs = append(uploadedGrafanaSlugs, slug)
	}
	log.Logger.Debug("Uploaded dashboard data to grafana")

	// Positive outcome. All dashboards were created both in db and grafana
	for index := range dashboardsDB {
		dashboardsDB[index].Slug = uploadedGrafanaSlugs[index]
	}
	log.Logger.Debug("Updating db entries of dashboards with corresponding" +
		" grafana slugs")
	updateErrorDB := clients.DatabaseManager.BulkUpdateDashboard(dashboardsDB)
	if updateErrorDB != nil {
		log.Logger.Errorf("Error updating db dashboard slugs '%s'", updateErrorDB)
		return nil, err
	}

	return VisualizationDashboardToResponse(visualizationDB, dashboardsDB), nil
}

// VisualizationDelete removes visualizations
func (h *V1Visualizations) VisualizationDelete(clients *common.ClientContainer,
	organizationID, visualizationSlug string) (
	*common.VisualizationWithDashboards, error) {
	log.Logger.Debug("getting data from db matching provided string")
	visualizationDB, dashboardsDB, err := clients.DatabaseManager.GetVisualizationWithDashboardsBySlug(
		visualizationSlug, organizationID)
	log.Logger.Debug("got data from db matching provided string")

	if err != nil {
		log.Logger.Errorf("Error getting data from db: '%s'", err)
		return nil, err
	}

	if visualizationDB == nil {
		log.Logger.Errorf("User requested visualization '%s' not found in db", visualizationSlug)
		return nil, common.NewUserDataError("No visualizations found")
	}

	removedDashboardsFromGrafana := []*models.Dashboard{}
	failedToRemoveDashboardsFromGrafana := []*models.Dashboard{}
	for index, dashboardDB := range dashboardsDB {
		if dashboardDB.Slug == "" {
			// in case grafana slug is empty - just remove dashboard from db
			removedDashboardsFromGrafana = append(removedDashboardsFromGrafana,
				dashboardsDB[index])
		} else {
			log.Logger.Debugf("Removing grafana dashboard '%s'", dashboardDB.Slug)
			err = clients.Grafana.DeleteDashboard(dashboardDB.Slug, organizationID)
			if err != nil {
				failedToRemoveDashboardsFromGrafana = append(
					failedToRemoveDashboardsFromGrafana, dashboardsDB[index])
			} else {
				removedDashboardsFromGrafana = append(
					removedDashboardsFromGrafana, dashboardsDB[index])
			}
		}
	}

	if len(failedToRemoveDashboardsFromGrafana) != 0 {
		log.Logger.Debug("Deleting dashboards from db")
		deletionError := clients.DatabaseManager.BulkDeleteDashboard(
			removedDashboardsFromGrafana)
		if deletionError != nil {
			log.Logger.Error(deletionError)
		}
		log.Logger.Debug("Deleted dashboards from db")

		result := VisualizationDashboardToResponse(visualizationDB,
			failedToRemoveDashboardsFromGrafana)
		return result, common.NewClientError("failed to remove data from grafana")
	}
	log.Logger.Debugf("removing visualization '%s' from db", visualizationSlug)
	err = clients.DatabaseManager.DeleteVisualization(visualizationDB)
	if err != nil {
		log.Logger.Error()
	}
	log.Logger.Debugf("removed visualization '%s' from db", visualizationSlug)
	return VisualizationDashboardToResponse(visualizationDB, dashboardsDB), nil
}
