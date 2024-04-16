## Install ingress nginx in a bare metal environment 

nginx ingress: Install version 1.3.0
Refer to ingress-nignx github to install a version compatible with your kubernetes version

* Unset tain to install nginx ingress on master node
```
# kubectl describe node master | grep Taints
Taints: node-role.kubernetes.io/control-plane:NoSchedule

# unset
 kubectl taint nodes -all node-role.kubernetes.io/master-

```
* Set affinity
    * label nodes to install nginx ingress on
```
kubectl label nodes <your-node-name> <key>=<value>
kubectl label nodes master type=lb

# lookup
kubectl get nodes --show-labels
```

* Install controller for bare metal
```
 curl -L -o ingress-nginx-controller.yml \
https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.3.0/deploy/static/provider/baremetal/deploy.yaml
```

* __controller__, **admission**, **admission-patch** The three pods are master or Add affinity to work on the same node
```
# ingress-nginx-controller.yml
...
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
              - key: type
                operator: In
                values:
                  - lb
```

```
$ kubectl apply -f ingress-nginx-controller.yml
$ kubectl get pods -n ingress-nginx -o wide

NAME                                        READY   STATUS      RESTARTS       AGE   IP               NODE     NOMINATED NODE   READINESS GATES
ingress-nginx-admission-create-8pjzr        0/1     Completed   0              55d   192.168.219.92   master   <none>           <none>
ingress-nginx-admission-patch-n8hmr         0/1     Completed   1              55d   192.168.219.91   master   <none>           <none>
ingress-nginx-controller-6444cb45b5-nplrt   1/1     Running     5 (4h3m ago)   55d   192.168.219.94   master   <none>           <none>

```

* The controller for bare metal behaves as a NodePort. In this state, external traffic cannot directly access the controller and must be changed to the LoadBalancer type.
```
$ kubectl patch svc ingress-nginx-controller -n \
ingress-nginx -p '{"spec": {"type": "LoadBalancer"}}'

$ kubectl get svc -n ingress-nginx
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.96.161.225    <pending>     80:30723/TCP,443:32540/TCP   1m
ingress-nginx-controller-admission   ClusterIP      10.107.233.129   <none>        443/TCP                      1m
```

* Although it has been changed to the LoadBalancer type, external access is still not possible because EXTERNAL-IP is not being assigned with a <pending> status. Installing MetalLB will allow external access, so let's proceed.

## Install MetalLB
[Installation Guide](https://metallb.universe.tf/installation/)      


* Modify K8S settings before installation
    * Change "mode" to "ipvs" and "ipvs.strictARP" to "true" 
```
$ kubectl edit configmap -n kube-system kube-proxy
===
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
ipvs:
  strictARP: true
```

* Install manifest 
```
$ curl -L -o namespace.yml https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/namespace.yaml

$ curl -L -o metallb.yml https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/metallb.yaml

$ vi metallb.yml
```

* Add node affinity to deployment/controller and daemonset/speaker and proceed with the installation
```
$ kubectl apply -f namespace.yml
$ kubectl apply -f metallb.yml
```

* secret is required for encryption when communicating between speakers
```
$ kubectl create secret generic -n metallb-system \
memberlist --from-literal=secretKey="$(openssl rand -base64 128)"
```

* configmap.yaml required
  * HOST_IP is the ipv4 address of the current node (master) 
  * If IP is 155.230.36.27, specify 155.230.36.27-155.230.36.27 
```
# configmap.yml
---
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: metallb-system
  name: config
data:
  config: |
    address-pools:
    - name: default
      protocol: layer2
      addresses:
      - [HOST_IP]
```

```
$ kubectl apply -f configmap.yml
```

* After waiting for a while, check that the EXTERNAL-IP of the Ingress Controller is assigned. 


* ingress example
```
# ingress-nginx.yml
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-test
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
  rules:
    - host: [HOST NAME or MASTER IP]
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: application
                port:
                  number: 80
```
```
$ kubectl apply -f ingress-nginx.yml
```