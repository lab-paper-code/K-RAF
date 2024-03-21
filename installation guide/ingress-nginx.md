## bare metal 환경에서 ingress nginx 설치 

nginx ingress : 1.3.0 버전 설치
ingress-nignx github 참고해서 kubernetes 버전과 호환되는 버전 설치

* master 노드에 nginx ingress설치 하기 위해서 tain 설정 해제
```
# kubectl describe node master | grep Taints
Taints: node-role.kubernetes.io/control-plane:NoSchedule

# 설정 해제
 kubectl taint nodes –all node-role.kubernetes.io/master-

```
* affinity 설정
    * nginx ingress를 설치할 노드에 레이블
```
kubectl label nodes <your-node-name> <key>=<value>
kubectl label nodes master type=lb

# 조회
kubectl get nodes --show-labels
```

* bare metal용 controller 설치
```
 curl -L -o ingress-nginx-controller.yml \
https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.3.0/deploy/static/provider/baremetal/deploy.yaml
```

* __controller__, **admission**, **admission-patch** 세 파드는 master or 같은 노드에서 작동하도록 affinity 추가
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

* 베어메탈용 컨트롤러는 NodePort 로 동작함. 이 상태에서 외부 트래픽은 컨트롤러로 바로 접근할 수 없어 LoadBalancer 타입으로 변경해 주어야 합니다.
```
$ kubectl patch svc ingress-nginx-controller -n \
ingress-nginx -p '{"spec": {"type": "LoadBalancer"}}'

$ kubectl get svc -n ingress-nginx
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.96.161.225    <pending>     80:30723/TCP,443:32540/TCP   1m
ingress-nginx-controller-admission   ClusterIP      10.107.233.129   <none>        443/TCP                      1m
```

* LoadBalancer 타입으로 변경되었지만, EXTERNAL-IP 가 <pending> 상태로 할당되고 있지 않아 외부 접근은 여전히 불가능합니다. MetalLB 를 설치하면 외부에서 접근할 수 있으니 진행해 보겠습니다.

## MetalLB 설치하기
[설치 가이드](https://metallb.universe.tf/installation/)      


* 설치 전 k8s 설정 수정
    * "mode" 를 "ipvs"로, "ipvs.strictARP"를 "true"로 변경 
```
$ kubectl edit configmap -n kube-system kube-proxy
===
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
ipvs:
  strictARP: true
```

* manifest 설치 
```
$ curl -L -o namespace.yml https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/namespace.yaml

$ curl -L -o metallb.yml https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/metallb.yaml

$ vi metallb.yml
```

* deployment/controller 및 daemonset/speaker 에 노드 어피니티를 추가해 준 후 설치를 진행
```
$ kubectl apply -f namespace.yml
$ kubectl apply -f metallb.yml
```

* speaker 간 통신시 암호화를 위해 secret 이 필요
```
$ kubectl create secret generic -n metallb-system \
memberlist --from-literal=secretKey="$(openssl rand -base64 128)"
```

* configmap.yaml 필요
    * HOST_IP는 현재 노드(마스터)의 ipv4  주소 입력 
    * ip가 155.230.36.27 일 경우 155.230.36.27-155.230.36.27로 지정 
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

* 조금 기다린 후, 인그레스 컨트롤러의 EXTERNAL-IP 가 할당되는 것 확인하면 끝 


* ingress 예시
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