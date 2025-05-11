package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +code-gen:validate=true

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Name   string `json:"name"`
	Gender string `json:"gender" v:"required,oneof=male female" export:"true"`

	List []string `json:"list" v:"required,min=1,unique,dive,required" export:"true"`
	Sub  SubFoo   `json:"sub,omitempty"`
}

// +kubebuilder:object:generate=true

type SubFoo struct {
	Labels map[string]string `json:"labels,omitempty"`
}
