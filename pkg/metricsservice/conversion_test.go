/*
Copyright 2026 The KEDA Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package metricsservice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/external_metrics"
	"k8s.io/metrics/pkg/apis/external_metrics/v1beta1"

	"github.com/kedacore/keda/v2/pkg/metricsservice/api"
)

func TestExternalMetricsProtoRoundTrip(t *testing.T) {
	window := int64(60)
	cases := []struct {
		name string
		in   *external_metrics.ExternalMetricValueList
	}{
		{
			name: "nil input returns empty list",
			in:   nil,
		},
		{
			name: "empty items",
			in:   &external_metrics.ExternalMetricValueList{},
		},
		{
			name: "single item, full fields",
			in: &external_metrics.ExternalMetricValueList{
				Items: []external_metrics.ExternalMetricValue{
					{
						MetricName:    "queue_depth",
						MetricLabels:  map[string]string{"queue": "orders", "region": "eu"},
						Timestamp:     metav1.NewTime(time.Date(2026, 5, 28, 12, 34, 56, 789000000, time.UTC)),
						WindowSeconds: &window,
						Value:         resource.MustParse("1500m"),
					},
				},
			},
		},
		{
			name: "multiple items, mixed quantity forms",
			in: &external_metrics.ExternalMetricValueList{
				Items: []external_metrics.ExternalMetricValue{
					{MetricName: "a", Value: resource.MustParse("0")},
					{MetricName: "b", Value: resource.MustParse("5")},
					{MetricName: "c", Value: resource.MustParse("100m")},
					{MetricName: "d", Value: resource.MustParse("3.14")},
					{MetricName: "e", Value: resource.MustParse("2Ki")},
				},
			},
		},
		{
			name: "nil WindowSeconds preserved",
			in: &external_metrics.ExternalMetricValueList{
				Items: []external_metrics.ExternalMetricValue{
					{MetricName: "no_window", Value: resource.MustParse("1"), WindowSeconds: nil},
				},
			},
		},
		{
			name: "nil MetricLabels preserved",
			in: &external_metrics.ExternalMetricValueList{
				Items: []external_metrics.ExternalMetricValue{
					{MetricName: "no_labels", Value: resource.MustParse("1"), MetricLabels: nil},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			proto := externalMetricsToProto(tc.in)
			got, err := protoToExternalMetrics(proto)
			require.NoError(t, err)

			if tc.in == nil {
				require.NotNil(t, got)
				require.Empty(t, got.Items)
				return
			}

			require.Equal(t, len(tc.in.Items), len(got.Items), "item count mismatch")
			for i := range tc.in.Items {
				want := tc.in.Items[i]
				have := got.Items[i]
				require.Equal(t, want.MetricName, have.MetricName, "MetricName at item %d", i)
				require.Equal(t, want.MetricLabels, have.MetricLabels, "MetricLabels at item %d", i)
				// Quantity equality via Cmp to handle canonical-form normalization.
				require.Truef(t, want.Value.Cmp(have.Value) == 0,
					"Value at item %d: want=%s have=%s", i, want.Value.String(), have.Value.String())
				if want.WindowSeconds == nil {
					require.Nil(t, have.WindowSeconds, "WindowSeconds at item %d", i)
				} else {
					require.NotNil(t, have.WindowSeconds, "WindowSeconds at item %d", i)
					require.Equal(t, *want.WindowSeconds, *have.WindowSeconds, "WindowSeconds at item %d", i)
				}
				require.Truef(t, want.Timestamp.Time.Equal(have.Timestamp.Time),
					"Timestamp at item %d: want=%s have=%s", i, want.Timestamp.Time, have.Timestamp.Time)
			}
		})
	}
}

// TestExternalMetricsProtoMalformedQuantity ensures a bad value on the wire
// surfaces as an error rather than silently producing a zero Quantity.
func TestExternalMetricsProtoMalformedQuantity(t *testing.T) {
	bad := externalMetricsToProto(&external_metrics.ExternalMetricValueList{
		Items: []external_metrics.ExternalMetricValue{
			{MetricName: "ok", Value: resource.MustParse("1")},
		},
	})
	bad.Items[0].Value = &api.Quantity{String_: "not-a-quantity"}

	_, err := protoToExternalMetrics(bad)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not-a-quantity")
}

// TestExternalMetricsProtoWireCompatWithUpstream checks that bytes produced by the upstream v1beta1.ExternalMetricValueList
// type decode cleanly through KEDA's own proto definition, and vice versa.
func TestExternalMetricsProtoWireCompatWithUpstream(t *testing.T) {
	window := int64(60)
	timestamp := metav1.NewTime(time.Date(2026, 5, 28, 12, 34, 56, 0, time.UTC))

	upstream := &v1beta1.ExternalMetricValueList{
		ListMeta: metav1.ListMeta{ResourceVersion: "12345"},
		Items: []v1beta1.ExternalMetricValue{
			{
				MetricName:    "queue_depth",
				MetricLabels:  map[string]string{"queue": "orders", "region": "eu"},
				Timestamp:     timestamp,
				WindowSeconds: &window,
				Value:         resource.MustParse("1500m"),
			},
			{
				MetricName: "no_window",
				Value:      resource.MustParse("5"),
			},
		},
	}

	t.Run("upstream bytes decode into KEDA schema", func(t *testing.T) {
		wireBytes, err := upstream.Marshal()
		require.NoError(t, err)

		decoded := &api.ExternalMetricValueList{}
		require.NoError(t, proto.Unmarshal(wireBytes, decoded))
		require.Len(t, decoded.Items, 2)

		require.Equal(t, "queue_depth", decoded.Items[0].MetricName)
		require.Equal(t, map[string]string{"queue": "orders", "region": "eu"}, decoded.Items[0].MetricLabels)
		require.NotNil(t, decoded.Items[0].Timestamp)
		require.True(t, decoded.Items[0].Timestamp.AsTime().Equal(timestamp.Time))
		require.NotNil(t, decoded.Items[0].Window)
		require.Equal(t, int64(60), *decoded.Items[0].Window)
		require.Equal(t, "1500m", decoded.Items[0].GetValue().GetString_())

		require.Equal(t, "no_window", decoded.Items[1].MetricName)
		require.Nil(t, decoded.Items[1].Window)
		require.Equal(t, "5", decoded.Items[1].GetValue().GetString_())
	})

	t.Run("KEDA bytes decode into upstream schema", func(t *testing.T) {
		kedaForm := externalMetricsToProto(&external_metrics.ExternalMetricValueList{
			Items: []external_metrics.ExternalMetricValue{
				{
					MetricName:    "queue_depth",
					MetricLabels:  map[string]string{"queue": "orders", "region": "eu"},
					Timestamp:     timestamp,
					WindowSeconds: &window,
					Value:         resource.MustParse("1500m"),
				},
				{MetricName: "no_window", Value: resource.MustParse("5")},
			},
		})

		wireBytes, err := proto.Marshal(kedaForm)
		require.NoError(t, err)

		decoded := &v1beta1.ExternalMetricValueList{}
		require.NoError(t, decoded.Unmarshal(wireBytes))
		require.Len(t, decoded.Items, 2)

		require.Equal(t, "queue_depth", decoded.Items[0].MetricName)
		require.Equal(t, map[string]string{"queue": "orders", "region": "eu"}, decoded.Items[0].MetricLabels)
		require.True(t, decoded.Items[0].Timestamp.Time.Equal(timestamp.Time))
		require.NotNil(t, decoded.Items[0].WindowSeconds)
		require.Equal(t, int64(60), *decoded.Items[0].WindowSeconds)
		require.True(t, decoded.Items[0].Value.Equal(resource.MustParse("1500m")))

		require.Equal(t, "no_window", decoded.Items[1].MetricName)
		require.Nil(t, decoded.Items[1].WindowSeconds)
		require.True(t, decoded.Items[1].Value.Equal(resource.MustParse("5")))
	})
}
