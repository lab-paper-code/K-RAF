package k8s

/*
import (
	"context"
	"fmt"
	"os"

	"github.com/rook/rook/pkg/client/clientset/versioned/scheme"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// webdav exec
// sed -i -e 's#Alias /uploads \"/uploads\"#Alias /<volumdID>/uploads \"/uploads\"#g' /etc/apache2/conf.d/dav.conf

func (client *K8sClient) ExecInPod(namespace string, volumeID string, command string) error {

	// pod, err := client.PodClient().Get(podName, metav1.GetOptions{})
	// if err!= nil {
	// 	panic(err)
	// }

	// New ~ using the passed context name
	// return &DeferredLoadingClientConfig
	// ( use most recent rules even when loading rules change )
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	// corev1client change to client
	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	// vd was webdavnamespace name
	// get DeploywebdavName : (%s%s, VolumeID, WebdavSuffix)
	podLabel := client.getDeployWebdavName(volumeID)
	pods, err := client.clientSet.CoreV1().Pods("vd").List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=" + podLabel,
	})

	if err != nil {
		panic(err)
	}

	var podName string
	// var podObj corev1.Pod
	for _, pod := range pods.Items {
		podName = pod.Name
		// podObj = pod
	}

	execCommand := []string{
		"sh",
		"-c",
		command,
	}

	fmt.Println(execCommand)
	req := coreclient.RESTClient().
		Post().
		Namespace(namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			// Container: podObj.Spec.Containers[0].Name,
			Command: execCommand,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
	// NewSPDYExecutor connects to the provided server and
	// upgrades the connection to multiplexed bidirectional streams

	// initiates the transport of the standard shell streams
	// can transport except for nil streams
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		print(err)
	}

	return nil
}

func (client *K8sClient) getPodName(volumeID string) string {

	podLabel := client.getDeployWebdavName(volumeID)
	pods, err := client.clientSet.CoreV1().Pods("vd").List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=" + podLabel,
	})

	if err != nil {
		panic(err)
	}

	var podName string
	for _, pod := range pods.Items {
		podName = pod.Name
	}
	fmt.Println(podName)

	return podName
}

func (client *K8sClient) WaitPodRun3(username string, volumeID string) error {
	// pod, err := client.PodClient().Get(podName, metav1.GetOptions{})
	// if err!= nil {
	// 	panic(err)
	// }

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	namespace := client.getDeployNamespace()
	podName := client.getPodName(volumeID)

	pod, err := client.clientSet.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	ctx := context.Background()
	watcher, err := coreclient.Pods(namespace).Watch(
		ctx,
		metav1.SingleObject(pod.ObjectMeta),
	)
	if err != nil {
		return err
	}

	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.ResultChan():
			pod := event.Object.(*corev1.Pod)

			if pod.Status.Phase == corev1.PodRunning {
				fmt.Printf("The POD is running")
				return nil
			}

		case <-ctx.Done():
			fmt.Printf("Exit from waitPodRunning for POD \"%s\" because the context is done", volumeID)
			return nil
		}
	}

	return nil
}
*/
