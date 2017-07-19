package v1handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ulule/deepcopier"
	"visualization-client"
	"visualization-client/http_endpoint/common"
)

// V1Visualizations implements part of handler interface
type V1Visualizations struct{}

// VisualizationDashboardToResponse transforms models to response format
func VisualizationDashboardToResponse(visualization *visualization.Visualization,
	dashboards []*visualization.Dashboard) *common.VisualizationWithDashboards {
	// This function is used, when we have to return visualization with
	// limited number of dashboards (for example in post method)
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
	data *map[visualization.Visualization][]*visualization.Dashboard) *[]common.VisualizationWithDashboards {
	// This function is used, when

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

	data, err := clients.DatabaseManager.QueryVisualizationsDashboards(
		"", name, organizationID, tags)
	if err != nil {
		return nil, err
	}

	return GroupedVisualizationDashboardToResponse(data), nil
}

func renderTemplates(templates []string, templateParamaters []interface{}) (
	[]string, error) {
	// this function takes Visualization data and returns rendered templates
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

	templates := []string{}
	templateParamaters := []interface{}{}
	dashboardNames := []string{}
	for _, dashboardData := range data.Dashboards {
		templates = append(templates, dashboardData.TemplateBody)
		templateParamaters = append(templateParamaters, dashboardData.TemplateParameters)
		dashboardNames = append(dashboardNames, dashboardData.Name)
	}

	renderedTemplates, err := renderTemplates(templates, templateParamaters)
	if err != nil {
		return nil, err
	}

	// create db entries for visualizations and dashboards
	visualizationDB, dashboardsDB, err := clients.DatabaseManager.CreateVisualizationsWithDashboards(
		data.Name, organizationID, data.Tags, dashboardNames, renderedTemplates)
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

	for _, renderedTemplate := range renderedTemplates {
		slug, grafanaUploadErr := clients.Grafana.UploadDashboard(
			[]byte(renderedTemplate), organizationID, false)
		if grafanaUploadErr != nil {
			// We can not create grafana dashboard using user-provided template

			updateDashboardsDB := []*visualization.Dashboard{}
			deleteDashboardsDB := []*visualization.Dashboard{}
			for index, slugToDelete := range uploadedGrafanaSlugs {
				grafanaDeletionErr := clients.Grafana.DeleteDashboard(slugToDelete, organizationID)
				// if already created dashboard was failed to delete -
				// corresponding db entry has to be updated with grafanaSlug
				// to guarantee consistency
				if grafanaDeletionErr != nil {
					dashboard := dashboardsDB[index]
					dashboard.Slug = uploadedGrafanaSlugs[index]
					updateDashboardsDB = append(
						updateDashboardsDB, dashboard)
				} else {
					deleteDashboardsDB = append(deleteDashboardsDB,
						dashboardsDB[index])
				}
			}

			// Delete dashboards, that were not uploaded to grafana
			deleteDashboardsDB = append(deleteDashboardsDB,
				dashboardsDB[len(uploadedGrafanaSlugs):]...)
			if len(updateDashboardsDB) > 0 {
				dashboardsToReturn := []*visualization.Dashboard{}
				dashboardsToReturn = append(dashboardsToReturn, updateDashboardsDB...)
				updateErrorDB := clients.DatabaseManager.BulkUpdateDashboard(
					updateDashboardsDB)
				if updateErrorDB != nil {
					fmt.Printf("Error during cleanup on grafana upload"+
						" error '%s'. Unable to update db entities of dashboards"+
						" with slugs of corresponding grafana dashboards for"+
						"dashboards not deleted from grafana '%s'",
						grafanaUploadErr, updateErrorDB)
				}
				fmt.Printf("Deleting db dashboards that are not uploaded" +
					" to grafana")
				deletionErrorDB := clients.DatabaseManager.BulkDeleteDashboard(
					deleteDashboardsDB)
				if deletionErrorDB != nil {
					fmt.Printf("due to failed deletion operation - extend" +
						" the slice of returned dashboards to user")
					dashboardsToReturn = append(dashboardsToReturn, deleteDashboardsDB...)
					fmt.Printf("Error during cleanup on grafana upload"+
						" error '%s'. Unable to delete entities of grafana "+
						"dashboards deleted from grafana '%s'",
						grafanaUploadErr, updateErrorDB)
				}
				result := VisualizationDashboardToResponse(
					visualizationDB, dashboardsToReturn)
				return result, common.NewClientError(
					"Unable to create new grafana dashboards, and remove old ones")
			}
			fmt.Printf("trying to delete visualization with " +
				"corresponding dashboards from database. dashboards have no " +
				"matching grafana uploads")
			visualizationDeletionErr := clients.DatabaseManager.DeleteVisualization(
				visualizationDB)
			if visualizationDeletionErr != nil {
				fmt.Printf("Unable to delete visualization entry " +
					"from db with corresponding dashboards entries. " +
					"all entries are returned to user")
				result := VisualizationDashboardToResponse(
					visualizationDB, updateDashboardsDB)
				return result, common.NewClientError(
					"Unable to create new grafana dashboards, and remove old ones")
			}
			fmt.Printf("All created data was deleted both from grafana " +
				"and from database without errors. original grafana error is returned")
			return nil, grafanaUploadErr
		}
		fmt.Printf("Created dashboard named '%s'", slug)
		uploadedGrafanaSlugs = append(uploadedGrafanaSlugs, slug)
	}
	fmt.Printf("Uploaded dashboard data to grafana")

	// Positive outcome. All dashboards were created both in db and grafana
	for index := range dashboardsDB {
		dashboardsDB[index].Slug = uploadedGrafanaSlugs[index]
	}
	fmt.Printf("Updating db entries of dashboards with corresponding" +
		" grafana slugs")
	updateErrorDB := clients.DatabaseManager.BulkUpdateDashboard(dashboardsDB)
	if updateErrorDB != nil {
		fmt.Printf("Error updating db dashboard slugs '%s'", updateErrorDB)
		return nil, err
	}

	return VisualizationDashboardToResponse(visualizationDB, dashboardsDB), nil
}

// VisualizationDelete removes visualizations
func (h *V1Visualizations) VisualizationDelete(clients *common.ClientContainer,
	organizationID, visualizationSlug string) (
	*common.VisualizationWithDashboards, error) {
	fmt.Printf("getting data from db matching provided string")
	visualizationDB, dashboardsDB, err := clients.DatabaseManager.GetVisualizationWithDashboardsBySlug(
		visualizationSlug, organizationID)
	fmt.Printf("got data from db matching provided string")

	if err != nil {
		fmt.Printf("Error getting data from db: '%s'", err)
		return nil, err
	}

	if visualizationDB == nil {
		fmt.Printf("User requested visualization '%s' not found in db", visualizationSlug)
		return nil, common.NewUserDataError("No visualizations found")
	}

	removedDashboardsFromGrafana := []*visualization.Dashboard{}
	failedToRemoveDashboardsFromGrafana := []*visualization.Dashboard{}
	for index, dashboardDB := range dashboardsDB {
		if dashboardDB.Slug == "" {
			// in case grafana slug is empty - just remove dashboard from db
			removedDashboardsFromGrafana = append(removedDashboardsFromGrafana,
				dashboardsDB[index])
		} else {
			fmt.Printf("Removing grafana dashboard '%s'", dashboardDB.Slug)
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
		fmt.Printf("Deleting dashboards from db")
		deletionError := clients.DatabaseManager.BulkDeleteDashboard(
			removedDashboardsFromGrafana)
		if deletionError != nil {
			fmt.Println(deletionError)
		}
		fmt.Printf("Deleted dashboards from db")

		result := VisualizationDashboardToResponse(visualizationDB,
			failedToRemoveDashboardsFromGrafana)
		return result, common.NewClientError("failed to remove data from grafana")
	}
	fmt.Printf("removing visualization '%s' from db", visualizationSlug)
	err = clients.DatabaseManager.DeleteVisualization(visualizationDB)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("removed visualization '%s' from db", visualizationSlug)
	return VisualizationDashboardToResponse(visualizationDB, dashboardsDB), nil
}
