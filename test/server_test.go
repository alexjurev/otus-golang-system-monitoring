package test

import (
	"context"
	"github.com/alexjurev/otus-golang-system-monitoring/api"
	"github.com/alexjurev/otus-golang-system-monitoring/internal/logger"
	metricloader "github.com/alexjurev/otus-golang-system-monitoring/internal/metrics/loader"
	internalgrpc "github.com/alexjurev/otus-golang-system-monitoring/internal/server/grpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	serverHost = "127.0.0.1"
	serverPort = 9006
)

func TestMain(m *testing.M) {
	_ = logger.PrepareLogger(logger.Config{Level: "ERROR"})

	code := m.Run()
	os.Exit(code)
}

func TestServerIncorrectArgument(t *testing.T) {
	startServer(t)
	type testCase struct {
		name string
		req  *api.GetMetricsRequest
	}

	testCases := []testCase{
		{
			name: "NotifyInterval: 0, AverageCalcInterval: 0",
			req:  &api.GetMetricsRequest{NotifyInterval: 0, AverageCalcInterval: 0},
		},
		{
			name: "NotifyInterval: 1, AverageCalcInterval: -1",
			req:  &api.GetMetricsRequest{NotifyInterval: 1, AverageCalcInterval: -1},
		},
		{
			name: "NotifyInterval: -1, AverageCalcInterval: 1",
			req:  &api.GetMetricsRequest{NotifyInterval: -1, AverageCalcInterval: 1},
		},
		{
			name: "NotifyInterval: -1, AverageCalcInterval: -1",
			req:  &api.GetMetricsRequest{NotifyInterval: -1, AverageCalcInterval: -1},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := createClient(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
			defer cancel()

			r, err := client.GetMetrics(ctx, tc.req)
			require.NoError(t, err)

			_, err = r.Recv()
			require.Error(t, err)
			require.Equal(t, codes.InvalidArgument, status.Convert(err).Code())
		})
	}
}

func startServer(t *testing.T) {
	t.Helper()

	c, _ := metricloader.Load(metricloader.Config{IgnoreUnavailable: true})
	grpcServer := internalgrpc.NewServer(internalgrpc.Config{
		Host: serverHost,
		Port: serverPort,
	}, c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = grpcServer.Start()
	}()

	// Wait stating of servers
	require.Eventually(t, func() bool {
		conn, err := grpc.Dial(
			net.JoinHostPort(serverHost, strconv.Itoa(serverPort)),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return false
		}
		client := api.NewMetricsClient(conn)

		m, err := client.GetMetrics(ctx, &api.GetMetricsRequest{NotifyInterval: -1, AverageCalcInterval: -1})
		require.NoError(t, err)
		_, err = m.Recv()

		if err != nil {
			s := status.Convert(err)
			if s != nil && s.Code() == codes.InvalidArgument {
				return true
			}
		}
		return false
	}, 5*time.Second, 200*time.Millisecond)

	t.Cleanup(func() {
		cancel()
		grpcServer.Stop()
	})
}

func createClient(t *testing.T) api.MetricsClient {
	t.Helper()
	conn, err := grpc.Dial(
		net.JoinHostPort(serverHost, strconv.Itoa(serverPort)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to create client")
	return api.NewMetricsClient(conn)
}
