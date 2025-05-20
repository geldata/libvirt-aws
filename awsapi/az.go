package awsapi

import (
	"encoding/xml"
	"net/http"

	"github.com/google/uuid"
)

type DescribeAvailabilityZonesResponse struct {
	XMLName              xml.Name            `xml:"DescribeAvailabilityZonesResponse"`
	Xmlns                string              `xml:"xmlns,attr"`
	RequestID            string              `xml:"requestId"`
	AvailabilityZoneInfo []*AvailabilityZone `xml:"availabilityZoneInfo>item"`
}

type AvailabilityZone struct {
	GroupName          string `xml:"groupName"`
	OptInStatus        string `xml:"optInStatus"`
	ZoneName           string `xml:"zoneName"`
	ZoneID             string `xml:"zoneId"`
	ZoneState          string `xml:"zoneState"`
	ZoneType           string `xml:"zoneType"`
	RegionName         string `xml:"regionName"`
	MessageSet         string `xml:"messageSet"`
	NetworkBorderGroup string `xml:"NetworkBorderGroup"`
	GroupLongName      string `xml:"GroupLongName"`
}

func (a *AWSAPI) DescribeAvailabilityZones(w http.ResponseWriter, r *http.Request) {
	azs := []*AvailabilityZone{}
	for _, suffix := range []string{"a", "b", "c"} {
		az := &AvailabilityZone{
			OptInStatus: "opt-in-not-required",
			ZoneName:    a.Region + suffix,
			ZoneID:      a.Region + suffix,
			ZoneState:   "available",
			RegionName:  a.Region,
		}
		azs = append(azs, az)
	}

	response := &DescribeAvailabilityZonesResponse{
		Xmlns:                xmlns,
		RequestID:            uuid.New().String(),
		AvailabilityZoneInfo: azs,
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}
