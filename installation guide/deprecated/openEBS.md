## OpenEBS 설치
[공식문서](https://openebs.io/docs/user-guides/jiva/jiva-install) 

* OpenEBS operator 설치
```
kubectl apply -f https://openebs.github.io/charts/openebs-operator.yaml
```


* Jiva operator 설치
```
kubectl apply -f  https://openebs.github.io/charts/hostpath-operator.yaml
kubectl apply -f https://openebs.github.io/charts/jiva-operator.yaml
```

* Jiva volume policy 생성 및 설치 

```
# jiva_storage_class.yaml
apiVersion: openebs.io/v1alpha1
kind: JivaVolumePolicy
metadata:
  name: example-jivavolumepolicy
  namespace: openebs
spec:
  replicaSC: openebs-hostpath
  target:
    # This sets the number of replicas for high-availability
    # replication factor <= no. of (CSI) nodes
    replicationFactor: 3
    # disableMonitor: false
    # auxResources:
    # tolerations:
    # resources:
    # affinity:
    # nodeSelector:
    # priorityClassName:
  # replica:
    # tolerations:
    # resources:
    # affinity:
    # nodeSelector:
    # priorityClassName:
```

```
kubectl apply -f jiva_storage_class.yaml
```

* storage class 생성 및 설치 
```
# jiva_storage_class.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: openebs-jiva-csi-sc
provisioner: jiva.csi.openebs.io
allowVolumeExpansion: true
parameters:
  cas-type: "jiva"
  policy: "example-jivavolumepolicy"

```
```
kubectl apply -f jiva_storage_class.yaml
```
