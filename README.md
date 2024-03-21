# K-RAF
A Kubernetes-based Resource Augmentation Framework for Edge Devices

<!-- 
[ACMHPDC'24] K-RAF: A Kubernetes-based Resource Augmentation Framework for Edge Devices
* [What is TS-trie](#What-is-TS-trie)
* [Settings](#Settings)
    *  [Environments](#Environments)
    *  [Public datasets](#Public-datasets)
* [Building from source](#Building-from-source)
    * [Build TS-trie indexing server](#build-ts-trie-indexing-server)
    * [Build front application server](#Buil-front-application-server)
* [Main Features](#Main-Features)
    * [Point query](#Point-query)
    * [Trajectory query](#trajectory-query)
* [How to Use](#How-to-Use)
    * [Experiment input form](#experiment-input-form)
    * [Experiment process](#experiment-process)
* [Sample Results](#Sample-Results)
    * [Point query](#point-query-1)
    * [Trajectory query](#trajectory-query-1) -->

## What is K-RAF

K-RAF assists edge devices to overcome their limited capabilities by provisioning virtualized computation and storage resources in a Kubernetes environment.

<p align="center">
  <img align="center" width="75%" src="figures/K-RAF_archtecture.png"></img>
</p>

K-RAF consists of three modules (Execution Migration Module, Monitoring Module, and Storage Module):
1) Execution Migration Module:
    - Utilizes migration features of K-RAF.
    - Receives IoT device's app execution environment info from the database.
    - Performs migration on the edge server by creating application pods on a Kubernetes cluster.
    - User input for GPU usage guides GPU allocation during migration.
    - Kubernetes' Custom Resource Definition partitions GPU VRAM for GPU-reliant applications.

2) Monitoring Module:
    - Consists of Prometheus, Grafana, and Node-exporter.
    - Node-exporter deployed as a pod on each node of the edge server to collect data.
    - Prometheus stores collected data in a database.
    - Grafana visualizes data from Prometheus via a dashboard.
    - Users can query Prometheus to retrieve specific data.

3) Storage Module:
    - Comprises K-RAF's volume-related functions and rook-ceph.
    - rook-ceph builds a storage pool using HDDs and SSDs in the edge server cluster.
    - SSDs classified as cache and HDDs as backing storage for cache-tiering.
    - Kubernetes creates a PV using storage from the pool.
    - K-RAF uses WebDAV Server pods to mount PV to IoT devices for virtual storage.



## Settings

### Environments


**Kubernetes-based Resource Augmentation Framework**
- [Kubernetes v1.22.9]
- [Docker v20.10.16]
- [rook-ceph v1.9.10]
- [davfs2 v1.6.1]

**Edge Server Cluster**
- [Intel i7-11700]
- [RAM 64 GB]
- [SSD 500 GB]
- [HDD 1 TB]
- [Ubuntu 19.04]

**Edge Device**
- [Raspberry Pi 4]
    - Installed davfs2 as WebDAV client
