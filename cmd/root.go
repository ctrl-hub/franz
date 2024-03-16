package cmd

import (
	"strings"

	"github.com/ctrl-hub/franz/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the kakfa metric exporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.Serve()
		},
	}
)

func Execute() {
	rootCmd.Flags().String("log_level", "info", "The log level to use: fatal|panic|warn|info|debug")
	viper.BindPFlag("log_level", rootCmd.Flags().Lookup("log_level"))
	viper.BindEnv("log_level")

	rootCmd.Flags().String("log_format", "logfmt", "How you want your logs formatted: logfmt|json")
	viper.BindPFlag("log_format", rootCmd.Flags().Lookup("log_format"))
	viper.BindEnv("log_format")

	rootCmd.Flags().String("metrics_path", "/metrics", "The path to serve metrics on")
	viper.BindPFlag("metrics_path", rootCmd.Flags().Lookup("metrics_path"))
	viper.BindEnv("metrics_path")

	rootCmd.Flags().Int("metrics_port", 3100, "Port to run the metrics server on")
	viper.BindPFlag("metrics_port", rootCmd.Flags().Lookup("metrics_port"))
	viper.BindEnv("metrics_port")

	rootCmd.Flags().Int("polling_interval_seconds", 10, "How often to poll the kafka cluster (in seconds)")
	viper.BindPFlag("polling_interval_seconds", rootCmd.Flags().Lookup("polling_interval_seconds"))
	viper.BindEnv("polling_interval_seconds")

	rootCmd.Flags().String("confluent_cluster_label", "confluent", "The label value to add for the cluster label in metrics")
	viper.BindPFlag("confluent_cluster_label", rootCmd.Flags().Lookup("confluent_cluster_label"))
	viper.BindEnv("confluent_cluster_label")

	rootCmd.Flags().String("confluent_endpoint", "", "The confluent endpoint for the cluster (e.g. `pkc-xxxx.region.provider.onfluent.cloud:9092`)")
	viper.BindPFlag("confluent_endpoint", rootCmd.Flags().Lookup("confluent_endpoint"))
	viper.BindEnv("confluent_endpoint")

	rootCmd.Flags().String("confluent_api_key", "", "The confluent API Key for SASL authentication (username)")
	viper.BindPFlag("confluent_api_key", rootCmd.Flags().Lookup("confluent_api_key"))
	viper.BindEnv("confluent_api_key")

	rootCmd.Flags().String("confluent_api_secret", "", "The confluent API Secret for SASL authentication (password)")
	viper.BindPFlag("confluent_api_secret", rootCmd.Flags().Lookup("confluent_api_secret"))
	viper.BindEnv("confluent_api_secret")

	level := logrus.DebugLevel
	switch strings.ToLower(viper.GetString("log_level")) {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "panic":
		level = logrus.PanicLevel
	case "fatal":
		level = logrus.FatalLevel
	}
	logrus.SetLevel(level)

	switch strings.ToLower(viper.GetString("log_format")) {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}

	if viper.Get("confluent_endpoint") == "" {
		logrus.WithField("flag", "confluent_endpoint").WithField("help", "Set the --confluent_endpoint flag or CONFLUENT_ENDPOINT env var").Fatal("missing value in config")
	}
	if viper.Get("confluent_api_key") == "" {
		logrus.WithField("flag", "confluent_api_key").WithField("help", "Set the --confluent_api_key flag or CONFLUENT_API_KEY env var").Fatal("missing value in config")
	}
	if viper.Get("confluent_api_secret") == "" {
		logrus.WithField("flag", "confluent_api_secret").WithField("help", "Set the --confluent_api_secret flag or CONFLUENT_API_SECRET env var").Fatal("missing value in config")
	}

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("could not start")
	}
}
