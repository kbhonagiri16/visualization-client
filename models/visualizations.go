package models

// Visualization represents visualization in db
type Visualization struct {
	ID             int    `xorm:"autoincr pk 'id'"`
	Slug           string `xorm:"slug"`
	Name           string `xorm:"name"`
	OrganizationID string `xorm:"organization_id"`
	Tags           string `xorm:"tags"`
}

// Dashboard represents dashboard in db
type Dashboard struct {
	ID               string `xorm:"pk 'id'"`
	Visualization    int    `xorm:"visualization_id"`
	Name             string `xorm:"name"`
	RenderedTemplate string `xorm:"rendered_template"`
	Slug             string `xorm:"slug"`
}

// DashboardTableName describes database table name (not to use reflect)
const DashboardTableName = "dashboard"

// DashboardVisualizationColumn describes database column name (not to use reflect)
const DashboardVisualizationColumn = "visualization_id"

// VisualizationTableName describes database table name (not to use reflect)
const VisualizationTableName = "visualization"

// VisualizationIDColumn describes database column name (not to use reflect)
const VisualizationIDColumn = "id"

// VisualizationTagsColumn describes database column name (not to use reflect)
const VisualizationTagsColumn = "tags"

// VisualizationSlugColumn describes database column name (not to use reflect)
const VisualizationSlugColumn = "slug"

// VisualizationNameColumn describes database column name (not to use reflect)
const VisualizationNameColumn = "name"

// VisualizationOrgColumn describes database column name (not to use reflect)
const VisualizationOrgColumn = "organization_id"
