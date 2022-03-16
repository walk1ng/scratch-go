package mq

import "sync"

type WorkChannel struct {
	sync.RWMutex
	Queue chan Message
}
