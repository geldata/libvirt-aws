package models

type Tag struct {
	ResourceName string `gorm:"uniqueIndex:rname_rtype_tname"`
	ResourceType string `gorm:"uniqueIndex:rname_rtype_tname"`
	TagName      string `gorm:"uniqueIndex:rname_rtype_tname"`
	TagValue     string
}

type IPAddress struct {
	AllocationID     string `gorm:"uniqueIndex:allocation_id"`
	IPAddress        string `gorm:"uniqueIndex:ip_address"`
	AssociationID    string `gorm:"uniqueIndex:association_id"`
	InstanceID       string
	PrivateIPAddress string
}

type PrivateIPAddress struct {
	IPAddress  string `gorm:"primaryKey:ip_address"`
	InstanceID string
	Interface  string
}

type DNSZone struct {
	ID      string `gorm:"primaryKey:id"`
	Name    string
	Comment string
}

type DNSChange struct {
	ID          string `gorm:"primaryKey:id"`
	SubmittedAt string
	Comment     string
}

type VolumeModification struct {
	ID            string `gorm:"primaryKey:id"`
	Modifications string
}

func All() []interface{} {
	return []interface{}{
		&Tag{},
		&IPAddress{},
		&PrivateIPAddress{},
		&DNSZone{},
		&DNSChange{},
		&VolumeModification{},
	}
}
