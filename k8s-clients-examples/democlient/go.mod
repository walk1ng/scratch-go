module democlient

go 1.15

require (
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	walk1ng.io/demo v0.0.0
)

replace walk1ng.io/demo => ../demo
