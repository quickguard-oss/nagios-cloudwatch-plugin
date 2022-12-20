package cloudwatch

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client/types"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/container"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/log"
)

type CloudWatch struct {
	client   types.Client
	duration int
	queries  []awstypes.MetricDataQuery
	timeout  int
	result   *cloudwatch.GetMetricDataOutput
}

func New(duration int, queriesStr string, timeout int) (CloudWatch, error) {
	client, err := container.GetCloudWatchClient()

	if err != nil {
		return CloudWatch{}, err
	}

	queries, err := parseQueries(queriesStr)

	if err != nil {
		return CloudWatch{}, err
	}

	return CloudWatch{
		client:   client,
		duration: duration,
		queries:  queries,
		timeout:  timeout,
	}, nil
}

func parseQueries(queries string) ([]awstypes.MetricDataQuery, error) {
	var q []awstypes.MetricDataQuery

	log.V(3).Trace().
		Str("package", "cloudwatch").
		Msg("parsing API queries")

	if err := json.Unmarshal([]byte(queries), &q); err != nil {
		return nil, errors.NewArgumentErrorWithError(err, "queries", queries)
	}

	log.V(3).Trace().
		Str("package", "cloudwatch").
		RawJSON("queries", []byte(queries)).
		Send()

	return q, nil
}

func (c *CloudWatch) GetMetricValues(now time.Time) ([]float64, error) {
	log.V(3).Trace().
		Str("package", "cloudwatch").
		Msg("calling GetMetricData API")

	if err := c.getMetricData(now); err != nil {
		return []float64{}, err
	}

	c.printResult()

	return c.result.MetricDataResults[0].Values, nil
}

func (c CloudWatch) LatestValue() (metricName string, value float64, timestamp time.Time) {
	m := c.result.MetricDataResults[0]

	if m.Label == nil {
		metricName = *m.Id
	} else {
		metricName = *m.Label
	}

	if len(m.Values) != 0 {
		value = m.Values[0]
	}

	if len(m.Timestamps) != 0 {
		timestamp = m.Timestamps[0]
	}

	return
}

func (c *CloudWatch) getMetricData(now time.Time) error {
	startTime := now.Add(-1 * time.Duration(c.duration) * time.Minute)

	ctx, cancel := context.WithDeadline(
		context.Background(),
		now.Add(time.Duration(c.timeout)*time.Second),
	)

	defer cancel()

	log.V(3).Trace().
		Str("package", "cloudwatch").
		Time("start_time", startTime).
		Time("end_time", now).
		Int("timeout", c.timeout).
		Msg("API parameters")

	result, err := c.client.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
		StartTime:         aws.Time(startTime),
		EndTime:           aws.Time(now),
		MetricDataQueries: c.queries,
	})

	if err != nil {
		return errors.NewCloudWatchError(err)
	}

	c.result = result

	return nil
}

func (c CloudWatch) printResult() {
	r, err := json.Marshal(c.result)

	if err != nil {
		log.V(2).Error().
			Err(err).
			Msg("failed to marshal GetMetricDataOutput")
	} else {
		log.V(2).Debug().
			RawJSON("GetMetricDataOutput", r).
			Msg("API call succeeds")
	}
}
