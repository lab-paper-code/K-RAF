# k8s 클러스터 구성
- kubernetes 패키지를 설치하고, master 노드와 worker 노드 간의 연결을 설정합니다.
- containerd 기반으로 kubernetes 클러스터를 구축합니다.
- 기존에 kubernetes가 설치된 환경이라면,
kubernetes 관련 프로세스 (kube*)를 모두 제거하고
연결이 끊어진 mount-point (/var/lib/kubelet/...)를 모두 제거합니다.

## 1.사전 준비

### 1.1.Hostname 설정
master/worker 노드임을 알 수 있게 hostname을 변경합니다. 재로그인시 반영됩니다.
```
sudo hostnamectl set-hostname master

sudo hostnamectl set-hostname worker1
```

### 1.2.Ubuntu OS 정보 확인(생략 가능)
```
nproc # 코어 확인
free -h # 메모리 확인
ifconfig -a # Mac 확인
sudo cat /sys/class/dmi/id/product_uuid # product_uuid 확인
```

### 1.3.swap 메모리 비활성화
kubelet을 사용하기 위해 swap 메모리를 비활성화합니다.
```
sudo swapoff -a
sudo sed -ri '/\sswap\s/s/^#?/#/' /etc/fstab

sudo free -m # swap 메모리가 0인지 확인
```

### 1.4.방화벽 설정
kubernetes 클러스터를 구성하는데 필요한 포트 번호와, http, https를 허용합니다.
```
sudo apt-get install firewalld
systemctl start firewalld
systemctl enable firewalld
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
  

firewall-cmd --permanent --add-port=2379-2380/tcp
firewall-cmd --permanent --add-port=6443/tcp
firewall-cmd --permanent --add-port=10250-10252/tcp
firewall-cmd --permanent --add-port=26443/tcp
firewall-cmd --permanent --add-port=30000-32767/tcp
  
firewall-cmd --permanent --add-port=8285/udp
firewall-cmd --permanent --add-port=8472/udp

firewall-cmd --reload
```

### 1.5.Kubernetes 클러스터 실행에 필요한 커널 모듈 로드
/etc/modules-load.d/ 디렉토리에 containerd.conf라는 파일을 생성하고 overlay 및 br_netfilter 커널 모듈을 부팅 시 로드해야 하는 모듈 목록에 추가합니다.
```
cat <<EOF | sudo tee /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter
```

  

### 1.6.sysctl 파라미터 정의
```
cat <<EOF | sudo tee /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

sudo sysctl --system # 시스템을 재시작하지 않고 sysctl 파라미터를 변경합니다.
```

### 1.7.SELinux 모드 수정
쿠버네티스가 파드 네트워크에 필요한 호스트 파일 시스템에 접근하기 위해 SELinux(보안 아키텍처)를 permissive 모드로 설정합니다.
```
setenforce 0

sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

sestatus
```

## 2.ContainerD 설치

### 2.1.containerd.io 설치
- containerd.io는 특정 패키지만 설치합니다. containerD의 모든 기능을 구현하려면 containerd를 설치하면 됩니다.
- 둘 중 어떤 것을 사용해도 상관없습니다.
- 라즈베리파이에서는 containerd.io 대신 containerd를 설치합니다.
```
sudo apt-get update
sudo apt-get install \
apt-transport-https \
ca-certificates \
curl \
gnupg \
lsb-release

# HTTPS 리포지토리와 함께 apt 패키지 관리자를 사용하는 데 필요한 패키지를 설치하고 Docker GPG 키를 추가

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update # 새로운 저장소가 추가되었으므로, 업데이트

sudo apt-get install containerd.io

sudo apt-get install containerd

sudo systemctl status containerd # 설치 확인
```

  

### 2.2.ContainerD 구성 파일 생성
```
sudo mkdir -p /etc/containerd
sudo containerd config default | sudo tee /etc/containerd/config.toml
```

### 2.3.config.toml 파일 수정
```
sudo nano /etc/containerd/config.toml # SystemdCgroup 을 검색해 SystemdCgroup = true로 수정

sudo systemctl restart containerd # 수정한 config.toml을 적용 및 재실행
```

## 3.k8s 설치 및 클러스터 가입

### 3.1.k8s 설치 준비
- Kubernetes GPG 키를 검색하여 시스템의 키링에 추가
- 시스템의 apt 소스 목록에 Kubernetes 레포지토리를 추가

```
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add

sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
```

### 3.2kubeadm, kubelet 및 kubectl 패키지 설치, 버전 고정
```
sudo apt-get install kubeadm kubelet kubectl -y

sudo apt-mark hold kubeadm kubelet kubectl containerd
```

### 3.3.k8s 클러스터 초기화 전 image 가져오기 (생략 가능)
미리 이미지를 가져오면 초기화 과정이 더 빠르고 안정적입니다.
```
sudo kubeadm config images pull
```

### 3.4.k8s 클러스터 초기화(Master)
```
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
```

위의 명령어를 실행하면 worker node를 클러스터에 가입하는데 사용할 수 있는 명령이 출력됩니다.

예를 들어,
```

kubeadm join 192.168.0.100:6443 --token abcdef.1234567890abcdef \
--discovery-token-ca-cert-hash sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
```
이 명령어는 복사해두고 이후에 사용합니다.

### 3.5.kubectl 구성
```
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

### 3.6.CNI 설정(Flannel)
kubernetes 클러스터에 Flannel을 배포합니다.
```
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

### 3.7.worker node를 클러스터에 가입
위에서 복사한 kubeadm join... 명령어를 worker node에서 실행합니다.
```
kubectl get nodes
```
를 이용하여 status가 Ready로 바뀌는지 확인합니다.

[참고](https://velog.io/@chan9708/k8ssettings)