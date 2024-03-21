[Rook 공식문서](https://rook.io/docs/rook/v1.9/quickstart.html)

rook version : 1.9.10
ceph version : 16.2.10

## [중요] k8s가 heterogeneous cluster로 구성되었을 때
* 예를 들어 k8s cluster가 서버 3대와 라즈베리파이 2대로 구성되었을 때
* 라즈베리파이의 자원을 사용하지 않게 하는 방법입니다.   

* 아래의 ceph 코드 clone은 동일합니다. clone 후 rook/deploy/examples 로 이동합니다.
  * crds, common, operator는 그대로 생성해줍니다.
```
cd rook/deploy/examples
kubectl create -f crds.yaml -f common.yaml -f operator.yaml
```

* 라즈베리 파이인 노드들에 대해 해당 label을 추가해줍니다.
```
kubectl label node <raspberrypi_node_name> rookrole=no-rook-node
```

* 이후 KSV repository의 k8s/rook-ceph의 cluster.yaml을 실행합니다.
```
cd ~/KSV/k8s/
kubectl create -f cluster.yaml
```

* cluster 구성 및 OSD pod가 성공적으로 배포된 후에  rook ceph 설치 - ceph toolbox 설치 부터 진행합니다.

## rook ceph 설치 

* ceph 코드 clone 받아 오기
```
git clone --single-branch --branch v1.9.10 https://github.com/rook/rook.git
```

* ceph cluster 구성
```
cd rook/deploy/examples
kubectl create -f crds.yaml -f common.yaml -f operator.yaml
kubectl create -f cluster.yaml
```

* rook operator 생성
```
cd deploy/examples
kubectl create -f crds.yaml -f common.yaml -f operator.yaml

# verify the rook-ceph-operator is in the `Running` state before proceeding
kubectl -n rook-ceph get pod
```

* ceph toolbox 설치
[공식문서](https://rook.io/docs/rook/v1.8/ceph-toolbox.html)
```
# toolbox 파드 런칭
kubectl create -f deploy/examples/toolbox.yaml

# toolbox 파드 설치 상태 확인
kubectl -n rook-ceph rollout status deploy/rook-ceph-tools
# 파드 런칭이 다 되면 rook-ceph-tools 파드에 접속
kubectl -n rook-ceph exec -it deploy/rook-ceph-tools -- bash

# 아래의 명령어로도 접속 가능
kubectl -n rook-ceph exec -it $(kubectl -n rook-ceph get pod -l "app=rook-ceph-tools" -o jsonpath='{.items[0].metadata.name}') -- bash
```
- 아래의 실행 명령어로 ceph 상태 확인 가능
    - ceph status
    - ceph osd status
    - ceph df
    - rados df
    - ceph version
    - rook version



* ceph 설치 확인
```
$ kubectl get cephcluster -A
NAMESPACE   NAME        DATADIRHOSTPATH   MONCOUNT   AGE   PHASE   MESSAGE                        HEALTH      EXTERNAL
rook-ceph   rook-ceph   /var/lib/rook     3          70m   Ready   Cluster created successfully   HEALTH_OK
```

**ceph cluster 구성 완료**


## Shared Filesystem
[Rook 공식문서](https://rook.io/docs/rook/v1.9/ceph-filesystem.html)

* shared filesystem의 정의인 filesystem.yaml 생성
```
apiVersion: ceph.rook.io/v1
kind: CephFilesystem
metadata:
  name: myfs
  namespace: rook-ceph
spec:
  metadataPool:
    replicated:
      size: 3
  dataPools:
    - name: replicated
      replicated:
        size: 3
  preserveFilesystemOnDelete: true
  metadataServer:
    activeCount: 1
    activeStandby: true
``` 
* filesystem 생성
  * rook/deploy/examples/filesystem.yaml
```
# Create the filesystem
kubectl create -f filesystem.yaml

# 생성 확인
kubectl -n rook-ceph get pod -l app=rook-ceph-mds

NAME                                      READY     STATUS    RESTARTS   AGE
rook-ceph-mds-myfs-7d59fdfcf4-h8kw9       1/1       Running   0          12s
rook-ceph-mds-myfs-7d59fdfcf4-kgkjp       1/1       Running   0          12s
```

* storageclass.yaml 생성
  * rook/deploy/examples/csi/cephfs/sotrageclass.yamlz
```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: rook-cephfs
# Change "rook-ceph" provisioner prefix to match the operator namespace if needed
provisioner: rook-ceph.cephfs.csi.ceph.com
parameters:
  # clusterID is the namespace where the rook cluster is running
  # If you change this namespace, also change the namespace below where the secret namespaces are defined
  clusterID: rook-ceph

  # CephFS filesystem name into which the volume shall be created
  fsName: myfs

  # Ceph pool into which the volume shall be created
  # Required for provisionVolume: "true"
  pool: myfs-replicated

  # The secrets contain Ceph admin credentials. These are generated automatically by the operator
  # in the same namespace as the cluster.
  csi.storage.k8s.io/provisioner-secret-name: rook-csi-cephfs-provisioner
  csi.storage.k8s.io/provisioner-secret-namespace: rook-ceph
  csi.storage.k8s.io/controller-expand-secret-name: rook-csi-cephfs-provisioner
  csi.storage.k8s.io/controller-expand-secret-namespace: rook-ceph
  csi.storage.k8s.io/node-stage-secret-name: rook-csi-cephfs-node
  csi.storage.k8s.io/node-stage-secret-namespace: rook-ceph

allowVolumeExpansion: true
reclaimPolicy: Delete
```

```
# 해당 파일은 rook/deploy/examples/csi/cephfs/storageclass.yaml에 존재 
# allowVolumeExpansion: true 만 추가 
kubectl create -f storageclass.yaml
```


* pvc 생성 테스트
```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: cephfs-rwx-pvc
spec:
  accessModes:
  - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
  storageClassName: rook-cephfs
```

### rook-ceph clean delete
[공식문서](https://rook.io/docs/rook/v1.9/ceph-teardown.html) 참조해서 삭제 할 것.
