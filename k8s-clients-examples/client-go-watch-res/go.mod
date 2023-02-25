module client-go-watch-res

go 1.15

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/zap v1.19.1
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/tools v0.1.6-0.20210820212750-d4cc65f0b2ff // indirect
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/structured-merge-diff/v4 v4.2.0 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	github.com/Azure/go-autorest v11.2.8+incompatible => github.com/Azure/go-autorest/autorest v0.9.6
	github.com/Azure/go-autorest/autorest v0.9.6 => github.com/Azure/go-autorest v11.2.8+incompatible
	k8s.io/api => k8s.io/api v0.19.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.2
	k8s.io/client-go => k8s.io/client-go v0.19.2 // Required by client-go
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.3
)
