package request

import (
	"time"
)

type RequestConfig struct {
	Timeout time.Duration
}

func NewRequestConfig() RequestConfig {
	return RequestConfig{
		Timeout: 10 * time.Second,
	}
}
