package k8s

import (
	"context"
	"fmt"

	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	webdavDeploymentNamePrefix string = "webdav"
	webdavDeploymentNamespace  string = objectNamespace
	webdavServiceNamePrefix    string = "webdav"
	webdavServiceNamespace     string = objectNamespace
	webdavIngressNamePrefix    string = "webdav"
	webdavIngressNamespace     string = objectNamespace

	webdavContainerVolumeName  string = "webdav-storage"
	webdavContainerPVMountPath string = "/uploads"
)

func (adapter *K8SAdapter) GetWebdavDeploymentName(volumeID string) string {
	return makeValidObjectName(webdavDeploymentNamePrefix, volumeID)
}

func (adapter *K8SAdapter) GetWebdavServiceName(volumeID string) string {
	return makeValidObjectName(webdavServiceNamePrefix, volumeID)
}

func (adapter *K8SAdapter) GetWebdavIngressName(volumeID string) string {
	return makeValidObjectName(webdavIngressNamePrefix, volumeID)
}

func (adapter *K8SAdapter) getWebdavDeploymentLabels(volume *types.Volume) map[string]string {
	labels := map[string]string{}
	labels["webdav-name"] = adapter.GetWebdavDeploymentName(volume.ID)
	labels["volume-id"] = volume.ID
	labels["device-id"] = volume.DeviceID
	return labels
}

func (adapter *K8SAdapter) getWebdavServiceLabels(volume *types.Volume) map[string]string {
	labels := map[string]string{}
	labels["webdav-name"] = adapter.GetWebdavServiceName(volume.ID)
	labels["volume-id"] = volume.ID
	labels["device-id"] = volume.DeviceID
	return labels
}

func (adapter *K8SAdapter) getWebdavIngressLabels(volume *types.Volume) map[string]string {
	labels := map[string]string{}
	labels["webdav-name"] = adapter.GetWebdavIngressName(volume.ID)
	labels["volume-id"] = volume.ID
	labels["device-id"] = volume.DeviceID
	return labels
}

func (adapter *K8SAdapter) GetWebdavIngressPath(volumeID string) string {
	return fmt.Sprintf("/%s", volumeID)
}

func (adapter *K8SAdapter) getWebdavContainers(device *types.Device, volume *types.Volume) []apiv1.Container {
	webdavPVMountPath := webdavContainerPVMountPath

	return []apiv1.Container{
		{
			Name:  "webdav",
			Image: "daclab/ksv-webdav:v2",
			Ports: []apiv1.ContainerPort{
				{
					ContainerPort: 80,
				},
			},
			LivenessProbe: &apiv1.Probe{
				ProbeHandler: apiv1.ProbeHandler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/",
						Port: intstr.FromInt(80),
					},
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				FailureThreshold:    3,
			},
			ReadinessProbe: &apiv1.Probe{
				ProbeHandler: apiv1.ProbeHandler{
					HTTPGet: &apiv1.HTTPGetAction{
						Path: "/",
						Port: intstr.FromInt(80),
					},
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				FailureThreshold:    3,
			},
			VolumeMounts: []apiv1.VolumeMount{
				{
					Name:      webdavContainerVolumeName,
					MountPath: webdavPVMountPath,
				},
			},
			Env: []apiv1.EnvVar{
				{
					Name:  "BASIC_AUTH",
					Value: "True",
				},
				{
					Name:  "WEBDAV_LOGGIN",
					Value: "info",
				},
				{
					Name: "WEBDAV_USERNAME",
					ValueFrom: &apiv1.EnvVarSource{
						SecretKeyRef: &apiv1.SecretKeySelector{
							LocalObjectReference: apiv1.LocalObjectReference{Name: adapter.GetSecretName(device)},
							Key:                  "username",
						},
					},
				},
				{
					Name: "WEBDAV_PASSWORD",
					ValueFrom: &apiv1.EnvVarSource{
						SecretKeyRef: &apiv1.SecretKeySelector{
							LocalObjectReference: apiv1.LocalObjectReference{Name: adapter.GetSecretName(device)},
							Key:                  "password",
						},
					},
				},
			},
		},
	}
}

func (adapter *K8SAdapter) getWebdavContainerVolumes(volume *types.Volume) []apiv1.Volume {
	pvcName := adapter.GetVolumeClaimName(volume.ID)

	containerVolumes := []apiv1.Volume{
		{
			Name: webdavContainerVolumeName,
			VolumeSource: apiv1.VolumeSource{
				PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
					ReadOnly:  false,
				},
			},
		},
	}
	return containerVolumes
}

func (adapter *K8SAdapter) createWebdavDeployment(device *types.Device, volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "createWebdavDeployment",
	})

	logger.Debug("received createWebdavDeployment()")

	webdavDeploymentName := adapter.GetWebdavDeploymentName(volume.ID)
	webdavDeploymentLabels := adapter.getWebdavDeploymentLabels(volume)
	webdavDeploymentNumReplicas := int32(1)

	webdavContainers := adapter.getWebdavContainers(device, volume)
	webdavContainerVolumes := adapter.getWebdavContainerVolumes(volume)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webdavDeploymentName,
			Labels:    webdavDeploymentLabels,
			Namespace: webdavDeploymentNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &webdavDeploymentNumReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"webdav-name": webdavDeploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   webdavDeploymentName,
					Labels: webdavDeploymentLabels,
				},
				Spec: apiv1.PodSpec{
					Containers:    webdavContainers,
					Volumes:       webdavContainerVolumes,
					RestartPolicy: apiv1.RestartPolicyAlways,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(webdavDeploymentNamespace)
	_, err := deploymentclient.Get(ctx, deployment.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist
		_, createErr := deploymentclient.Create(ctx, deployment, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := deploymentclient.Update(ctx, deployment, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) deleteWebdavDeployment(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteWebdavDeployment",
	})

	logger.Debug("received deleteWebdavDeployment()")

	webdavDeploymentName := adapter.GetWebdavDeploymentName(volumeID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(webdavDeploymentNamespace)
	err := deploymentclient.Delete(ctx, webdavDeploymentName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) createWebdavService(volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "createWebdavService",
	})

	logger.Debug("received createWebdavService()")

	webdavServiceName := adapter.GetWebdavServiceName(volume.ID)
	webdavServiceLabels := adapter.getWebdavServiceLabels(volume)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webdavServiceName,
			Labels:    webdavServiceLabels,
			Namespace: webdavServiceNamespace,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port:     int32(80),
					Protocol: apiv1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"webdav-name": adapter.GetWebdavDeploymentName(volume.ID),
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(webdavServiceNamespace)
	_, err := serviceClient.Get(ctx, service.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist
		_, createErr := serviceClient.Create(ctx, service, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := serviceClient.Update(ctx, service, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) deleteWebdavService(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteWebdavService",
	})

	logger.Debug("received deleteWebdavService()")

	webdavServiceName := adapter.GetWebdavServiceName(volumeID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(webdavServiceNamespace)
	err := serviceClient.Delete(ctx, webdavServiceName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) createWebdavIngress(volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "createWebdavIngress",
	})

	logger.Debug("received createWebdavIngress()")

	webdavIngressName := adapter.GetWebdavIngressName(volume.ID)
	webdavIngressLabels := adapter.getWebdavIngressLabels(volume)

	pathPrefix := networkingv1.PathTypePrefix

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webdavIngressName,
			Labels:    webdavIngressLabels,
			Namespace: webdavIngressNamespace,
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
									Path:     adapter.GetWebdavIngressPath(volume.ID),
									PathType: &pathPrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: adapter.GetWebdavServiceName(volume.ID),
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
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

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	ingressClient := adapter.clientSet.NetworkingV1().Ingresses(webdavIngressNamespace)
	_, err := ingressClient.Get(ctx, ingress.GetName(), metav1.GetOptions{})
	if err != nil {
		// does not exist
		_, createErr := ingressClient.Create(ctx, ingress, metav1.CreateOptions{})
		if createErr != nil {
			return createErr
		}
	} else {
		// exist -> update
		_, updateErr := ingressClient.Update(ctx, ingress, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (adapter *K8SAdapter) deleteWebdavIngress(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteWebdavIngress",
	})

	logger.Debug("received deleteWebdavIngress()")

	webdavIngressName := adapter.GetWebdavIngressName(volumeID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	ingressClient := adapter.clientSet.NetworkingV1().Ingresses(webdavIngressNamespace)
	err := ingressClient.Delete(ctx, webdavIngressName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) CreateWebdav(device *types.Device, volume *types.Volume) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreateWebdav",
	})

	logger.Debug("received CreateWebdav()")

	err := adapter.createWebdavDeployment(device, volume)
	if err != nil {
		return err
	}

	err = adapter.createWebdavService(volume)
	if err != nil {
		return err
	}

	err = adapter.createWebdavIngress(volume)
	if err != nil {
		panic(err)
	}

	return nil
}

func (adapter *K8SAdapter) DeleteWebdav(volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "DeleteWebdav",
	})

	logger.Debug("received DeleteWebdav()")

	err := adapter.deleteWebdavIngress(volumeID)
	if err != nil {
		return err
	}

	err = adapter.deleteWebdavService(volumeID)
	if err != nil {
		return err
	}

	err = adapter.deleteWebdavDeployment(volumeID)
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) EnsureDeleteWebdav(volumeID string) {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "EnsureDeleteWebdav",
	})

	logger.Debug("received EnsureDeleteWebdav()")

	adapter.deleteWebdavIngress(volumeID)
	adapter.deleteWebdavService(volumeID)
	adapter.deleteWebdavDeployment(volumeID)
}
