package Storage

import "time"

type ClientStorage struct {
	clients []Client
}

type Client struct {
	ip         string
	port       string
	lastActive time.Time
	id         string
}
