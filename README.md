# Nagios CloudWatch Plugin

The Nagios CloudWatch Plugin enables Nagios to monitor various AWS services such as EC2, RDS, and S3, as well as custom metrics created in CloudWatch.

## Installation

```console
$ go install github.com/quickguard-oss/nagios-cloudwatch-plugin/cmd/check_cloudwatch@latest
```

or download a pre-built binary on our [releases page](https://github.com/quickguard-oss/nagios-cloudwatch-plugin/releases).

## AWS Credentials

This plugin uses [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2).

To specify AWS credentials, see [the official documentation](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials).

## Usage

```console
$ check_cloudwatch -q <queries> -w <range> -c <range> -p <datapoints>
                   [-d <duration>] [-t <timeout>] [-C] [-v]
```

Options:

```
  -q, --queries JSON     An array of MetricDataQuery objects in JSON format.
                         See the AWS GetMetricData API reference for details.
  -w, --warning range    Set the warning range for the metric.
  -c, --critical range   Set the critical range for the metric.
  -p, --datapoints n/m   Set the number of data points 'm' and the threshold 'n' for determining
                         a monitoring status. If 'n' or more of the 'm' data points are in the warning
                         or critical range, the status will be considered unhealthy. Should be
                         specified in the format 'n/m'.
                          (default "1/1")
  -d, --duration int     Set the duration in minutes for which to retrieve metrics.
                          (default 60)
  -t, --timeout int      Set the time in seconds before the plugin times out.
                          (default 10)
  -C, --classic-output   Print status message in classic format.
  -v, --verbose count    Enable extra information, with up to 3 verbosity levels.
  -V, --version          Print version information.
  -h, --help             Print detailed help information.
```

See [Nagios guidelines](http://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT) for the format of warning/critical ranges.

Example:

```console
$ cat ./queries.json
[
  {
    "Id": "m1",
    "MetricStat": {
      "Metric": {
        "Namespace": "AWS/EBS",
        "MetricName": "BurstBalance",
        "Dimensions": [
          {
            "Name": "VolumeId",
            "Value": "YOUR_VOLUME_ID"
          }
        ]
      },
      "Period": 60,
      "Stat": "Average"
    },
    "ReturnData": false
  },
  {
    "Id": "e1",
    "Label": "BurstUsage",
    "Expression": "DIFF(m1)"
  }
]

$ check_cloudwatch -q "$(< ./queries.json)" -w '-5.0:5.0' -c '-10.0:10.0' -p 3/5 -d 6 -C
CLOUDWATCH OK: BurstUsage = 0.052259259259301416 | value=0.052259259259301416;-5.0:5.0;-10.0:10.0;;
```

## Queries

The query format is an array of [MetricDataQuery](https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_MetricDataQuery.html).

The first metric in the returned set is used for alerting.

## Output

By default, this plugin outputs a status line in JSON format.

```console
$ check_cloudwatch -q "..." -w '-5.0:5.0' -c '-10.0:10.0' -p 3/5 -d 6
{"level":"info","service":"CLOUDWATCH","status":"OK","time":"2022-12-13T16:01:59+09:00","message":"BurstUsage = 0.052259259259301416 | value=0.052259259259301416;-5.0:5.0;-10.0:10.0;;"}
```

To display a status line in the [classic Nagios style](http://nagios-plugins.org/doc/guidelines.html#AEN33), use the `-C` flag.

## Missing data

When there are missing data points in the retrieved metrics, only the existing data points are used to determine the monitoring status, and the missing data points are ignored.

If the number of data points obtained is less than the value specified by the `-p` flag, the monitoring status results to `UNKNOWN`.

To ensure stable monitoring of metrics that may have missing data points, it is recommended to use the [`FILL` or `TIME_SERIES` function](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/using-metric-math.html#metric-math-syntax-functions-list) in conjunction with the target metrics.


```json
[
  {
    "Id": "m1",
    "MetricStat": {
      "Metric": {
        "Namespace": "AWS/ApiGateway",
        "MetricName": "IntegrationLatency",
        "Dimensions": [
          {
            "Name": "ApiName",
            "Value": "YOUR_API_NAME"
          }
        ]
      },
      "Period": 300,
      "Stat": "Average"
    },
    "ReturnData": false
  },
  {
    "Id": "e1",
    "Label": "IntegrationLatency",
    "Expression": "MAX([FILL(m1, LINEAR), TIME_SERIES(0)])"
  }
]
```


## License

MIT
