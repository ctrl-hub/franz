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
	go collect()
	ticker := time.NewTicker(5 * time.Second)
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
	logrus.WithField("port", viper.Get("metrics_port")).WithField("path", viper.GetString("metrics_path")).Info("starting metrics server")
	http.Handle(viper.GetString("metrics_path"), handler)
	return http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("metrics_port")), nil)
}
