package scheme

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/suzi1037/pcmd/app/demo/apis/v1"
	"github.com/suzi1037/pcmd/app/demo/apis/v1beta1"
)

var Scheme = runtime.NewScheme()

func init() {
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	// add schema
	utilruntime.Must(v1.AddToScheme(Scheme))
	utilruntime.Must(v1beta1.AddToScheme(Scheme))
}
