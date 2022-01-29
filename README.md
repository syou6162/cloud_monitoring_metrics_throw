## cloud_monitoring_metrics_throw
Cloud monitoring metrics is useful service for monitoring, but it is painful to use it from command line. `cloud_monitoring_metrics_throw` is a tool to post cloud monitoring metrics easily from command line.

## Install

```
go install github.com/syou6162/cloud_monitoring_metrics_throw@latest
```

## Usage
You can post metrics to cloud monitoring via STDIN.

```
cat ./sample.json | cloud_monitoring_metrics_throw --project my-project --metricName sample_metric 
```

Before posting metrics, it is it is recommended to define [metric descriptor](https://cloud.google.com/monitoring/custom-metrics/). Following is a sample definition by terraform.

```terraform
resource "google_monitoring_metric_descriptor" "sample_metric_descriptor" {
  description  = "Sample metric"
  display_name = "sample_metric"
  type         = "custom.googleapis.com/sample_metric"
  metric_kind  = "GAUGE"
  value_type   = "INT64"

  labels {
    key         = "sample_key"
    value_type  = "STRING"
    description = "Sample label key"
  }
}
```
