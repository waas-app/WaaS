package datastore

import (
	"context"

	"github.com/hjoshi123/WaaS/infra/database"
	"github.com/hjoshi123/WaaS/model"
	"github.com/hjoshi123/WaaS/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type DeviceStore interface {
	Save(ctx context.Context, device *model.Device) error
	List(ctx context.Context, owner string) ([]*model.Device, error)
	Get(ctx context.Context, owner string, name string) (*model.Device, error)
	GetByPublicKey(ctx context.Context, publicKey string) (*model.Device, error)
	Delete(ctx context.Context, device *model.Device) error
}

type deviceStore struct{}

func NewDeviceStore() DeviceStore {
	return &deviceStore{}
}

func (s *deviceStore) Save(ctx context.Context, device *model.Device) error {
	db := database.Instance(ctx)
	if err := db.Save(device).Error; err != nil {
		util.Logger(ctx).Error("Failed to save device", zap.Error(err))
		return err
	}
	return nil
}

func (s *deviceStore) List(ctx context.Context, owner string) ([]*model.Device, error) {
	db := database.Instance(ctx)
	devices := make([]*model.Device, 0)
	var err error
	if owner != "" {
		err = db.Where("owner = ?", owner).Find(&devices).Error
	} else {
		err = db.Find(&devices).Error
	}

	if err != nil {
		util.Logger(ctx).Error("Failed to list devices", zap.Error(err))
		return nil, err
	}

	util.Logger(ctx).Debug("List devices", zap.Int("count", len(devices)))
	return devices, nil
}

func (s *deviceStore) Get(ctx context.Context, owner string, name string) (*model.Device, error) {
	db := database.Instance(ctx)
	device := new(model.Device)
	if err := db.Where("owner = ? AND name = ?", owner, name).First(&device).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to read device")
	}
	return device, nil
}

func (s *deviceStore) GetByPublicKey(ctx context.Context, publicKey string) (*model.Device, error) {
	db := database.Instance(ctx)
	device := new(model.Device)
	if err := db.Where("public_key = ?", publicKey).First(&device).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to read device")
	}

	return device, nil
}

func (s *deviceStore) Delete(ctx context.Context, device *model.Device) error {
	db := database.Instance(ctx)
	if err := db.Delete(device).Error; err != nil {
		util.Logger(ctx).Error("Failed to delete device", zap.Error(err))
		return err
	}

	return nil
}
