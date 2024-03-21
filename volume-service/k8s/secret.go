package k8s

import (
	"context"
	"fmt"

	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	secretNamePrefix string = "secret"
	secretNamespace  string = objectNamespace
)

func (adapter *K8SAdapter) GetSecretName(device *types.Device) string {
	return makeValidObjectName(secretNamePrefix, device.ID)
}

func (adapter *K8SAdapter) getSecretLabels(device *types.Device) map[string]string {
	labels := map[string]string{}
	labels["device-id"] = device.ID
	return labels
}

func (adapter *K8SAdapter) CreateSecret(device *types.Device) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreateSecret",
	})

	logger.Debug("received CreateSecret()")

	secretName := adapter.GetSecretName(device)
	secretLabels := adapter.getSecretLabels(device)

	// create secret
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Labels:    secretLabels,
			Namespace: secretNamespace,
		},
		Data: map[string][]byte{
			"username": []byte(device.ID),
			"password": []byte(device.Password),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	secretclient := adapter.clientSet.CoreV1().Secrets(secretNamespace)
	_, err := secretclient.Get(ctx, secret.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist
		_, createErr := secretclient.Create(ctx, secret, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := secretclient.Update(ctx, secret, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	fmt.Println("Secret created")

	return nil

}
