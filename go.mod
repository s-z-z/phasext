module github.com/s-z-z/phasext

go 1.23.0

toolchain go1.23.6

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/fatih/color v1.15.0
	github.com/go-ini/ini v1.67.0
	github.com/go-playground/validator/v10 v10.23.0
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/olekukonko/tablewriter v1.0.7
	github.com/pkg/errors v0.9.1
	github.com/samber/lo v1.47.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.12.0
	golang.org/x/sync v0.11.0
	google.golang.org/grpc v1.65.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/apimachinery v0.31.2
	k8s.io/klog/v2 v2.130.1
	k8s.io/kubernetes v1.31.2
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/olekukonko/errors v0.0.0-20250405072817-4e6d85265da6 // indirect
	github.com/olekukonko/ll v0.0.8 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.31.2 // indirect
	k8s.io/client-go v0.31.2 // indirect
	k8s.io/cluster-bootstrap v0.0.0 // indirect
	k8s.io/component-base v0.31.2 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/kube-proxy v0.0.0 // indirect
	k8s.io/kubelet v0.0.0 // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.31.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.31.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.31.2
	k8s.io/apiserver => k8s.io/apiserver v0.31.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.31.2
	k8s.io/client-go => k8s.io/client-go v0.31.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.31.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.31.2
	k8s.io/code-generator => k8s.io/code-generator v0.31.2
	k8s.io/component-base => k8s.io/component-base v0.31.2
	k8s.io/component-helpers => k8s.io/component-helpers v0.31.2
	k8s.io/controller-manager => k8s.io/controller-manager v0.31.2
	k8s.io/cri-api => k8s.io/cri-api v0.31.2
	k8s.io/cri-client => k8s.io/cri-client v0.31.2
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.31.2
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.31.2
	k8s.io/endpointslice => k8s.io/endpointslice v0.31.2
	k8s.io/kms => k8s.io/kms v0.31.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.31.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.31.2
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.31.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.31.2
	k8s.io/kubectl => k8s.io/kubectl v0.31.2
	k8s.io/kubelet => k8s.io/kubelet v0.31.2
	k8s.io/metrics => k8s.io/metrics v0.31.2
	k8s.io/mount-utils => k8s.io/mount-utils v0.31.2
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.31.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.31.2
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.31.2
	k8s.io/sample-controller => k8s.io/sample-controller v0.31.2
)
