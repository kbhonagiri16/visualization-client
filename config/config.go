package config

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

// constants
const configName = "visualization-api"
const configDirPath = "/etc/platformvisibility/visualization-api"

/*
  CONFIG NAMES
  The following logic has 3 options to specify config value
    1 - conf file
    2 - env variable
    3 - command line argument
  Conf file values are overriden by env variables,
  Env variables are overriden by command line arguments

  Example of usage:
    given name "mysql.port":
      * config parser would split this string by dot and would try
        to find "port" option in "mysql" section
      * env variables parser would replace "." to "_" and would
        capitalize the strigs, then it would try to find the variable
        named "MYSQL_PORT"
      * command line parser would replace "." to "-" and look for
        "--mysql-port" option
*/
const logLevelConfigName = "log.level"
const logFileConfigName = "log.path"

const mysqlPortConfigName = "mysql.port"
const mysqlPasswordConfigName = "mysql.password"
const mysqlHostConfigName = "mysql.host"
const mysqlUserConfigName = "mysql.username"
const mysqlDatabaseConfigName = "mysql.database"

const grafanaURLConfigName = "grafana.url"
const grafanaUserConfigName = "grafana.username"
const grafanaPasswordConfigName = "grafana.password"

const httpPortConfigName = "http_endpoint.port"

// #nosec <- linter thinks that secret is hardcoded, in fact it is setting name
const httpSecretConfigName = "http_endpoint.jwt_secret"

const openstackAuthURLConfigName = "openstack.auth_url"
const openstackUsernameConfigName = "openstack.username"
const openstackPasswordConfigName = "openstack.password"
const openstackProjectConfigName = "openstack.project_name"
const openstackDomainConfigName = "openstack.domain_name"

// VisualizationAPIConfig is a struct that keeps all application config options
type VisualizationAPIConfig struct {
	// logging settings
	LogFilePath  string
	LogLevel     string
	ConsoleDebug bool

	// mysql settings
	MysqlHost         string
	MysqlPassword     string
	MysqlUsername     string
	MysqlDatabaseName string
	MysqlPort         int

	// http_endpoint settings
	HTTPPort  int
	JWTSecret string

	// openstack settings
	OpenstackAuthURL  string
	OpenstackUsername string
	OpenstackPassword string
	OpenstackProject  string
	OpenstackDomain   string

	// grafana settings
	GrafanaURL      string
	GrafanaUsername string
	GrafanaPassword string
}

var (
	singleToneConfig *VisualizationAPIConfig
)

var flagReplacer = strings.NewReplacer(".", "-", "_", "-")
var envReplacer = strings.NewReplacer(".", "_")

var _ = flag.String(flagReplacer.Replace(logFileConfigName), "",
	"Path to log file")
var _ = flag.Int(flagReplacer.Replace(mysqlPortConfigName), 0,
	"Port mysql server is listening on")
var _ = flag.String(flagReplacer.Replace(mysqlPasswordConfigName), "",
	"Password to authenticate on mysql server")
var _ = flag.String(flagReplacer.Replace(mysqlHostConfigName), "",
	"Host mysql server is running on")
var _ = flag.String(flagReplacer.Replace(mysqlUserConfigName), "",
	"Username to authenticate on mysql server")
var _ = flag.String(flagReplacer.Replace(mysqlDatabaseConfigName), "",
	"Database to use on mysql server")
var _ = flag.String(flagReplacer.Replace(grafanaURLConfigName), "",
	"URL for Grafana server")
var _ = flag.String(flagReplacer.Replace(grafanaUserConfigName), "",
	"Username for Grafana server")
var _ = flag.String(flagReplacer.Replace(grafanaPasswordConfigName), "",
	"Password for Grafana server")
var _ = flag.Bool("debug", false, "display debug messages in stdout")
var _ = flag.Int(flagReplacer.Replace(httpPortConfigName), 0,
	"Port to serve http API")
var _ = flag.String(flagReplacer.Replace(httpSecretConfigName), "",
	"Secret to use for JsonWebToken signature")

var _ = flag.String(flagReplacer.Replace(openstackAuthURLConfigName), "",
	"Auth url of openstack keystone")
var _ = flag.String(flagReplacer.Replace(openstackUsernameConfigName), "",
	"Username to auth in openstack keystone")
var _ = flag.String(flagReplacer.Replace(openstackPasswordConfigName), "",
	"Password to auth in openstack keystone")
var _ = flag.String(flagReplacer.Replace(openstackProjectConfigName), "",
	"Project name to auth in openstack keystone")
var _ = flag.String(flagReplacer.Replace(openstackDomainConfigName), "",
	"Domain name to auth in openstack keystone")

func initializeCommandLineFlags() error {

	flagsToBind := []string{
		logFileConfigName,
		mysqlHostConfigName,
		mysqlPortConfigName,
		mysqlUserConfigName,
		mysqlDatabaseConfigName,
		mysqlPasswordConfigName,
		grafanaURLConfigName,
		grafanaUserConfigName,
		grafanaPasswordConfigName,
		httpPortConfigName,
		httpSecretConfigName,
		openstackAuthURLConfigName,
		openstackUsernameConfigName,
		openstackPasswordConfigName,
		openstackProjectConfigName,
		openstackDomainConfigName,
	}
	for _, configName := range flagsToBind {
		err := viper.BindPFlag(configName, flag.Lookup(
			flagReplacer.Replace(configName)))
		if err != nil {
			return err
		}
	}

	err := viper.BindPFlag("logging.consoleDebug", flag.Lookup(
		"debug"))

	return err
}

func parseLoggingValues() error {
	// parseLoggingValues aggregates logic for getting values from viper
	// framework. Parsing options by groups reduces cyclomatic complexity of
	// InitializeConfig function
	logFileConfigValue := viper.GetString(
		logFileConfigName)
	if logFileConfigValue == "" {
		return NewParseError(
			"logPath", "path", "log", "LOG_PATH", "--log-path")
	}
	singleToneConfig.LogFilePath = logFileConfigValue

	logLevelConfigValue := viper.GetString(
		logLevelConfigName)
	if logLevelConfigValue == "" {
		return NewParseError(
			"logLevel", "level", "log", "LOG_LEVEL", "")
	}
	singleToneConfig.LogLevel = logLevelConfigValue

	return nil
}

func parseMysqlValues() error {
	// parseLoggingValues aggregates logic for getting values from viper
	// framework. Parsing options by groups reduces cyclomatic complexity of
	// InitializeConfig function
	mysqlHostConfigValue := viper.GetString(
		mysqlHostConfigName)
	if mysqlHostConfigValue == "" {
		return NewParseError(
			"mysqlHost", "host", "mysql", "MYSQL_HOST", "--mysql-host")
	}
	singleToneConfig.MysqlHost = mysqlHostConfigValue

	mysqlUserConfigValue := viper.GetString(
		mysqlUserConfigName)
	if mysqlUserConfigValue == "" {
		return NewParseError(
			"mysqlUser", "username", "mysql", "MYSQL_USERNAME", "--mysql-username")
	}
	singleToneConfig.MysqlUsername = mysqlUserConfigValue

	mysqlPasswordConfigValue := viper.GetString(
		mysqlPasswordConfigName)
	if mysqlPasswordConfigValue == "" {
		return NewParseError(
			"MysqlPassword", "password", "mysql", "MYSQL_PASSWORD", "--mysql-password")
	}
	singleToneConfig.MysqlPassword = mysqlPasswordConfigValue

	mysqlDatabaseConfigValue := viper.GetString(
		mysqlDatabaseConfigName)
	if mysqlDatabaseConfigValue == "" {
		return NewParseError(
			"MysqlDatabase", "database", "mysql", "MYSQL_DATABASE", "--mysql-database")
	}
	singleToneConfig.MysqlDatabaseName = mysqlDatabaseConfigValue

	mysqlPortConfigValue := viper.GetInt(
		mysqlPortConfigName)
	if mysqlPortConfigValue == 0 {
		return NewParseError(
			"MysqlPort", "port", "mysql", "MYSQL_PORT", "--mysql-port")
	}
	singleToneConfig.MysqlPort = mysqlPortConfigValue

	return nil
}

func parseHTTPEndpointValues() error {
	httpPortConfigValue := viper.GetInt(
		httpPortConfigName)
	if httpPortConfigValue == 0 {
		return NewParseError(
			"httpEndpointPort", "port", "http_endpoint", "HTTP_ENDPOINT_PORT",
			"--http-endpoint-port")
	}
	singleToneConfig.HTTPPort = httpPortConfigValue

	httpSecretConfigValue := viper.GetString(
		httpSecretConfigName)
	if httpSecretConfigValue == "" {
		return NewParseError(
			"httpEndpointSecret", "secret", "http_endpoint",
			"HTTP_ENDPOINT_JWT_SECRET", "--http-endpoint-jwt-secret")
	}
	singleToneConfig.JWTSecret = httpSecretConfigValue

	return nil
}

func parseOpenstackValues() error {
	openstackAuthURLConfigValue := viper.GetString(
		openstackAuthURLConfigName)
	if openstackAuthURLConfigValue == "" {
		return NewParseError(
			"OpenstackAuthURL", "auth_url", "openstack",
			"OPENSTACK_AUTH_URL", "--openstack-auth-url")
	}
	singleToneConfig.OpenstackAuthURL = openstackAuthURLConfigValue

	openstackPasswordConfigValue := viper.GetString(
		openstackPasswordConfigName)
	if openstackPasswordConfigValue == "" {
		return NewParseError(
			"OpenstackPassword", "password", "openstack",
			"OPENSTACK_PASSWORD", "--openstack-password")
	}
	singleToneConfig.OpenstackPassword = openstackPasswordConfigValue

	openstackUsernameConfigValue := viper.GetString(
		openstackUsernameConfigName)
	if openstackUsernameConfigValue == "" {
		return NewParseError(
			"OpenstackUsername", "username", "openstack",
			"OPENSTACK_USERNAME", "--openstack-username")
	}
	singleToneConfig.OpenstackUsername = openstackUsernameConfigValue

	openstackProjectConfigValue := viper.GetString(
		openstackProjectConfigName)
	if openstackProjectConfigValue == "" {
		return NewParseError(
			"OpenstackProject", "project", "openstack",
			"OPENSTACK_PROJECT_NAME", "--openstack-project-name")
	}
	singleToneConfig.OpenstackProject = openstackProjectConfigValue

	openstackDomainConfigValue := viper.GetString(
		openstackDomainConfigName)
	if openstackDomainConfigValue == "" {
		return NewParseError(
			"OpenstackDomain", "domain_name", "openstack",
			"OPENSTACK_DOMAIN_NAME", "--openstack-domain-name")
	}
	singleToneConfig.OpenstackDomain = openstackDomainConfigValue

	return nil
}

func parseGrafanaValues() error {
	// parseLoggingValues aggregates logic for getting values from viper
	// framework. Parsing options by groups reduces cyclomatic complexity of
	// InitializeConfig function
	grafanaURLConfigValue := viper.GetString(
		grafanaURLConfigName)
	if grafanaURLConfigValue == "" {
		return NewParseError(
			"grafanaUrl", "url", "grafana", "GRAFANA_URL", "--grafana-url")
	}
	singleToneConfig.GrafanaURL = grafanaURLConfigValue

	grafanaUserConfigValue := viper.GetString(
		grafanaUserConfigName)
	if grafanaUserConfigValue == "" {
		return NewParseError(
			"grafanaUser", "username", "grafana", "GRAFANA_USERNAME", "--grafana-username")
	}
	singleToneConfig.GrafanaUsername = grafanaUserConfigValue

	grafanaPasswordConfigValue := viper.GetString(
		grafanaPasswordConfigName)
	if grafanaPasswordConfigValue == "" {
		return NewParseError(
			"grafanaPassword", "password", "grafana", "GRAFANA_PASSSWORD", "--grafana-password")
	}
	singleToneConfig.GrafanaPassword = grafanaPasswordConfigValue

	return nil
}

// InitializeConfig parses application configuration from config file, env
// variables and console flags. parsed configs are stored in module level variable
func InitializeConfig() error {
	var err error

	// assign address of default initialized VisualizationApiConfig
	// to global pointer
	singleToneConfig = &VisualizationAPIConfig{}

	// initialize path to conf files
	viper.SetConfigName(configName)
	viper.AddConfigPath(configDirPath)

	// set automatic parse of environment variables
	// env variables have higher priority then config file values,
	// but are overriden with command line flags
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(envReplacer)

	// initialize viper to read parameters passed by commandline via
	// argv[], commandline variables have higher priority
	// then env variables and config file values
	err = initializeCommandLineFlags()
	if err != nil {
		return err
	}

	// parse configs using viper lib
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	// parsing of values from viper framework to struct is split to groups to
	// reduce cyclomatic complexity of InitializeConfig function
	err = parseLoggingValues()
	if err != nil {
		return err
	}
	err = parseMysqlValues()
	if err != nil {
		return err
	}
	err = parseGrafanaValues()
	if err != nil {
		return err
	}
	err = parseHTTPEndpointValues()
	if err != nil {
		return err
	}
	err = parseOpenstackValues()
	if err != nil {
		return err
	}

	// console debug has default values - no need to check
	singleToneConfig.ConsoleDebug = viper.GetBool(
		"logging.consoleDebug")
	return nil
}

// GetConfig returns structure with all application configuration
func GetConfig() *VisualizationAPIConfig {
	return singleToneConfig
}
