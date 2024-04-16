# Configure the K8S cluster
- Install the kubernetes packages, and establish connections between the master and worker nodes.
- Build a kubernetes cluster based on containerd.
- If you already have kubernetes installed,
remove all kubernetes-related processes (kube*) and
Remove all disconnected mount-points (/var/lib/kubelet/...).

## 1.Preparation

### 1.1.Set Hostname
Change the hostname to recognize that you are the master/worker node. This will be reflected on re-login.
```
sudo hostnamectl set-hostname master

sudo hostnamectl set-hostname worker1
```

### 1.2.Check Ubuntu OS information (can be omitted)
```
nproc # Check core
free -h # Check memory
ifconfig -a # Check Mac
sudo cat /sys/class/dmi/id/product_uuid # check product_uuid
```

### 1.3.Disable swap memory
Disable swap memory to use the kubelet.
```
sudo swapoff -a
sudo sed -ri '/\sswap\s/s/^#?/#/' /etc/fstab

sudo free -m # Verify that swap memory is 0
```

### 1.4.Set up a firewall
Allow the port numbers, http, and https required to configure the kubernetes cluster.
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

### 1.5.Load kernel modules required to run a Kubernetes cluster
Create a file called containerd.conf in the /etc/modules-load.d/ directory and add the overlay and br_netfilter kernel modules to the list of modules that should be loaded at boot time.
```
cat <<EOF | sudo tee /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter
```

  

### 1.6.sysctl parameter definition
```
cat <<EOF | sudo tee /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

sudo sysctl --system # Changes the sysctl parameters without restarting the system.
```

### 1.7.SELinux Mode Modifications
Set SELinux (Security Architecture) to permissive mode in order for Kubernetes to access the host file system required by the Pod network.
```
setenforce 0

sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

sestatus
```

## Install 2.ContainerD

### 2.1.install containerd.io
- containerd.io only installs certain packages. If you want to implement all the features of containerD, you can install containerd.
- It doesn't matter which one you use.
- On Raspberry Pi, install containerd instead of containerd.io.
```
sudo apt-get update
sudo apt-get install \
apt-transport-https \
ca-certificates \
curl \
gnupg \
lsb-release

# Install the packages required to use the apt package manager with the HTTPS repository and add the Docker GPG key

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update # Now that the new repositories have been added, update them

sudo apt-get install containerd.io

sudo apt-get install containerd

sudo systemctl status containerd # Verify installation
```

  

### 2.2.Create ContainerD configuration file
```
sudo mkdir -p /etc/containerd
sudo containerd config default | sudo tee /etc/containerd/config.toml
```

### Modify the 2.3.config.toml file
```
sudo nano /etc/containerd/config.toml # Search for SystemdCgroup and modify it to SystemdCgroup = true

sudo systemctl restart container # Apply the modified config.toml and restart container
```

## Install 3.k8s and join the cluster

### Prepare to install 3.1.k8s
- Retrieve the Kubernetes GPG key and add it to your system's keyring
- Add the Kubernetes repository to your system's apt sources list

```
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add

sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
```

### Install 3.2kubeadm, kubelet, and kubectl packages, version fixed
```
sudo apt-get install kubeadm kubelet kubectl -y

sudo apt-mark hold kubeadm kubelet kubectl containerd
```

### 3.3.k8s Importing images before cluster initialization (can be omitted)
Importing an image in advance makes the initialization process faster and more reliable.
```
sudo kubeadm config images pull
```

### Initializing a 3.4.k8s cluster (Master)
```
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
```

Running the above command will output commands that can be used to join the worker node to the cluster.

For example,
```

kubeadm join 192.168.0.100:6443 --token abcdef.1234567890abcdef \
--discovery-token-ca-cert-hash sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
```
Copy this command and use it later.

### Configure 3.5.kubectl
```
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

### 3.6.CNI Setup (Flannel)
Deploy Flannel in your kubernetes cluster.
```
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

### 3.7.Join the worker node to the cluster
Run the kubeadm join... command copied above on the worker node.
```
kubectl get nodes
```
and verify that the status changes to Ready.

[Note](https://velog.io/@chan9708/k8ssettings)