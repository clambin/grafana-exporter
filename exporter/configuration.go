package exporter

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"path/filepath"
	"strings"
)

// Configuration holds all parameters retrieved from the command line
type Configuration struct {
	Debug     bool
	URL       string
	Token     string
	Out       string
	Direct    bool
	Namespace string
	Command   string
	Folders   []string
}

const (
	cmdDashboards            = "dashboards"
	cmdDataSources           = "datasources"
	cmdDashboardProvisioning = "dashboard-provisioning"
)

// GetConfiguration parses the command line arguments. If showErrors is true, any errors are displayed
func GetConfiguration(args []string, showErrors bool) (cfg *Configuration, err error) {
	cfg = new(Configuration)

	app := kingpin.New(filepath.Base(args[0]), "grafana provisioning export")
	app.Flag("debug", "Log debug messages").Short('d').BoolVar(&cfg.Debug)
	app.Flag("url", "Grafana API URL").Short('u').Required().StringVar(&cfg.URL)
	app.Flag("token", "Grafana API Token (must have admin access)").Short('t').Required().StringVar(&cfg.Token)
	app.Flag("out", "Output directory").Short('o').Default(".").StringVar(&cfg.Out)
	app.Flag("direct", "Write config files directory (default: write as K8S ConfigMaps").BoolVar(&cfg.Direct)
	app.Flag("namespace", "K8s Namespace").Short('n').Default("monitoring").StringVar(&cfg.Namespace)

	app.Command(cmdDataSources, "export Grafana data sources")
	app.Command(cmdDashboardProvisioning, "export Grafana dashboard provisioning data")

	dashboards := app.Command(cmdDashboards, "export Grafana dashboards")
	folders := dashboards.Flag("folders", "Comma-separated list of folders to export").Short('f').String()

	cfg.Command, err = app.Parse(args[1:])

	if err == nil {
		if *folders != "" {
			cfg.Folders = strings.Split(*folders, ",")
		}
	}

	if err != nil && showErrors {
		app.Usage(args[1:])
	}

	return
}
