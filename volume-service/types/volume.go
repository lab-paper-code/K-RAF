package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/xerrors"
)

const (
	volumeIDPrefix string = "vol"
)

// Volume represents a volume, holding all necessary info. about volume
type Volume struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	DeviceID   string    `json:"device_id"`
	VolumeSize int64     `json:"volume_size"` // in bytes
	Mounted    bool      `json:"mounted"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func ValidateVolumeID(id string) error {
	if len(id) == 0 {
		return xerrors.Errorf("empty volume id")
	}

	prefix := fmt.Sprintf("%s_", volumeIDPrefix)

	if !strings.HasPrefix(id, prefix) {
		return xerrors.Errorf("invalid volume id - %s", id)
	}
	return nil
}

// NewVolumeID creates a new Volume ID
func NewVolumeID() string {
	return fmt.Sprintf("%s_%s", volumeIDPrefix, xid.New().String())
}
