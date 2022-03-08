module using-controller-runtime-client

go 1.15

require (
	github.com/kudobuilder/kudo v0.19.0
	github.com/spotahome/redis-operator v1.1.1
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/controller-runtime v0.7.2
)

replace (
	k8s.io/api => k8s.io/api v0.19.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.2
	k8s.io/client-go => k8s.io/client-go v0.19.2 // Required by client-go
)
