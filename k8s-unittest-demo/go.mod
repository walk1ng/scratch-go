module k8s-unittest-demo

go 1.15

require (
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/gomega v1.10.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.23.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog/v2 v2.8.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7 // indirect
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.0 // indirect

)

replace (
	github.com/Azure/go-autorest v11.2.8+incompatible => github.com/Azure/go-autorest/autorest v0.9.6
	github.com/Azure/go-autorest/autorest v0.9.6 => github.com/Azure/go-autorest v11.2.8+incompatible
	k8s.io/api => k8s.io/api v0.19.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.2
	k8s.io/client-go => k8s.io/client-go v0.19.2 // Required by client-go
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.3
)
