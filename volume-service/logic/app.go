package logic

import (
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
)

func (logic *Logic) ListApps() ([]types.App, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListApps",
	})

	logger.Debug("received ListApps()")

	return logic.dbAdapter.ListApps()
}

func (logic *Logic) GetApp(appID string) (types.App, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetApp",
	})

	logger.Debug("received GetApp()")

	return logic.dbAdapter.GetApp(appID)
}

func (logic *Logic) CreateApp(app *types.App) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateApp",
	})

	logger.Debug("received CreateApp()")

	return logic.dbAdapter.InsertApp(app)
}

func (logic *Logic) ListAppRuns(deviceID string) ([]types.AppRun, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListAppRuns",
	})

	logger.Debug("received ListAppRuns()")

	return logic.dbAdapter.ListAppRuns(deviceID)
}

func (logic *Logic) ListAllAppRuns() ([]types.AppRun, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListAllAppRuns",
	})

	logger.Debug("received ListAllAppRuns()")

	return logic.dbAdapter.ListAllAppRuns()
}

func (logic *Logic) GetAppRun(appRunID string) (types.AppRun, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetAppRun",
	})

	logger.Debug("received GetAppRun()")

	return logic.dbAdapter.GetAppRun(appRunID)
}

func (logic *Logic) ExecuteApp(appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ExecuteApp",
	})

	logger.Debug("received ExecuteApp()")

	app, err := logic.GetApp(appRun.AppID)
	if err != nil {
		return err
	}

	device, err := logic.dbAdapter.GetDevice(appRun.DeviceID)
	if err != nil {
		return err
	}

	volume, err := logic.dbAdapter.GetVolume(appRun.VolumeID)
	if err != nil {
		return err
	}

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.CreateApp()")
	} else {
		logger.Debugf("creating App %s for device %s, volume %s", app.Name, device.ID, volume.ID)
		err = logic.k8sAdapter.CreateApp(&device, &volume, &app, appRun)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.InsertAppRun(appRun)
}

func (logic *Logic) TerminateAppRun(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "TerminateAppRun",
	})

	logger.Debug("received TerminateAppRun()")

	appRun, err := logic.dbAdapter.GetAppRun(appRunID)
	if err != nil {
		return err
	}

	device, err := logic.dbAdapter.GetDevice(appRun.DeviceID)
	if err != nil {
		return err
	}

	volume, err := logic.dbAdapter.GetVolume(appRun.VolumeID)
	if err != nil {
		return err
	}

	if logic.config.NoKubernetes {
		logger.Debug("bypass k8sAdapter.DeleteApp()")
	} else {
		logger.Debugf("stopping App Run %s for device %s, volume %s", appRun.ID, device.ID, volume.ID)
		err = logic.k8sAdapter.DeleteApp(appRunID)
		if err != nil {
			return err
		}
	}

	return logic.dbAdapter.UpdateAppRunTermination(appRunID, true)
}
