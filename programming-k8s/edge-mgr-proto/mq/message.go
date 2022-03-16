package mq

type Verb string

const (
	DELETE Verb = "delete"
	Update Verb = "update"
	Patch  Verb = "patch"
	Create Verb = "create"
)

type Message struct {
	Namespace string
	Name      string
	Verb      Verb
}
