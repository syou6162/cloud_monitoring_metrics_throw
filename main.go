package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

type CustomMertic struct {
	Labels map[string]string `json:"labels"`
	Value  float64           `json:"value"`
}

func getMetricDescriptor(ctx context.Context, client monitoring.MetricClient, projectId, metricName string) (*metric.MetricDescriptor, error) {
	req := &monitoringpb.GetMetricDescriptorRequest{
		Name: fmt.Sprintf("projects/%s/metricDescriptors/custom.googleapis.com/%s", projectId, metricName),
	}
	desc, err := client.GetMetricDescriptor(ctx, req)
	if err != nil {
		return nil, err
	}
	return desc, nil
}

func throwMetrics(ctx context.Context, client monitoring.MetricClient, desc *metric.MetricDescriptor, labels map[string]string, value float64) error {
	pt := monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{Int64Value: int64(value)},
		},
	}

	ts := monitoringpb.TimeSeries{
		Metric: &metric.Metric{
			Type:   desc.Type,
			Labels: labels,
		},
		Resource: &monitoredres.MonitoredResource{
			Type: "global",
		},
		MetricKind: desc.MetricKind,
		ValueType:  desc.ValueType,
		Points:     []*monitoringpb.Point{&pt},
	}
	req := &monitoringpb.CreateTimeSeriesRequest{
		Name:       desc.Name,
		TimeSeries: []*monitoringpb.TimeSeries{&ts},
	}

	return client.CreateTimeSeries(ctx, req)
}

func main() {
	var (
		project    = flag.String("project", "", "GCP project")
		metricName = flag.String("metricName", "", "Name of metric")
	)
	flag.Parse()

	ctx := context.Background()
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	desc, err := getMetricDescriptor(ctx, *client, *project, *metricName)
	if err != nil {
		log.Fatal(err)
	}

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		if err := stdin.Err(); err != nil {
			log.Fatal(err)
		}

		var metric CustomMertic
		if err := json.Unmarshal(stdin.Bytes(), &metric); err != nil {
			log.Fatal(err)
		}
		if err := throwMetrics(ctx, *client, desc, metric.Labels, metric.Value); err != nil {
			log.Fatal(err)
		}
	}
}
