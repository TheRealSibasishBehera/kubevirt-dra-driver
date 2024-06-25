# Enable Support for Dynamic Resource Allocation in KubeVirtCI

This guide will walk you through the steps to enable support for dynamic resource allocation in KubeVirtCI.

## Step 1: Start KubeVirt CI

1. **Download the KubeVirtCI repository:**

    ```bash
    git clone https://github.com/kubevirt/kubevirtci.git
    cd kubevirtci
    ```

2. **Start a multi-node Kubernetes cluster with 2 NICs:**

    ```bash
    export KUBEVIRT_PROVIDER=k8s-1.30
    export KUBEVIRT_NUM_NODES=1
    export KUBEVIRT_NUM_SECONDARY_NICS=1
    make cluster-up
    ```

3. **Set the KUBEVIRTCI_TAG environment variable:**

    ```bash
    export KUBEVIRTCI_TAG=$(curl -L -Ss https://storage.googleapis.com/kubevirt-prow/release/kubevirt/kubevirtci/latest)
    ```

4. **Verify that the nodes are up and running:**

    ```bash
    cluster-up/kubectl.sh get nodes
    ```

## Step 2: Enable Dynamic Resource Allocation

1. **SSH into the control plane node:**

    ```bash
    cluster-up/ssh.sh node01
    ```

2. **Edit the manifest files to add the necessary parameters:**

    - **kube-apiserver:** Edit `/etc/kubernetes/manifests/kube-apiserver.yaml` and add the following to the `command` section under the `container` section:

      ```yaml
      - --feature-gates=DynamicResourceAllocation=true
      - --runtime-config=resource.k8s.io/v1alpha2=true
      ```

    - **kube-controller-manager:** Edit `/etc/kubernetes/manifests/kube-controller-manager.yaml` and add the following to the `command` section under the `container` section:

      ```yaml
      - --feature-gates=DynamicResourceAllocation=true
      ```

    - **kube-scheduler:** Edit `/etc/kubernetes/manifests/kube-scheduler.yaml` and add the following to the `command` section under the `container` section:

      ```yaml
      - --feature-gates=DynamicResourceAllocation=true
      ```

    - **kubelet config:** Edit `/var/lib/kubelet/config.yaml` and add the following:

      ```yaml
      featureGates:
        DynamicResourceAllocation: true
      ```

## Step 3: Verification

1. **Check if your Kubernetes cluster supports dynamic resource allocation:**

    ```bash
    cluster-up/kubectl.sh get resourceclasses
    ```

    - If your cluster supports dynamic resource allocation, the response is either a list of `ResourceClass` objects or:

      ```bash
      No resources found
      ```

    - If not supported, this error is printed instead:

      ```bash
      error: the server doesn't have a resource type "resourceclasses"
      ```

By following these steps, you should have successfully enabled dynamic resource allocation in KubeVirtCI.
