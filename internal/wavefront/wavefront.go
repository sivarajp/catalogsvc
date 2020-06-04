package wavefront

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
	"github.com/sivarajp/catalogsvc/pkg/logger"
	"github.com/wavefronthq/go-metrics-wavefront/reporting"
	"github.com/wavefronthq/wavefront-sdk-go/application"
	"github.com/wavefronthq/wavefront-sdk-go/senders"
)

func InitWavefront() reporting.WavefrontMetricsReporter {

	// Create a direct sender
	directCfg := &senders.DirectConfiguration{
		Server:               "https://demo.wavefront.com",
		Token:                "400a1c6d-247f-46c9-acd7-53306b58525f",
		BatchSize:            10000,
		MaxBufferSize:        50000,
		FlushIntervalSeconds: 5,
	}

	sender, err := senders.NewDirectSender(directCfg)
	if err != nil {
		panic(err)
	}

	reporter := reporting.NewReporter(
		sender,
		application.New("siva-demo-app", "catalog-service"),
		reporting.Source("siva-dev"),
		reporting.Prefix("siva.acme"),
		reporting.LogErrors(true),
		reporting.RuntimeMetric(true),
	)

	return reporter
}

func WavefrontEmitter(reporter reporting.WavefrontMetricsReporter) gin.HandlerFunc {

	return func(g *gin.Context) {
		logger.Logger.Infof("Inside middleware")
		pointTags := make(map[string]string)
		// Start timer
		start := time.Now()

		// Process request
		g.Next()

		// Stop timer
		end := time.Now()

		logger.Logger.Infof("Time difference between calls%d", end.Sub(start).Milliseconds)
		// latency := end.Sub(start)
		statusCode := g.Writer.Status()
		// bytesOut := c.Writer.Size()
		// bytesIn := c.Request.ContentLength

		pointTags["path"] = g.Request.URL.Path
		pointTags["clientIP"] = g.ClientIP()
		pointTags["method"] = g.Request.Method
		pointTags["userAgent"] = g.Request.UserAgent()

		logger.Logger.Infof("tags", pointTags)
		var c interface{}
		var apiType string
		if g.Request.Method == "GET" {
			if g.Request.URL.Path == "/products" {
				apiType := "ListProduct"
				c = reporter.GetMetric(apiType, pointTags)
				if c == nil {
					c = metrics.NewCounter()
					reporter.RegisterMetric(apiType, c, pointTags)
				}
			} else {
				apiType := "GetProduct"
				c = reporter.GetMetric(apiType, pointTags)
				if c == nil {
					c = metrics.NewCounter()
					reporter.RegisterMetric(apiType, c, pointTags)
				}
			}
		} else if g.Request.Method == "POST" {
			apiType := "CreateProduct"
			c = reporter.GetMetric(apiType, pointTags)
			if c == nil {
				c = metrics.NewCounter()
				reporter.RegisterMetric(apiType, c, pointTags)
			}
		}
		c.(metrics.Counter).Inc(1)
		fmt.Println(apiType)
		// reporter.GetOrRegisterMetric("m1", getCounter, map[string]string{"tag1": "tag"})
		// reporter.GetOrRegisterMetric("m2", createCounter, map[string]string{"application": "tag"})

		// // Send metrics
		// // <metricName> <metricValue> [<timestamp>] source=<source> [pointTags]
		// sender.SendMetric(strings.Join([]string{w.MetricPrefix, ".latency"}, ""), float64(latency.Milliseconds()), end.Unix(), w.Source, w.PointTags)
		// sender.SendMetric(strings.Join([]string{w.MetricPrefix, ".bytes.in"}, ""), float64(bytesIn), end.Unix(), w.Source, w.PointTags)
		// sender.SendMetric(strings.Join([]string{w.MetricPrefix, ".bytes.out"}, ""), float64(bytesOut), end.Unix(), w.Source, w.PointTags)
		// switch {
		// case statusCode > 199 && statusCode < 300:
		// 	sender.SendDeltaCounter(strings.Join([]string{w.MetricPrefix, ".status.success"}, ""), 1, w.Source, w.PointTags)
		// case statusCode > 299 && statusCode < 400:
		// 	sender.SendDeltaCounter(strings.Join([]string{w.MetricPrefix, ".status.redirection"}, ""), 1, w.Source, w.PointTags)
		// case statusCode > 399 && statusCode < 500:
		// 	sender.SendDeltaCounter(strings.Join([]string{w.MetricPrefix, ".status.error.client"}, ""), 1, w.Source, w.PointTags)
		// case statusCode > 499 && statusCode < 600:
		// 	sender.SendDeltaCounter(strings.Join([]string{w.MetricPrefix, ".status.error.server"}, ""), 1, w.Source, w.PointTags)
		// }
	}

}
