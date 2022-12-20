# Nagios CloudWatch Plugin

Nagios plugin to check AWS CloudWatch metrics using GetMetricData API.

## Installation

```console
$ go get github.com/quickguard-oss/nagios-cloudwatch-plugin
```

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
-q, --queries string      Array of MetricDataQuery objects in JSON format. See AWS GetMetricData API reference.
-w, --warning string      Metric range to result in warning status.
-c, --critical string     Metric range to result in critical status.
-p, --datapoints string   'n/m'; resulting to anomaly status if n of m datapoints are out of warning or critical range. (default "1/1")
-d, --duration int        Duration in minutes for which to retrieve metrics. (default 60)
-t, --timeout int         Seconds before plugin times out. (default 10)
-C, --classic-output      Print status message in classic format.
-v, --verbose count       Extra information. Up to 3 verbosity levels.
-V, --version             Print version information.
-h, --help                Print detailed help screen.
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

This plugin outputs a status line in JSON format by default.

```console
$ check_cloudwatch -q "..." -w '-5.0:5.0' -c '-10.0:10.0' -p 3/5 -d 6
{"level":"info","service":"CLOUDWATCH","status":"OK","time":"2022-12-13T16:01:59+09:00","message":"BurstUsage = 0.052259259259301416 | value=0.052259259259301416;-5.0:5.0;-10.0:10.0;;"}
```

Use the `-C` option if you want to display a status line in the [classic Nagios style](http://nagios-plugins.org/doc/guidelines.html#AEN33).

## Missing data

When evaluating alerts, missing datapoints are ignored and only existing ones are used in the evaluation.

## License

MIT
