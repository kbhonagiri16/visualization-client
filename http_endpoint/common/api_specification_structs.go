package common

// VisualizationPOSTData - POST data expected by visualization api
type VisualizationPOSTData struct {
	Name       string `json:"name"`
	Dashboards []struct {
		Name               string      `json:"name"`
		TemplateName       string      `json:"templateName"`
		TemplateVersion    int         `json:"templateVersion"`
		TemplateBody       string      `json:"templateBody"`
		TemplateParameters interface{} `json:"templateParameters"`
	} `json:"dashboards"`
	Tags map[string]interface{} `json:"tags"`
}

// VisualizationResponseEntry describes what data would be returned to user
type VisualizationResponseEntry struct {
	Slug string `json:"id"`
	Name string `json:"name"`
	Tags string `json:"tags"`
}

// DashboardResponseEntry describes what data would be returned to user
type DashboardResponseEntry struct {
	Name             string `json:"name"`
	RenderedTemplate string `json:"renderedTemplate"`
	Slug             string `json:"id"`
}

// VisualizationWithDashboards aggregates VisualizationResponseEntry and DashboardResponseEntry
type VisualizationWithDashboards struct {
	*VisualizationResponseEntry
	Dashboards []*DashboardResponseEntry `json:"dashboards"`
}
