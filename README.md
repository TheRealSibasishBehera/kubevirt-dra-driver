# Kubevirt DRA Driver [WIP]

<p align="center">
<img src="https://github.com/kubevirt/community/raw/main/logo/KubeVirt_icon.png" width="100">
</p>

This repository contains an resource driver for use with the [Dynamic
Resource Allocation
(DRA)](https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation/)
feature of Kubernetes in KubeVirt.

It is intended to demonstrate how Host devices can be allocated usign DRA , instead of the existing device-plugin framework . It can be deployed in [KubeVirtCI](https://github.com/kubevirt/kubevirtci) after activating
DRA following [this guide](doc/CI_SETUP.md)  
## About resource driver

Kubevirt DRA Driver is a better alternative to the Device Plugin Framerwork which is used in KubeVirt as a part of the virt-handler component , which enable device provisioning in Virutal Machine Instances(VMIs) in Kubevirt.


### Prerequisites

* [GNU Make 3.81+](https://www.gnu.org/software/make/)
* [GNU Tar 1.34+](https://www.gnu.org/software/tar/)
* [docker v20.10+ (including buildx)](https://docs.docker.com/engine/install/)

### Documentation

- [How to setup KubeVirt CI with DRA enabled](doc/CI_SETUP.md)
- [How to deploy and use KubeVirt DRA resource driver](doc/USAGE.md)
- [How to change and rebuild the driver](doc/BUILD.md)