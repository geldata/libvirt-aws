package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/geldata/libvirt-aws/awsapi"
	"github.com/geldata/libvirt-aws/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getTestDB() *gorm.DB {
	db, err := db.NewDB(&db.DBOpts{
		DBFile: ":memory:",
		Config: &gorm.Config{},
	})
	if err != nil {
		log.Fatalf("failed to open test database: %v", err)
	}
	return db
}

func TestNewAWSEmulatorServer(t *testing.T) {
	tests := []struct {
		name string
		opts *AWSEmulatorServerOpts
		want *AWSEmulatorServer
	}{
		{
			name: "default options",
			opts: &AWSEmulatorServerOpts{},
			want: &AWSEmulatorServer{
				BindTo: "",
				Port:   5100,
				Addr:   ":5100",
				Debug:  false,
				Region: "us-east-2",
			},
		},
		{
			name: "custom options",
			opts: &AWSEmulatorServerOpts{
				BindTo: "169.254.169.254",
				Port:   9090,
				Debug:  true,
				Region: "us-west-2",
			},
			want: &AWSEmulatorServer{
				BindTo: "169.254.169.254",
				Port:   9090,
				Addr:   "169.254.169.254:9090",
				Debug:  true,
				Region: "us-west-2",
			},
		},
	}
	testDB := getTestDB()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAWSEmulatorServer(testDB, tt.opts)
			assert.Equal(t, tt.want.BindTo, got.BindTo)
			assert.Equal(t, tt.want.Port, got.Port)
			assert.Equal(t, tt.want.Addr, got.Addr)
			assert.Equal(t, tt.want.Debug, got.Debug)
			assert.Equal(t, tt.want.Region, got.Region)
		})
	}
}

func TestDescribeAvailabilityZones(t *testing.T) {
	tests := []struct {
		name       string
		region     string
		wantRegion string
	}{
		{
			name:       "canada central 1",
			region:     "ca-central-1",
			wantRegion: "ca-central-1",
		},
		{
			name:       "no region",
			wantRegion: "us-east-2",
		},
	}
	testDB := getTestDB()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Set("Action", "DescribeAvailabilityZones")
			form.Set("Version", "2016-11-15")
			w := httptest.NewRecorder()
			s := NewAWSEmulatorServer(testDB, &AWSEmulatorServerOpts{
				Region: tt.region,
			})
			awsapi.DescribeAvailabilityZones(s.Region, w, nil)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "<DescribeAvailabilityZonesResponse")
			assert.Contains(t, w.Body.String(), fmt.Sprintf("<regionName>%s</regionName>", tt.wantRegion))
			for _, suffix := range []string{"a", "b", "c"} {
				assert.Contains(t, w.Body.String(), fmt.Sprintf("<zoneName>%s%s</zoneName>", tt.wantRegion, suffix))
				assert.Contains(t, w.Body.String(), fmt.Sprintf("<zoneId>%s%s</zoneId>", tt.wantRegion, suffix))
			}
		})
	}
}
