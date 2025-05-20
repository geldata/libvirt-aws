//go:build integration
// +build integration

package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/assert"
)

var (
	portMap = make(map[int]bool)
)

func getPort() int {
	for i := 5100; i < 5200; i++ {
		if _, ok := portMap[i]; !ok {
			portMap[i] = true
			return i
		}
	}
	return 5100
}

func TestDescribeAvailabilityZonesIntegration(t *testing.T) {
	region := "us-east-2"
	port := getPort()
	addr := fmt.Sprintf("http://localhost:%d/", port)
	opts := &AWSEmulatorServerOpts{
		Port:   port,
		DBOpts: testDBOpts(),
	}
	server, err := NewAWSEmulatorServer(opts)
	assert.NoError(t, err)
	go server.Start()

	// make a request to healthz to wait for the server to start
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Get(addr + "healthz")
	if err != nil {
		t.Fatalf("failed to get healthz: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	client := ec2.NewFromConfig(aws.Config{}, func(o *ec2.Options) {
		o.BaseEndpoint = aws.String(addr)
	})

	t.Run("test describe azs", func(t *testing.T) {
		got, err := client.DescribeAvailabilityZones(context.Background(), &ec2.DescribeAvailabilityZonesInput{})
		assert.NoError(t, err)
		assert.Equal(t, len(got.AvailabilityZones), 3)
		azMap := make(map[int]string)
		azMap[0] = "a"
		azMap[1] = "b"
		azMap[2] = "c"
		for i, az := range got.AvailabilityZones {
			assert.Equal(t, *az.ZoneName, region+azMap[i])
			assert.Equal(t, *az.ZoneId, region+azMap[i])
			assert.Equal(t, *az.RegionName, region)
		}
	})

	server.Stop()
}

func TestDescribeAvailabilityZonesCustomRegionIntegration(t *testing.T) {
	region := "us-west-2"
	port := getPort()
	addr := fmt.Sprintf("http://localhost:%d/", port)
	opts := &AWSEmulatorServerOpts{
		Port:   port,
		Region: region,
		DBOpts: testDBOpts(),
	}
	server, err := NewAWSEmulatorServer(opts)
	assert.NoError(t, err)
	go server.Start()

	// make a request to healthz to wait for the server to start
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Get(addr + "healthz")
	if err != nil {
		t.Fatalf("failed to get healthz: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	client := ec2.NewFromConfig(aws.Config{}, func(o *ec2.Options) {
		o.BaseEndpoint = aws.String(addr)
	})

	t.Run("test describe azs custom region", func(t *testing.T) {
		got, err := client.DescribeAvailabilityZones(context.Background(), &ec2.DescribeAvailabilityZonesInput{})
		assert.NoError(t, err)
		assert.Equal(t, len(got.AvailabilityZones), 3)
		azMap := make(map[int]string)
		azMap[0] = "a"
		azMap[1] = "b"
		azMap[2] = "c"
		for i, az := range got.AvailabilityZones {
			assert.Equal(t, *az.ZoneName, region+azMap[i])
			assert.Equal(t, *az.ZoneId, region+azMap[i])
			assert.Equal(t, *az.RegionName, region)
		}
	})

	server.Stop()
}
