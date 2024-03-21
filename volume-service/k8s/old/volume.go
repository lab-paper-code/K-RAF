package k8s

/*
import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	resourcev1 "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pvcSuffix         string = "-pvc"
	volumeNamespace   string = "vd"
	storageClassName  string = "rook-cephfs"
	volumeSizeDefault string = "20Gi"

	k8sTimeout time.Duration = 30 * time.Second
)


---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pod1-pvc
  namespace: ksv
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 20Gi
  storageClassName: rook-cephfs


// getPVCName makes pvc name
func (client *K8sClient) getPVCName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, pvcSuffix)
}

func (client *K8sClient) getPVCLabels(username string, volumeID string) map[string]string {
	return map[string]string{
		"username":  username,
		"volume-id": volumeID,
	}
}

func (client *K8sClient) getVolumeNamespace() string {
	return volumeNamespace
}

func (client *K8sClient) getStorageClassName() string {
	return storageClassName
}

func (client *K8sClient) getDefaultVolumeSize() resourcev1.Quantity {
	quantity, _ := resourcev1.ParseQuantity(volumeSizeDefault)
	return quantity
}

// CreatePVC creates a pvc for the given volumeID
func (client *K8sClient) CreatePVC(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreatePVC",
	})

	logger.Debugf("Creating a PVC for user %s, volume id %s", username, volumeID)

	scName := client.getStorageClassName()

	claim := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      client.getPVCName(volumeID),
			Labels:    client.getPVCLabels(username, volumeID),
			Namespace: client.getVolumeNamespace(),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			StorageClassName: &scName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: client.getDefaultVolumeSize(),
				},
			},
		},
	}

	pvcclient := client.clientSet.CoreV1().PersistentVolumeClaims(client.getVolumeNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := pvcclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		_, err = pvcclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err, "\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a PVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a PVC for user %s, volume id %s", username, volumeID)
	} else {
		_, err = pvcclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a PVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a PVC for user %s, volume id %s", username, volumeID)
	}

	return nil
}
*/
