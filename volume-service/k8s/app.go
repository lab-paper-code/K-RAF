package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

const (
	appDeploymentNamePrefix string = "app"
	appDeploymentNamespace  string = objectNamespace
	appServiceNamePrefix    string = "app"
	appServiceNamespace     string = objectNamespace
	appIngressNamePrefix    string = "app"
	appIngressNamespace     string = objectNamespace

	appContainerVolumeName  string = "app-storage"
	appContainerPVMountPath string = "/uploads"
)

func (adapter *K8SAdapter) GetAppDeploymentName(appRunID string) string {
	return makeValidObjectName(appDeploymentNamePrefix, appRunID)
}

func (adapter *K8SAdapter) GetAppServiceName(appRunID string) string {
	return makeValidObjectName(appServiceNamePrefix, appRunID)
}

func (adapter *K8SAdapter) GetAppIngressName(appRunID string) string {
	return makeValidObjectName(appIngressNamePrefix, appRunID)
}

func (adapter *K8SAdapter) getAppDeploymentLabels(appRun *types.AppRun) map[string]string {
	labels := map[string]string{}
	labels["app-name"] = adapter.GetAppDeploymentName(appRun.ID)
	labels["app-id"] = appRun.AppID
	labels["apprun-id"] = appRun.ID
	labels["volume-id"] = appRun.VolumeID
	labels["device-id"] = appRun.DeviceID
	return labels
}

func (adapter *K8SAdapter) getAppServiceLabels(appRun *types.AppRun) map[string]string {
	labels := map[string]string{}
	labels["app-name"] = adapter.GetAppServiceName(appRun.ID)
	labels["app-id"] = appRun.AppID
	labels["apprun-id"] = appRun.ID
	labels["volume-id"] = appRun.VolumeID
	labels["device-id"] = appRun.DeviceID
	return labels
}

func (adapter *K8SAdapter) getAppIngressLabels(appRun *types.AppRun) map[string]string {
	labels := map[string]string{}
	labels["app-name"] = adapter.GetAppIngressName(appRun.ID)
	labels["app-id"] = appRun.AppID
	labels["apprun-id"] = appRun.ID
	labels["volume-id"] = appRun.VolumeID
	labels["device-id"] = appRun.DeviceID
	return labels
}

func (adapter *K8SAdapter) GetAppIngressPath(appRunID string) string {
	return fmt.Sprintf("/app/%s", appRunID)
}

func (adapter *K8SAdapter) getAppContainers(app *types.App, device *types.Device, volume *types.Volume) []apiv1.Container {
	containerPorts := []apiv1.ContainerPort{}
	for _, port := range app.OpenPorts {
		containerPorts = append(containerPorts, apiv1.ContainerPort{
			Name:          fmt.Sprintf("cont-port-%d", port),
			ContainerPort: int32(port),
		})
	}
	gpuFlag := "0"
	// set to 1 if app requires GPU
	if app.RequireGPU {
		gpuFlag = "1"
	}

	cmdString := app.Commands
	commands := strings.Split(cmdString, " ")

	argString := app.Arguments
	arguments := strings.Split(argString, " ")

	// Create a container object
	container := apiv1.Container{
		Name:            "app",
		Image:           app.DockerImage,
		ImagePullPolicy: "IfNotPresent",
		Ports:           containerPorts,
		VolumeMounts: []apiv1.VolumeMount{
			{
				Name:      appContainerVolumeName,
				MountPath: appContainerPVMountPath,
			},
		},
		SecurityContext: &apiv1.SecurityContext{
			Privileged: pointer.Bool(true),
		},
	}

	// Conditionally set Command and Args if they are not empty
	if cmdString != "" {
		container.Command = commands
	}
	if argString != "" {
		container.Args = arguments
	}

	// Conditionally set GPU Limits if gpuFlag is not "0"
	if gpuFlag != "0" {
		container.Resources = apiv1.ResourceRequirements{
			Limits: apiv1.ResourceList{
				"nvidia.com/gpu": resource.MustParse(gpuFlag),
			},
		}
	}

	return []apiv1.Container{container}
}

/*
func (adapter *K8SAdapter) getAppContainers(app *types.App, device *types.Device, volume *types.Volume) []apiv1.Container {
	containerPorts := []apiv1.ContainerPort{}
	for _, port := range app.OpenPorts {
		containerPorts = append(containerPorts, apiv1.ContainerPort{
			Name:          fmt.Sprintf("cont-port-%d", port),
			ContainerPort: int32(port),
		})
	}
	gpuFlag := "0"
	// set to 1 if app requires GPU
	if app.RequireGPU {
		gpuFlag = "1"

	}

	cmdString := app.Commands

	commands := strings.Split(cmdString, " ")
	// split commands into slice

	argString := app.Arguments
	// variable to store app.Arguments

	arguments := strings.Split(argString, " ")
	// split arguments into slice



	return []apiv1.Container{
		{
			Name:            "app",
			Image:           app.DockerImage,
			ImagePullPolicy: "IfNotPresent",
			Ports:           containerPorts,

			VolumeMounts: []apiv1.VolumeMount{
				{
					Name:      appContainerVolumeName,
					MountPath: appContainerPVMountPath,
				},
			},
			Resources: apiv1.ResourceRequirements{
				Limits: apiv1.ResourceList{
					"nvidia.com/gpu": resource.MustParse(gpuFlag),
				},
			},
			// add command to container
			Command: commands,

			// add argument to container
			Args: arguments,

			SecurityContext: &apiv1.SecurityContext{
				Privileged: pointer.Bool(true),
			},
		},
	}
}
*/

func (adapter *K8SAdapter) getAppContainerVolumes(volume *types.Volume) []apiv1.Volume {
	pvcName := adapter.GetVolumeClaimName(volume.ID)

	containerVolumes := []apiv1.Volume{
		{
			Name: appContainerVolumeName,
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

func (adapter *K8SAdapter) createAppDeployment(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "createAppDeployment",
	})

	logger.Debug("received createAppDeployment()")

	appDeploymentName := adapter.GetAppDeploymentName(appRun.ID)
	appDeploymentLabels := adapter.getAppDeploymentLabels(appRun)
	deployReplicas := int32(1)

	appContainers := adapter.getAppContainers(app, device, volume)
	appContainerVolumes := adapter.getAppContainerVolumes(volume)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appDeploymentName,
			Labels:    appDeploymentLabels,
			Namespace: appDeploymentNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deployReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app-name": appDeploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   appDeploymentName,
					Labels: appDeploymentLabels,
				},
				Spec: apiv1.PodSpec{
					Containers:    appContainers,
					Volumes:       appContainerVolumes,
					RestartPolicy: apiv1.RestartPolicyAlways,
				}, //spec
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(appDeploymentNamespace)
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

func (adapter *K8SAdapter) deleteAppDeployment(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteAppDeployment",
	})

	logger.Debug("received deleteAppDeployment()")

	appDeploymentName := adapter.GetAppDeploymentName(appRunID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	deploymentclient := adapter.clientSet.AppsV1().Deployments(appDeploymentNamespace)
	err := deploymentclient.Delete(ctx, appDeploymentName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) createAppService(app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "createAppService",
	})

	logger.Debug("received createAppService()")

	appServiceName := adapter.GetAppServiceName(appRun.ID)
	appServiceLabels := adapter.getAppServiceLabels(appRun)

	servicePorts := []apiv1.ServicePort{}
	for _, port := range app.OpenPorts {
		servicePorts = append(servicePorts, apiv1.ServicePort{
			Name:     fmt.Sprintf("svc-port-%d", port),
			Port:     int32(port),
			Protocol: apiv1.ProtocolTCP,
		})
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appServiceName,
			Labels:    appServiceLabels,
			Namespace: appServiceNamespace,
		},
		Spec: apiv1.ServiceSpec{
			Ports: servicePorts,
			Selector: map[string]string{
				"app-name": adapter.GetAppDeploymentName(appRun.ID),
			},
			Type: apiv1.ServiceTypeClusterIP,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(volumeNamespace)
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

func (adapter *K8SAdapter) deleteAppService(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteAppService",
	})

	logger.Debug("received deleteAppService()")

	appServiceName := adapter.GetAppServiceName(appRunID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	serviceClient := adapter.clientSet.CoreV1().Services(appServiceNamespace)
	err := serviceClient.Delete(ctx, appServiceName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) createAppIngress(app *types.App, appRun *types.AppRun) error {

	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sAdapter",
		"function": "createAppIngress",
	})

	logger.Debug("received createAppIngress()")

	appIngressName := adapter.GetAppIngressName(appRun.ID)
	appIngressLabels := adapter.getAppIngressLabels(appRun)

	pathPrefix := networkingv1.PathTypePrefix

	serviceBackendPort := 0
	if len(app.OpenPorts) > 0 {
		serviceBackendPort = app.OpenPorts[0]
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appIngressName,
			Labels:    appIngressLabels,
			Namespace: appIngressNamespace,
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
									Path:     adapter.GetAppIngressPath(appRun.ID),
									PathType: &pathPrefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: adapter.GetAppServiceName(appRun.ID),
											Port: networkingv1.ServiceBackendPort{
												Number: int32(serviceBackendPort),
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

	ingressClient := adapter.clientSet.NetworkingV1().Ingresses(volumeNamespace)
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

func (adapter *K8SAdapter) deleteAppIngress(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "deleteAppIngress",
	})

	logger.Debug("received deleteAppIngress()")

	appIngressName := adapter.GetAppIngressName(appRunID)

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	ingressClient := adapter.clientSet.NetworkingV1().Ingresses(appIngressNamespace)
	err := ingressClient.Delete(ctx, appIngressName, *metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) CreateApp(device *types.Device, volume *types.Volume, app *types.App, appRun *types.AppRun) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "CreateApp",
	})

	logger.Debug("received CreateApp()")

	err := adapter.createAppDeployment(device, volume, app, appRun)
	if err != nil {
		return err
	}

	err = adapter.createAppService(app, appRun)
	if err != nil {
		return err
	}

	err = adapter.createAppIngress(app, appRun)
	if err != nil {
		panic(err)
	}

	return nil
}

func (adapter *K8SAdapter) DeleteApp(appRunID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "DeleteApp",
	})

	logger.Debug("received DeleteApp()")

	err := adapter.deleteAppIngress(appRunID)
	if err != nil {
		return err
	}

	err = adapter.deleteAppService(appRunID)
	if err != nil {
		return err
	}

	err = adapter.deleteAppDeployment(appRunID)
	if err != nil {
		return err
	}

	return nil
}

func (adapter *K8SAdapter) EnsureDeleteApp(appRunID string) {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "EnsureDeleteApp",
	})

	logger.Debug("received EnsureDeleteApp()")

	adapter.deleteAppIngress(appRunID)
	adapter.deleteAppService(appRunID)
	adapter.deleteAppDeployment(appRunID)
}

/*

func (adapter *K8SAdapter) sendAppCommand(podName string, command ...string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "sendAppCommand",
	})

	logger.Debug("received sendAppCommand()")

	args := append([]string{"exec", "-it", podName, "--"}, command...)

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (adapter *K8SAdapter) ExecuteAppCommand(appRunID string) {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8SAdapter",
		"function": "ExecAppCmd",
	})

	logger.Debug("received ExecuteAppCommand()")

	podName := appRunID // get podName with appRunID
	command := []string{"bash"}

	err := adapter.sendAppCommand(podName, command...)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Command executed successfully!")
	}

}
*/

/*
   //TODO:  k8s resource들 생성한 후
   //1. webdav pod으로 exec 명령어로 sed -i -e 's#Alias /uploads \"/uploads\"#Alias /<volumeID>/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf 명령어 실행
   //2. app pod으로 http://ip:60000/hello_flask?ip=<ip> 해서 dom ip 알려주기


   type Output struct {
       Mount  string       `json:mountPath`
       Device types.Device `json: device`
   }


   Mount := "http://155.230.36.27/" + volumeID + "/uploads"
   // 지금과 다름
   device = types.Device{
       IP:       input.IP,
       ID:       volumeID,
       Username: input.Username,
       Password: input.Password,
       Storage:  input.Storage,
   }


   output := Output{
       Mount:  Mount,
       Device: device,
   }
*/
