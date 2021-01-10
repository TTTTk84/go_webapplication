package main

import (
	"time"
)

// room との通信のため
type message struct {
	Name string
	Message string
	When time.Time
	AvatarURL string
}
