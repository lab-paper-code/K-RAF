package k8s

/*
import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	deployWebdavSuffix    string = "-webdav"
	deployWebdavNamespace string = "vd"
	deployWebdavSecret    string = "webdav-secret"
	deployAppSuffix       string = "-app"
)

//TODO: 한번에 deploy 하려면 V1DepoymentList 사용
// https://github.com/kubernetes-client/go/blob/master/kubernetes/docs/V1DeploymentList.md

// WEBDAV

apiVersion: apps/v1
kind: Deployment
metadata:
  name: pod1-webdav   #변경
  namespace: ksv
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pod1-webdav#변경
  revisionHistoryLimit: 2
  template:
    metadata:
      labels:
        app: pod1-webdav   #변경
    spec:
      containers:
      - name: webdav
        image: yechae/ksv-webdav:v1
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 10
          periodSeconds: 10
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 10
          periodSeconds: 10
          failureThreshold: 3
        resources:
          requests:
            memory: "100Mi"
            cpu: "100m"
          limits:
            memory: "200Mi"
            cpu: "200m"
        volumeMounts:
        - mountPath: /uploads
          name: webdav-storage
        env:
          - name: BASIC_AUTH
            value: "True"
          - name: WEBDAV_LOGGIN
            value: "info"
          - name: WEBDAV_USERNAME
            valueFrom:
              secretKeyRef:
                name: pod1-secret   #변경
                key: "user"
          - name: WEBDAV_PASSWORD
            valueFrom:
              secretKeyRef:
                name: pod1-secret
                key: "password"
      volumes:
      - name: webdav-storage
        persistentVolumeClaim:
          claimName: pod1-pvc    #변경
      restartPolicy: Always
      # uncomment if registry keys are specified
      #imagePullSecrets:
      #- name: <secret_name>

func (client *K8sClient) getDeployWebdavName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, deployWebdavSuffix)
}
func (client *K8sClient) getDeployAppName(volumeID string) string {
	return fmt.Sprintf("%s%s", volumeID, deployAppSuffix)
}

func (client *K8sClient) getDeployLabels(username string, volumeID string) map[string]string {
	return map[string]string{
		"username":  username,
		"volume-id": volumeID,
	}
}

func (client *K8sClient) getDeployNamespace() string {
	return volumeNamespace
}

// CreateWebdavDeploy creates a deploy for the given volumeID
func (client *K8sClient) CreateWebdavDeploy(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreateWebdavDeploy",
	})

	logger.Debugf("Creating a Webdav Deploy for user %s, volume id %s", username, volumeID)

	deployWebdavName := client.getDeployWebdavName(volumeID)
	deployReplicas := int32(1)

	claim := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployWebdavName,
			Labels:    client.getDeployLabels(username, volumeID),
			Namespace: client.getDeployNamespace(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deployReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deployWebdavName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deployWebdavName,
					Labels: map[string]string{
						"app": deployWebdavName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "webdav",
							Image: "yechae/ksv-webdav:v2",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(80),
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
								FailureThreshold:    3,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(80),
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
								FailureThreshold:    3,
							},
							//  Resources: corev1.ResourceRequirements{
							//  	Requests: corev1.ResourceList{

							//  	}
							//  }
							//meatav1, appsv1, corev1
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/uploads",
									Name:      "webdav-storage",
								},
							},
							Env: []corev1.EnvVar{
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
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{
											Name: "webdav-secret"},
											Key: "user",
										},
									},
								},
								{
									Name: "WEBDAV_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{
											Name: "webdav-secret"},
											Key: "password",
										},
									},
								},
							},
						}, //Continers
					}, //Continers
					Volumes: []corev1.Volume{
						{
							Name: "webdav-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: client.getPVCName(volumeID),
								},
							},
						},
					},
					RestartPolicy: "Always",
				}, //spec
			},
		},
	}

	deployclient := client.clientSet.AppsV1().Deployments(client.getDeployNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := deployclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		_, err = deployclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err, "\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a Webdav Deploy for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a Webdav Deploy for user %s, volume id %s", username, volumeID)
	} else {
		_, err = deployclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a Webdav Deploy for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a Webdav Deploy for user %s, volume id %s", username, volumeID)
	}

	return nil
}

//APP
/*
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pod1-app #변경
  namespace: ksv

spec:
  replicas: 1
  selector:
    matchLabels:
      app: pod1-app #변경
  revisionHistoryLimit: 2
  template:
    metadata:
      labels:
        app: pod1-app #변경
    spec:
      containers:
        - name: app-image
          #image: yechae/ksv-app:v3
          image: yechae/kube-flask:v4
          imagePullPolicy: IfNotPresent
          ports:
          - containerPort: 5000
          volumeMounts:
          - mountPath: "/mnt"
            name: volumes
          resources:
            requests:
              cpu: "250m"
            limits:
              cpu: "500m"

      volumes:
      - name: volumes
        persistentVolumeClaim:
          claimName: pod1-pvc #변경
      restartPolicy: Always

// CreateAppDeploy creates a App deploy for the given volumeID
func (client *K8sClient) CreateAppDeploy(username string, volumeID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "k8s",
		"struct":   "K8sClient",
		"function": "CreateAppDeploy",
	})

	logger.Debugf("Creating a App Deploy for user %s, volume id %s", username, volumeID)

	deployAppName := client.getDeployAppName(volumeID)
	deployReplicas := int32(1)

	claim := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployAppName,
			Labels:    client.getDeployLabels(username, volumeID),
			Namespace: client.getDeployNamespace(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deployReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deployAppName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deployAppName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "app-image",
							Image:           "yechae/ksv-app:v4",
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 5000,
								},
							},
							// Resources: corev1.ResourceRequirements{
							// 	Requests: map[string]string{
							// 		cpu: "250m",

							// 	},
							// 	Limits: map[string]string{
							// 		cpu: "500m",
							// 	},
							// },
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/mnt",
									Name:      "volumes",
								},
							},
						}, //Continers
					}, //Continers
					Volumes: []corev1.Volume{
						{
							Name: "volumes",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: client.getPVCName(volumeID),
								},
							},
						},
					},
					RestartPolicy: "Always",
				}, //spec
			},
		},
	}

	deployclient := client.clientSet.AppsV1().Deployments(client.getDeployNamespace())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), k8sTimeout)
	defer cancel()

	_, err := deployclient.Get(ctx, claim.GetName(), metav1.GetOptions{})

	if err != nil {
		// failed to get an existing claim
		_, err = deployclient.Create(ctx, claim, metav1.CreateOptions{})
		if err != nil {
			print(err, "\n")
			// failed to create one
			log.Fatal(err)
			logger.Errorf("Failed to create a App Deploy for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Created a App Deploy for user %s, volume id %s", username, volumeID)
	} else {
		_, err = deployclient.Update(ctx, claim, metav1.UpdateOptions{})
		if err != nil {
			// failed to create one
			logger.Errorf("Failed to update a App Deploy for user %s, volume id %s", username, volumeID)
			return err
		}

		logger.Debugf("Updated a App Deploy for user %s, volume id %s", username, volumeID)
	}

	return nil
}
*/
