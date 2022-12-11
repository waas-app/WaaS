package model

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	Owner             string     `json:"owner" gorm:"type:varchar(100);unique_index:key;primary_key"`
	OwnerName         string     `json:"owner_name"`
	OwnerEmail        string     `json:"owner_email"`
	Name              string     `json:"name" gorm:"type:varchar(100);unique_index:key;primary_key"`
	PublicKey         string     `json:"public_key" gorm:"unique_index"`
	Address           string     `json:"address"`
	CreatedAt         time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"column:updated_at"`
	LastHandshakeTime *time.Time `json:"last_handshake_time"`
	ReceiveBytes      int64      `json:"received_bytes"`
	TransmitBytes     int64      `json:"transmit_bytes"`
	Endpoint          string     `json:"endpoint"`
}

func (d *Device) TableName() string {
	return "devices"
}

func (d *Device) AfterSave(tx *gorm.DB) error {
	return nil
}

func (d *Device) AfterDelete(tx *gorm.DB) error {
	return nil
}
