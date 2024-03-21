package k8s

/*
import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ingWebdavSuffix     string = "-webdav-ing"
	ingAppSuffix        string = "-app-ing"
	ingNamespace        string = "vd"
	ingWebdavPathSuffix string = "/"
	ingAppPathSuffix    string = "/app/"
)


Ingress example

apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pod1-ingress # 변경
  namespace: ksv
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "150"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "150"
spec:
  rules:
    - host:
      http:
        paths:
          - path: /pod1  # 변경 # volumeID 로
            pathType: Prefix
            backend:
              service:
                name: webdav-pod1-svc # 변경
                port:
                  number: 80


// getAppSvcName makes appIngress name
func (client *K8sClient) getAppIngressName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, ingAppSuffix)
}

func (client *K8sClient) getAppIngressPath(volumeID string) string {
	return fmt.Sprintf("%s%s", ingAppPathSuffix, volumeID)
}

func (client *K8sClient) getIngressNamespace() string {
	return ingNamespace
}

func (client *K8sClient) CreateAppIngress(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreateAppIngress",
	})

	logger.Debugf("Creating a App Ingress for user %s, volume id %s", username, volumeID)
	pathPrefix := networkingv1.PathTypePrefix

	claim := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      client.getAppIngressName(volumeID),
			Namespace: client.getIngressNamespace(),
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                       "nginx",
				"nginx.ingress.kubernetes.io/proxy-connect-timeout": "150",
				"nginx.ingress.kubernetes.io/proxy-read-timeout":    "150",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     client.getAppIngressPath(volumeID),
									PathType: &pathPrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: client.getAppSVCName(volumeID),
											Port: networkingv1.ServiceBackendPort{
												Number: 60000,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// claim := &extensionsv1beta1.Ingress{
	// 	ApiVersion: "networking.k8s.io/v1",
	// 	Kind: "Ingress",
	// 	ObjectMeta: metav1.ObjectMetat{
	// 		Name: client.getAppIngressName(volumeID),
	// 		Namespace: client.getIngressNamespace(),
	// 		Annotations: {
	// 			kubernetes.io/ingress.class: "nginx",
	// 			nginx.ingress.kubernetes.io/proxy-connect-timeout: "150",
	// 			nginx.ingress.kubernetes.io/proxy-read-timeout: "150",

	// 		},
	// 	},
	// 	Spec: []extensionsv1beta1.IngressRule{
	// 		Http: []extensionsv1beta1.HttpIngressPath{
	// 			Path: client.getAppIngressPath(volumeID),
	// 			Backend: extensionsv1beta1.IngressBackend{
	// 				ServiceName: client.getAppSVCName(volumeID),
	// 				ServicePort: map[string]String{
	// 					60000
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	appIngclient := client.clientSet.NetworkingV1().Ingresses(client.getVolumeNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := appIngclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		_, err = appIngclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err, "\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a appSVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a appSVC for user %s, volume id %s", username, volumeID)
	} else {
		_, err = appIngclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a appSVC for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a appSVC for user %s, volume id %s", username, volumeID)
	}

	return nil
}
*/
