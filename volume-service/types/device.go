package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/xerrors"
)

const (
	deviceIDPrefix string = "dev"
)

// Device represents a device, holding all necessary info. about device
type Device struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	IP          string    `json:"ip"`
	Password    string    `json:"password"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func ValidateDeviceID(id string) error {
	if len(id) == 0 {
		return xerrors.Errorf("empty device id")
	}

	prefix := fmt.Sprintf("%s_", deviceIDPrefix)

	if !strings.HasPrefix(id, prefix) {
		return xerrors.Errorf("invalid device id - %s", id)
	}
	return nil
}

// NewDeviceID creates a new Device ID
func NewDeviceID() string {
	return fmt.Sprintf("%s_%s", deviceIDPrefix, xid.New().String())
}

func (device *Device) CheckAuthKey(authKey string) bool {
	expectedAuthKey := GetAuthKey(device.ID, device.Password)
	return expectedAuthKey == authKey
}

func (device *Device) GetRedacted() Device {
	dev := Device{}

	// copy
	dev = *device
	dev.Password = "<REDACTED>"
	return dev
}
