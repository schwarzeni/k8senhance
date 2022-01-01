package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudEdgeService is a specification for a CloudEdgeService resource
type CloudEdgeService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CloudEdgeServiceSpec `json:"spec"`
}

// CloudEdgeServiceSpec is the spec for a CloudEdgeService resource
type CloudEdgeServiceSpec struct {
	Selector map[string]string `json:"selector"`
	Port     *int16            `json:"port"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudEdgeServiceList is a list of CloudEdgeService resources
type CloudEdgeServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CloudEdgeService `json:"items"`
}
