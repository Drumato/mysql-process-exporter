package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MySQLProcessSecondsGaugeVec *prometheus.GaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "mysql_process_exporter",
		Name:      "mysql_process_seconds",
		Help:      "MySQL process seconds",
	}, []string{
		"db_host",
		"id",
		"user",
		"host",
		"db",
		"command",
		"state",
		"info",
	})
)

// When updating metrics, it is not clear what labels should be used, so we pass them via a struct to make it easier to understand.
type MySQLProcessSecondsGaugeVecLabels struct {
	// DBHost is the host of the database
	DBHost string
	// ID is the ID of the process
	ID string
	// User is the user of the process
	User string
	// Host is the host of the process
	Host string
	// DB is the database of the process
	DB string
	// Command is the command of the process
	Command string
	// State is the state of the process
	State string
	// Info is the info of the process
	Info string
}

func InitializeMetrics() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(MySQLProcessSecondsGaugeVec)
	return reg
}

func UpdateMySQLProcessSecondsGaugeVec(labels MySQLProcessSecondsGaugeVecLabels, value float64) {
	MySQLProcessSecondsGaugeVec.With(prometheus.Labels{
		"db_host": labels.DBHost,
		"id":      labels.ID,
		"user":    labels.User,
		"host":    labels.Host,
		"db":      labels.DB,
		"command": labels.Command,
		"state":   labels.State,
		"info":    labels.Info,
	}).Set(value)
}
