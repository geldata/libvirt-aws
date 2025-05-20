package awsapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAWSAPI(nil, tt.region)
			w := httptest.NewRecorder()
			api.DescribeAvailabilityZones(w, nil)
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
