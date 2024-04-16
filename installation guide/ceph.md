[Rook Official Document](https://rook.io/docs/rook/v1.9/quickstart.html)

rook version : 1.9.10
ceph version : 16.2.10

## [Important] When k8s is configured as a heterogeneous cluster.
* For example, when a k8s cluster consists of three servers and two Raspberry Pi's.
* This is a way to avoid using the resources of the Raspberry Pi.   

* Clone the CEPH code below is the same, after cloning go to rook/deploy/examples.
  * It will create crds, common, and operator as it is.
```
cd rook/deploy/examples
kubectl create -f crds.yaml -f common.yaml -f operator.yaml
```

* Add the corresponding labels for the nodes that are Raspberry Pi's.
```
kubectl label node <raspberrypi_node_name> rookrole=no-rook-node
```

* Then run cluster.yaml from k8s/rook-ceph in the KSV repository.
```
cd ~/KSV/k8s/
kubectl create -f cluster.yaml
```

* After the cluster configuration and OSD pods are successfully deployed, proceed with the rook ceph installation - ceph toolbox installation.

## Install Luke Ceph 

* Get Ceph code clone
```
git clone --single-branch --branch v1.9.10 https://github.com/rook/rook.git
```

* Configure a Ceph cluster
```
cd rook/deploy/examples
kubectl create -f crds.yaml -f common.yaml -f operator.yaml
kubectl create -f cluster.yaml
```

* Create a Luke operator
```
cd deploy/examples
kubectl create -f crds.yaml -f common.yaml -f operator.yaml

# make sure the rook-ceph-operator is in the `running` state before proceeding
kubectl -n rook-ceph get pod
```

* Install ceph toolbox
[Official Document](https://rook.io/docs/rook/v1.8/ceph-toolbox.html)
```
# Launch the toolbox Pod
kubectl create -f deploy/examples/toolbox.yaml

# Check the toolbox Pod installation status
kubectl -n rook-ceph rollout status deploy/rook-ceph-tools
# Access the rook-ceph-tools Pod when the Pod is finished rolling out
kubectl -n rook-ceph exec -it deploy/rook-ceph-tools -- bash

# You can also access it with the following command
kubectl -n rook-ceph exec -it $(kubectl -n rook-ceph get pod -l "app=rook-ceph-tools" -o jsonpath='{.items[0].metadata.name}') -- bash
```
- You can check the CEPH status by running the following command
    - ceph status
    - ceph osd status
    - ceph df
    - rados df
    - ceph version
    - rook version



* Verify CEPH installation
```
$ kubectl get cephcluster -A
NAMESPACE   NAME        DATADIRHOSTPATH   MONCOUNT   AGE   PHASE   MESSAGE                        HEALTH      EXTERNAL
rook-ceph   rook-ceph   /var/lib/rook     3          70m   Ready   Cluster created successfully   HEALTH_OK
```

**Complete CEPH cluster configuration


## Shared Filesystem
[Rook Official Document](https://rook.io/docs/rook/v1.9/ceph-filesystem.html)

* Create filesystem.yaml, the definition of the shared file system
```
apiVersion: ceph.rook.io/v1
Type: CephFilesystem
Metadata:
  Name: myfs
  namespace: rook-ceph
spec:
  Metadatapool:
    replicated:
      size: 3
  dataPools:
    - name: replicated
      replicated:
        size: 3
  Delete retention file system: true
  MetadataServers:
    activeCount: 1
    activeStandby: true
``` 
* Create a filesystem
  * rook/deploy/examples/filesystem.yaml
```
# Create a filesystem
kubectl create -f filesystem.yaml

# Verify creation
kubectl -n rook-ceph get pod -l app=rook-ceph-mds

NAME                                      READY     STATUS    RESTARTS   AGE
rook-ceph-mds-myfs-7d59fdfcf4-h8kw9       1/1       Running   0          12s
rook-ceph-mds-myfs-7d59fdfcf4-kgkjp       1/1       Running   0          12s
```

* Create storageclass.yaml
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
# that file exists in rook/deploy/examples/csi/cephfs/storageclass.yaml 
# add only allowVolumeExpansion: true 
kubectl create -f storageclass.yaml
```


* Test the PVC creation
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
[Official Document](https://rook.io/docs/rook/v1.9/ceph-teardown.html) to delete.
