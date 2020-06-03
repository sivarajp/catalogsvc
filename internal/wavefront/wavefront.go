package wavefront

import (
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
		FlushIntervalSeconds: 1,
	}

	sender, err := senders.NewDirectSender(directCfg)
	if err != nil {
		panic(err)
	}

	reporter := reporting.NewReporter(
		sender,
		application.New("siva-demo-app", "test-service"),
		reporting.Source("go-metrics-test"),
		reporting.Prefix("siva.prefix"),
		reporting.LogErrors(true),
		reporting.RuntimeMetric(true),
	)

	return reporter
}
