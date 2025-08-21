package telemetry

import "backend-poc/backoffice/extensions/configuration"

type OpenTelemetry struct {
	ApplicationName string
	Namespace       string
	InstanceID      string
	Version         string
	Endpoint        string
	Environment     string
	SamplingRate    float64
}

type OpenTelemetryOption func(*OpenTelemetry)

func NewOpenTelemetry(opts ...OpenTelemetryOption) OpenTelemetry {
	parsedConfig, _ := configuration.Load[config]()()

	applicationName := parsedConfig.ServiceName

	if applicationName == "" {
		applicationName = parsedConfig.ApplicationName
	}

	environment := parsedConfig.OtelEnv

	if environment == "" {
		environment = parsedConfig.Env
	}

	openTelemetry := OpenTelemetry{
		ApplicationName: applicationName,
		Namespace:       parsedConfig.Namespace,
		InstanceID:      parsedConfig.InstanceID,
		Version:         parsedConfig.Version,
		Endpoint:        parsedConfig.Endpoint,
		Environment:     environment,
		SamplingRate:    parsedConfig.SamplingRate,
	}

	for _, opt := range opts {
		opt(&openTelemetry)
	}

	return openTelemetry
}

func ApplicationName(applicationName string) OpenTelemetryOption {
	return func(o *OpenTelemetry) {
		o.ApplicationName = applicationName
	}
}

type config struct {
	ServiceName     string  `conf:"env:OTEL_SERVICE_NAME"`
	ApplicationName string  `conf:"env:APPLICATION_NAME"`
	Namespace       string  `conf:"env:OTEL_NAMESPACE,default:default"`
	InstanceID      string  `conf:"env:HOSTNAME,default:#1"`
	Version         string  `conf:"env:OTEL_SERVICE_VERSION"`
	Endpoint        string  `conf:"env:OTEL_EXPORTER_OTLP_ENDPOINT,default:localhost:4317"`
	OtelEnv         string  `conf:"env:OTEL_DEPLOYMENT_ENVIRONMENT"`
	Env             string  `conf:"env:ENV,default:dev"`
	SamplingRate    float64 `conf:"env:OTEL_SAMPLING_RATIO,default:0.1"`
}
