package logic

import (
	"github.com/lab-paper-code/ksv/volume-service/commons"
	"github.com/lab-paper-code/ksv/volume-service/db"
	"github.com/lab-paper-code/ksv/volume-service/k8s"
)

type Logic struct {
	config *commons.Config

	dbAdapter  *db.DBAdapter
	k8sAdapter *k8s.K8SAdapter
}

// Start starts Logic
func Start(config *commons.Config, dbAdapter *db.DBAdapter, k8sAdapter *k8s.K8SAdapter) (*Logic, error) {
	logic := &Logic{
		config:     config,
		dbAdapter:  dbAdapter,
		k8sAdapter: k8sAdapter,
	}

	return logic, nil
}

// Stop stops Logic
func (logic *Logic) Stop() error {
	return nil
}
