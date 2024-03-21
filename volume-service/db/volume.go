package db

import (
	"github.com/lab-paper-code/ksv/volume-service/types"
	"golang.org/x/xerrors"
)

func (adapter *DBAdapter) ListVolumes(deviceID string) ([]types.Volume, error) {
	volumes := []types.Volume{}
	result := adapter.db.Where("device_id = ?", deviceID).Find(&volumes)
	if result.Error != nil {
		return nil, result.Error
	}

	return volumes, nil
}

func (adapter *DBAdapter) ListAllVolumes() ([]types.Volume, error) {
	volumes := []types.Volume{}
	result := adapter.db.Find(&volumes)
	if result.Error != nil {
		return nil, result.Error
	}

	return volumes, nil
}

func (adapter *DBAdapter) GetVolume(volumeID string) (types.Volume, error) {
	var volume types.Volume
	result := adapter.db.Where("id = ?", volumeID).First(&volume)
	if result.Error != nil {
		return volume, result.Error
	}

	return volume, nil
}

func (adapter *DBAdapter) InsertVolume(volume *types.Volume) error {
	result := adapter.db.Create(volume)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert a volume")
	}

	return nil
}

func (adapter *DBAdapter) DeleteVolume(volume *types.Volume) error {
	result := adapter.db.Delete(volume)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to delete a volume")
	}

	return nil
}

func (adapter *DBAdapter) UpdateVolumeSize(volumeID string, size int64) error {
	var record types.Volume
	result := adapter.db.Where("id = ?", volumeID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.VolumeSize = size

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateVolumeMount(volumeID string, mounted bool) error {
	var record types.Volume
	result := adapter.db.Where("id = ?", volumeID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Mounted = mounted

	adapter.db.Save(&record)

	return nil
}
