package k8s

import (
	"time"

	"github.com/lab-paper-code/ksv/volume-service/commons"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	operationTimeout time.Duration = 30 * time.Second
	objectNamespace  string        = "ksv"
)

type K8SAdapter struct {
	config     *commons.Config
	kubeConfig *rest.Config
	clientSet  *kubernetes.Clientset
}

// Start starts K8SAdapter
func Start(config *commons.Config) (*K8SAdapter, error) {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"function": "Start",
	})

	kubeConfigPath, err := commons.ExpandHomeDir(config.KubeConfigPath)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	service := &K8SAdapter{
		config:     config,
		kubeConfig: kubeConfig,
		clientSet:  clientSet,
	}

	return service, nil
}

// Stop stops K8SAdapter
func (adapter *K8SAdapter) Stop() error {
	return nil
}
