package prometheusconvert

import (
	"time"

	"github.com/grafana/agent/component/prometheus/remotewrite"
	"github.com/grafana/agent/converter/internal/common"
	prom_config "github.com/prometheus/prometheus/config"
)

func appendPrometheusRemoteWrite(pb *prometheusBlocks, promConfig *prom_config.Config) *remotewrite.Exports {
	remoteWriteArgs := toRemotewriteArguments(promConfig)
	block := common.NewBlockWithOverride([]string{"prometheus", "remote_write"}, "default", remoteWriteArgs)
	pb.prometheusRemoteWriteBlocks = append(pb.prometheusRemoteWriteBlocks, block)
	return &remotewrite.Exports{
		Receiver: common.ConvertAppendable{Expr: "prometheus.remote_write.default.receiver"},
	}
}

func toRemotewriteArguments(promConfig *prom_config.Config) *remotewrite.Arguments {
	externalLabels := promConfig.GlobalConfig.ExternalLabels.Map()
	if len(externalLabels) == 0 {
		externalLabels = nil
	}

	return &remotewrite.Arguments{
		ExternalLabels: externalLabels,
		Endpoints:      GetEndpointOptions(promConfig.RemoteWriteConfigs),
		WALOptions:     remotewrite.DefaultWALOptions,
	}
}

func GetEndpointOptions(remoteWriteConfigs []*prom_config.RemoteWriteConfig) []*remotewrite.EndpointOptions {
	endpoints := make([]*remotewrite.EndpointOptions, 0)

	for _, remoteWriteConfig := range remoteWriteConfigs {
		endpoint := &remotewrite.EndpointOptions{
			Name:                 remoteWriteConfig.Name,
			URL:                  remoteWriteConfig.URL.String(),
			RemoteTimeout:        time.Duration(remoteWriteConfig.RemoteTimeout),
			Headers:              remoteWriteConfig.Headers,
			SendExemplars:        remoteWriteConfig.SendExemplars,
			SendNativeHistograms: remoteWriteConfig.SendNativeHistograms,
			HTTPClientConfig:     toHttpClientConfig(&remoteWriteConfig.HTTPClientConfig),
			QueueOptions:         toQueueOptions(&remoteWriteConfig.QueueConfig),
			MetadataOptions:      toMetadataOptions(&remoteWriteConfig.MetadataConfig),
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

func toQueueOptions(queueConfig *prom_config.QueueConfig) *remotewrite.QueueOptions {
	return &remotewrite.QueueOptions{
		Capacity:          queueConfig.Capacity,
		MaxShards:         queueConfig.MaxShards,
		MinShards:         queueConfig.MinShards,
		MaxSamplesPerSend: queueConfig.MaxSamplesPerSend,
		BatchSendDeadline: time.Duration(queueConfig.BatchSendDeadline),
		MinBackoff:        time.Duration(queueConfig.MinBackoff),
		MaxBackoff:        time.Duration(queueConfig.MaxBackoff),
		RetryOnHTTP429:    queueConfig.RetryOnRateLimit,
	}
}

func toMetadataOptions(metadataConfig *prom_config.MetadataConfig) *remotewrite.MetadataOptions {
	return &remotewrite.MetadataOptions{
		Send:              metadataConfig.Send,
		SendInterval:      time.Duration(metadataConfig.SendInterval),
		MaxSamplesPerSend: metadataConfig.MaxSamplesPerSend,
	}
}
