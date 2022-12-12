package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/infra/red"
	"gorm.io/gorm"
)

type DevicePayload struct {
	Type   string  `json:"type"`
	Device *Device `json:"device"`
}

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

func (d *Device) AfterCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}

	payload := new(DevicePayload)
	payload.Type = config.DevicesCreate

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := red.Message{
		Topic:   payload.Type,
		Payload: string(p),
	}

	pubsubHandler, err := red.GetPubsubClientHandler()
	if err != nil {
		return err
	}

	err = pubsubHandler.Publish(ctx, msg)
	return err
}

func (d *Device) BeforeDelete(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}

	payload := new(DevicePayload)
	payload.Type = config.DevicesDelete

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := red.Message{
		Topic:   payload.Type,
		Payload: string(p),
	}

	pubsubHandler, err := red.GetPubsubClientHandler()
	if err != nil {
		return err
	}

	err = pubsubHandler.Publish(ctx, msg)
	return err
}
