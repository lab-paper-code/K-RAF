package k8s

/*
import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	svcWebdavSuffix string = "-webdav-svc"
	svcAppSuffix    string = "-app-svc"
	svcNamespace    string = "vd"
)


apiVersion: v1
kind: Service
metadata:
  name: pod1-webdav-svc #변경
  namespace: ksv(webdav namespace) #변경
spec:
  ports:
  - port: 80
    protocol: TCP
  selector:
    app: pod1-webdav-svc #변경
---
apiVersion: v1
kind: Service
metadata:
  name: pod1-app-svc #변경
  namespace: ksv
spec:
  ports:
  - port: 60000
    protocol: TCP
  selector:
    app: pod1-app-svc #변경
  type: Clust


// getAppSvcName makes appSvc name
func (client *K8sClient) getAppSVCName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, svcAppSuffix)
}

func (client *K8sClient) getSvcNamespace() string {
	return svcNamespace
}

//Create SVC

func (client *K8sClient) CreateSVC(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreateSVC",
	})

	logger.Debugf("Creating a SVC for user %s, volume id %s", username, volumeID)

	claim := &corev1.ServiceList{
		Items: []corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      client.getWebdavSVCName(volumeID),
					Namespace: client.getSvcNamespace(),
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeClusterIP,
					Ports: []corev1.ServicePort{
						{
							Port:     int32(80),
							Protocol: corev1.ProtocolTCP,
						},
					},
					Selector: map[string]string{
						"app": client.getDeployWebdavName(volumeID),
					},
				},
			},


			///
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      client.getAppSVCName(volumeID),
					Namespace: client.getSvcNamespace(),
				},
				Spec: corev1.ServiceSpec{ // ServiceSpec describes the attributes that a user creates on a service.
					Ports: []corev1.ServicePort{
						{
							Port:     int32(60000),
							Protocol: corev1.ProtocolTCP,
						},
					},
					Selector: map[string]string{
						"app": client.getDeployAppName(volumeID),
					},
					Type: "ClusterIP",
				},
				 //identical code with CreateApp
			},
		},
	}

	SVCclient := client.clientSet.CoreV1().Services(client.getSvcNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := SVCclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		print("\nCREATE!!\n")
		_, err = SVCclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err, "\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a SVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a SVC for user %s, volume id %s", username, volumeID)
	} else {
		print("\n UPDATE!!")
		_, err = SVCclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a SVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a SVC for user %s, volume id %s", username, volumeID)
	}
	print("\n func FINISH")

	return nil
}

func (client *K8sClient) CreateAppSVC(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreateAppSVC",
	})

	logger.Debugf("Creating a AppSVC for user %s, volume id %s", username, volumeID)

	claim := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      client.getAppSVCName(volumeID),
			Namespace: client.getSvcNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     int32(60000),
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": client.getDeployAppName(volumeID),
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	appSVCclient := client.clientSet.CoreV1().Services(client.getVolumeNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := appSVCclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		_, err = appSVCclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err, "\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a AppSVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a AppSVC for user %s, volume id %s", username, volumeID)
	} else {
		_, err = appSVCclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a AppSVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a AppSVC for user %s, volume id %s", username, volumeID)
	}

	return nil
}
*/
