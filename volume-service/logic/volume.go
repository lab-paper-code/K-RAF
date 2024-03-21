package logic

import (
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
)

func (logic *Logic) ListVolumes(deviceID string) ([]types.Volume, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListVolumes",
	})

	logger.Debug("received ListVolumes()")

	return logic.dbAdapter.ListVolumes(deviceID)
}

func (logic *Logic) ListAllVolumes() ([]types.Volume, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListAllVolumes",
	})

	logger.Debug("received ListAllVolumes()")

	return logic.dbAdapter.ListAllVolumes()
}

func (logic *Logic) GetVolume(volumeID string) (types.Volume, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetVolume",
	})

	logger.Debug("received GetVolume()")

	return logic.dbAdapter.GetVolume(volumeID)
}

func (logic *Logic) CreateVolume(volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateVolume",
	})

	logger.Debug("received CreateVolume()")

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.CreateVolume()")
	} else {
		err := logic.k8sAdapter.CreateVolume(volume)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.InsertVolume(volume)
}

func (logic *Logic) ResizeVolume(volumeID string, size int64) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ResizeVolume",
	})

	logger.Debug("received ResizeVolume()")

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.ResizeVolume()")
	} else {
		err := logic.k8sAdapter.ResizeVolume(volumeID, size)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.UpdateVolumeSize(volumeID, size)
}

func (logic *Logic) DeleteVolume(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "DeleteVolume",
	})

	logger.Debug("received DeleteVolume()")

	volume, err := logic.dbAdapter.GetVolume(volumeID)
	if err != nil {
		return err
	}

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.DeleteVolume()")
	} else {
		// already mounted -> cannot delete
		if volume.Mounted {
			logger.Debugf("volume %s is already mounted, please unmount before you delete volume", volume.ID)
			return nil
		}

		logger.Debugf("Deleting volume %s for device %s", volume.ID, volume.DeviceID)
		err = logic.k8sAdapter.DeleteVolume(volumeID)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.DeleteVolume(&volume)
}

func (logic *Logic) MountVolume(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "MountVolume",
	})

	logger.Debug("received MountVolume()")

	volume, err := logic.dbAdapter.GetVolume(volumeID)
	if err != nil {
		return err
	}

	// already mounted
	if volume.Mounted {
		return nil
	}

	device, err := logic.dbAdapter.GetDevice(volume.DeviceID)
	if err != nil {
		return err
	}

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.CreateWebdav()")
	} else {
		logger.Debugf("creating Webdav for device %s, volume %s", device.ID, volume.ID)
		err = logic.k8sAdapter.CreateWebdav(&device, &volume)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.UpdateVolumeMount(volumeID, true)
}

func (logic *Logic) UnmountVolume(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UnmountVolume",
	})

	logger.Debug("received UnmountVolume()")

	volume, err := logic.dbAdapter.GetVolume(volumeID)
	if err != nil {
		return err
	}

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.CreateWebdav()")
	} else {
		// already unmounted
		if !volume.Mounted {
			logic.k8sAdapter.EnsureDeleteWebdav(volumeID)
			return nil
		}

		logger.Debugf("stopping Webdav for device %s, volume %s", volume.DeviceID, volume.ID)
		err = logic.k8sAdapter.DeleteWebdav(volumeID)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.UpdateVolumeMount(volumeID, false)
}
