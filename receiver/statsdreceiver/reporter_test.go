// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package statsdreceiver

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/obsreport/obsreporttest"
)

func TestReporterObservability(t *testing.T) {
	receiverID := component.NewIDWithName(typeStr, "fake_receiver")
	tt, err := obsreporttest.SetupTelemetry(receiverID)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tt.Shutdown(context.Background()))
	}()

	reporter, err := newReporter(tt.ToReceiverCreateSettings())
	require.NoError(t, err)

	ctx := reporter.OnDataReceived(context.Background())

	reporter.OnMetricsProcessed(ctx, 17, nil)

	require.NoError(t, tt.CheckReceiverMetrics("tcp", 17, 0))

	// Below just exercise the error paths.
	err = errors.New("fake error for tests")
	reporter.OnTranslationError(ctx, err)
	reporter.OnMetricsProcessed(ctx, 10, err)

	require.NoError(t, tt.CheckReceiverMetrics("tcp", 17, 10))
}
