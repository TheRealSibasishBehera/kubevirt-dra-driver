# Kubevirt DRA Driver Deployment

This guide provides steps to deploy the Kubevirt DRA Driver on your Kubernetes cluster.

## Prerequisites

- A running KubeVirtCI  following [this guide](doc/CI_SETUP.md)
- `kubectl` installed and configured to interact with your cluster

## Deployment Steps
1. Set the `KUBECONFIG` environment variable:

```bash
export KUBECONFIG=$(/path/to/your/ci/directory/kubeconfig.sh)
```

2. Apply the Custom Resource Definitions (CRDs):

```bash
kubectl apply -f ../deployments/native/kubevirt-dra-driver/crds/nas.pci.resource.kubevirt.io_nodeallocationstates.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/crds/pci.resource.kubevirt.io_pciclaimparameters.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/crds/pci.resource.kubevirt.io_pciclassparameters.yaml
```
3. Apply other Kubernetes objects:

```bash
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/serviceaccount.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/clusterrole.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/clusterrolebinding.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/resourceclass.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/controller.yaml
kubectl apply -f ../deployments/native/kubevirt-dra-driver/templates/kubeletplugin.yaml
```

4. Apply the demo Pod:

```bash
kubectl apply -f ../demo/pci-test1.yaml
```

5. Verify that the Pod is running:

```bash
kubectl get pods -A
```

6. Verify that the Pod has the PCI device assigned:

```bash
kubectl exec -it <test-pod> -- lspci
```bash
