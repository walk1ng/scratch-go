module democlient

go 1.15

require (
	k8s.io/api v0.20.0
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v0.20.0
	walk1ng.io/demo v0.0.0
)

replace walk1ng.io/demo => ../demo
