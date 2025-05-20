package awsapi

import "gorm.io/gorm"

const (
	xmlns         = "http://ec2.amazonaws.com/doc/2016-11-15/"
	DefaultRegion = "us-east-2"
)

type AWSAPI struct {
	db     *gorm.DB
	Region string
}

func NewAWSAPI(db *gorm.DB, region string) *AWSAPI {
	if region == "" {
		region = DefaultRegion
	}
	return &AWSAPI{
		db:     db,
		Region: region,
	}
}
