package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Serve() error {

	// create the prometheus handler
	handler := metricsHandler()

	// run our first metric collection on start, then at a predefined tick
	//nolint
	go collect()
	ticker := time.NewTicker(time.Duration(viper.GetInt("polling_interval_seconds")) * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err := collect()
				if err != nil {
					logrus.WithError(err).Error("failed to collect metrics")
				}
			}
		}
	}()

	// start the metrics server
	logrus.
		WithField("poll_interval", fmt.Sprintf("%ds", viper.GetInt("polling_interval_seconds"))).
		WithField("port", viper.Get("metrics_port")).
		WithField("path", viper.GetString("metrics_path")).
		Info("starting metrics server")

	http.Handle(viper.GetString("metrics_path"), handler)
	return http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("metrics_port")), nil)
}
