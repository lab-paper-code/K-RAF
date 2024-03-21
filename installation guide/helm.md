# helm , helm chart 설치 
* ubuntu 20.04 환경에서 실행하였습니다.
* kubernetes가 설치되어 있어야 합니다.
* 모든 과정은 ksv 네임스페이스에서 이루어집니다.

## 1.helm 설치
linux ubuntu 20.04 환경에서 설치합니다.
아래의 커맨드로 파일을 다운로드한 후 스크립트를 실행합니다.
```
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

아래 명령어를 통해 설치여부와 버전을 확인할 수 있습니다.
```
root@master:~helm version
version.BuildInfo{Version:"v3.9.2", GitCommit:"1addefbfe665c350f4daf868a9adc5600cc064fd", GitTreeState:"clean", GoVersion:"go1.17.12"}
```

## 2.helm 사용하기
### 2.1.helm 기본 사용법
#### 2.1.1.helm repo add
helm에서 chart로 미리 배포된 repository를 사용하기 위해서는 repo를 추가해야 합니다. 
repo add 명령어로 repository를 다운받을 수 있습니다. 
```
helm repo add [NAME] [URL] [flags]
```
#### 2.1.2.helm repo list
아래 list 명령어로 실행중인 chart의 리스트를 확인할 수 있습니다.이 때 namespace를 확인하고 기입해야 합니다.
```
root@master:~helm list -A
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                            APP VERSION
jsonexporter    ksv             6               2022-11-01 23:06:29.259408297 +0900 KST deployed        prometheus-json-exporter-0.4.0   v0.5.0
prometheus      ksv             4               2022-10-31 20:19:34.776034436 +0900 KST deployed        kube-prometheus-stack-39.11.0    0.58.0
pvc-autoresizer ksv             1               2022-10-24 16:11:48.440464113 +0900 KST deployed        pvc-autoresizer-0.5.0            0.5.0     
```

#### 2.1.3.helm pull
helm install을 통해서 프로그램을 설치할 수 있지만 세부 설정(values.yaml 또는 templates 디렉토리 내부의 템플릿)을 수정하기 위해서는 먼저 chart를 다운로드해야 합니다.
```
helm pull [chart URL | repo/chartname] [...] [flags]
```

#### 2.1.4.helm install
차트를 설치합니다.
```
helm install [NAME] [CHART] [flags]
```
### 2.2.관련 프로그램 설치
#### 2.2.1.Prometheus 사용하기
prometheus, prometheus-operator, grafana, node-exporter 등 관련된 프로그램을 일괄적으로 설치할 수 있는 kube-prometheus-stack을 사용해 설치를 진행합니다. 

kube-prometheus-stack 파일을 다운로드합니다.
```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm list -A
helm pull prometheus-community/kube-prometheus-stack
tar -xzf kube-prometheus-stack-55.5.0.tgz
cd kube-prometheus-stack
```

nodeSelector(nodeAffinity, nodeAntiaffinity 기능도 마찬가지) 기능을 이용하기 위해서는 노드별 label 설정이 필요합니다.

노드 전체 label 확인
```
kubectl get nodes --show-label
```

노드에 label 추가
```
kubectl label nodes {노드 이름} {label key}={label value}
```

노드에 label 제거
```
kubectl label nodes {노드 이름} {label key}-
```

values.yaml 파일을 수정해서 nodeSelector 기능을 활용합니다.
kube-prometheus-stack 디렉토리 내 values.yaml에서 nodeSelector를 검색합니다.
{label key}-{label value} 쌍을 추가합니다.

kube-prometheus-stack/charts 디렉토리에는 종속 chart 관련 파일들이 있습니다.
grafana/, kube-state-metrics/, prometheus-node-exporter/, prometheus-windows-exporter/ 아래 values.yaml에서도 nodeSelector를 추가합니다.

kube-prometheus-stack/ 경로에서
```
helm install prometheus . -n ksv
```
로 프로그램을 설치합니다.


#### 2.2.2.Json-exporter 사용하기
json-exporter는 위의 prometheus-community의 repo 내부에 포함되어 있습니다. 
pull을 통해 파일을 가져옵니다.

```

helm pull prometheus-community/prometheus-json-exporter
tar -xzf prometheus-json-exporter-0.9.0.tgz
```

마찬가지로 prometheus-json-exporter/의 values.yaml에 nodeSelector 를 추가해줍니다.
prometheus-json-exporter/ 경로에서
```
helm install json-exporter . -n ksv
```
로 프로그램을 설치합니다.

Json-exporter에서 사용할 API서버 사용하기 (실디바이스)
- 사용자
- python 3.8 버전 이상이 필요합니다.


#### 2.2.3.Pvc-autoresizer 사용하기 
pvc-autoresizer는 따로 repository를 추가해야 합니다. 
```
helm repo add pvc-autoresizer https://topolvm.github.io/pvc-autoresizer
helm repo update
helm pull pvc-autoresizer/pvc-autoresizer
tar -xzf pvc-autoresizer-0.10.1.tgz
```
pvc-autoresizer/values.yaml 과 pvc-autoresizer/charts/cert-manager/values.yaml에 nodeSelector를 추가합니다.

pvc-autoreizer/ 경로에서 
```
helm install json-exporter . -n ksv
```
로 설치합니다.

#### 2.2.4.Grafana 접속 & 대시보드 적용
그라파나에 접속하기 전 그라파나 서비스의 타입을 확인해야 합니다. 
타입이 NodePort일 경우 노출된 포트로 접근이 가능하지만 ClusterIP일 경우 이를 NodePort로 변경해야 합니다. 
```
root@master:~kubectl get service -n ksv
...
prometheus-grafana                        NodePort    10.101.233.184   <none>        80:32240/TCP                 25d
prometheus-kube-prometheus-alertmanager   ClusterIP   10.96.31.174     <none>        9093/TCP                     25d
prometheus-kube-prometheus-operator       ClusterIP   10.108.93.187    <none>        443/TCP                      25d
...
```

타입 변경시 helm chart의 values.yaml 파일에서 type을 바꿀 수 있고 아래의 명령어를 통해서도 변경이 가능합니다. 
```
kubectl patch svc prometheus-grafana -n ksv -p '{"spec": {"type": "NodePort"}}'
```

초기 로그인 정보인 admin/prom-operator 를 입력하면 아래와 같은 화면이 나옵니다. 
   
![grafana_main](./img/grafana_main.png) 
  
좌측 상단의 Dashboards 버튼에서 최하단의 import 버튼을 클릭합니다. 

![import_btn](./img/import_btn.png)

import 화면이 뜨고 upload JSON file 버튼을 클릭 후 [DEALLAB] pod disk usage.json 파일을 열고 적용하면 대시보드를 확인할 수 있습니다. 이때 메트릭 및 디자인은 수정할 수 있습니다.

![import_json](./img/import_json.png) 
![dashboard](./img/dashboard.png) 

[참고](https://ksr930.tistory.com/315): helm values.yaml 수정
[참고](https://ksr930.tistory.com/298): helm install