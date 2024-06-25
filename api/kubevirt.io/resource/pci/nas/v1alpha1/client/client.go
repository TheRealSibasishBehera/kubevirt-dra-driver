package client

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nascrd "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"
	nasclient "kubevirt.io/kubevirt-dra-driver/pkg/kubevirt.io/resource/clientset/versioned/typed/nas/v1alpha1"
)

type Client struct {
	nas    *nascrd.NodeAllocationState
	client nasclient.NasV1alpha1Interface
}

func New(nas *nascrd.NodeAllocationState, client nasclient.NasV1alpha1Interface) *Client {
	return &Client{
		nas,
		client,
	}
}

func (c *Client) GetOrCreate(ctx context.Context) error {
	err := c.Get(ctx)
	if err == nil {
		return nil
	}
	if errors.IsNotFound(err) {
		return c.Create(ctx)
	}
	return err
}

func (c *Client) Create(ctx context.Context) error {
	crd := c.nas.DeepCopy()
	crd, err := c.client.NodeAllocationStates(c.nas.Namespace).Create(ctx, crd, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	*c.nas = *crd
	return nil
}

func (c *Client) Delete(ctx context.Context) error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	err := c.client.NodeAllocationStates(c.nas.Namespace).Delete(ctx, c.nas.Name, deleteOptions)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (c *Client) Update(ctx context.Context, spec *nascrd.NodeAllocationStateSpec) error {
	crd := c.nas.DeepCopy()
	crd.Spec = *spec
	crd, err := c.client.NodeAllocationStates(c.nas.Namespace).Update(ctx, crd, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	*c.nas = *crd
	return nil
}

func (c *Client) UpdateStatus(ctx context.Context, status string) error {
	crd := c.nas.DeepCopy()
	crd.Status = status
	crd, err := c.client.NodeAllocationStates(c.nas.Namespace).Update(ctx, crd, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	*c.nas = *crd
	return nil
}

func (c *Client) Get(ctx context.Context) error {
	crd, err := c.client.NodeAllocationStates(c.nas.Namespace).Get(ctx, c.nas.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	*c.nas = *crd
	return nil
}
