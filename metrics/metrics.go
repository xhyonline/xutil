package metrics

import (
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/xhyonline/xutil/logger"
)

func Init(gatewayURL, serviceName string) {
	if gatewayURL == "" || serviceName == "" {
		logger.Errorf("metrics push url or service_name undefined.")
		return
	}
	timer := time.NewTicker(1 * time.Minute)
	// get hostname
	hostname, _ := os.Hostname()
	labels := map[string]string{
		"service_type": "golang",
		"hostname":     hostname,
	}

	go func() {
		for {
			<-timer.C

			goC := collectors.NewGoCollector()
			processC := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})

			pusher := push.New(gatewayURL, serviceName).Collector(goC).Collector(processC)
			for k, v := range labels {
				pusher = pusher.Grouping(k, v)
			}
			if err := pusher.Add(); err != nil {
				logger.Error(err)
			}
		}
	}()
}
