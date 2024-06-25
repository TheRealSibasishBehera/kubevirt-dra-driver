package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeviceSelector allows one to match on a specific type of Device as part of the class.
type DeviceSelector struct {
	Type              string `json:"type"`
	ResourceName      string `json:"resourceName"`
	PCIVendorSelector string `json:"pciVendorSelector"`
}

// DeviceClassParametersSpec is the spec for the DeviceClassParametersSpec CRD.
type DeviceClassParametersSpec struct {
	DeviceSelector []DeviceSelector `json:"deviceSelector,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Obect
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Cluster

// DeviceClassParameters holds the set of parameters provided when creating a resource class for this driver.
type DeviceClassParameters struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DeviceClassParametersSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeviceClassParametersList represents the "plural" of a DeviceClassParameters CRD object.
type DeviceClassParametersList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DeviceClassParameters `json:"items"`
}