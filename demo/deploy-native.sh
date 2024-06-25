#!/bin/bash

# Apply CRDs

kubectl apply -f ../deployments/native/kubevirt-dra-driver/crds/nas.pci.resource.kubevirt.io_nodeallocationstates.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/crds/pci.resource.kubevirt.io_pciclaimparameters.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/crds/pci.resource.kubevirt.io_deviceclassparameters.yaml


# Apply other Kubernetes objects
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/namespace.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/serviceaccount.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/clusterrole.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/clusterrolebinding.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/resourceclass.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/controller.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/kubeletplugin.yaml