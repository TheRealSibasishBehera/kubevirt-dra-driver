package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PciClaimParametersSpec is the spec for the PciClaimParameters CRD.
type PciClaimParametersSpec struct {
	//matches the ResourceName from class parameters
	DeviceName string `json:"deviceName"`
	Count      int    `json:"count"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced

// PciClaimParameters holds the set of parameters provided when creating a resource claim for a Pci.
type PciClaimParameters struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PciClaimParametersSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PciClaimParametersList represents the "plural" of a PciClaimParameters CRD object.
type PciClaimParametersList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PciClaimParameters `json:"items"`
}
