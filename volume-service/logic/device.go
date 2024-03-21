package logic

import (
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
)

func (logic *Logic) ListDevices() ([]types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListDevices",
	})

	logger.Debug("received ListDevices()")

	return logic.dbAdapter.ListDevices()
}

func (logic *Logic) GetDevice(deviceID string) (types.Device, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetDevice",
	})

	logger.Debug("received GetDevice()")

	return logic.dbAdapter.GetDevice(deviceID)
}

func (logic *Logic) CreateDevice(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateDevice",
	})

	logger.Debug("received CreateDevice()")

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.CreateSecret()")
	} else {
		err := logic.k8sAdapter.CreateSecret(device)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.InsertDevice(device)
}

func (logic *Logic) UpdateDeviceIP(deviceID string, ip string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceIP",
	})

	logger.Debug("received UpdateDeviceIP()")

	return logic.dbAdapter.UpdateDeviceIP(deviceID, ip)
}

func (logic *Logic) UpdateDevicePassword(deviceID string, password string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDevicePassword",
	})

	logger.Debug("received UpdateDevicePassword()")

	return logic.dbAdapter.UpdateDevicePassword(deviceID, password)
}
