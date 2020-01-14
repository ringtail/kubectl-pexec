package main

import (
	"k8s.io/api/core/v1"
	"context"
	"time"
)

// exec in multi pods and return the msg of total result
func Pexec(pods []v1.Pod, container string, command string, timeout int64) (e error, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))

	cancel()
	return
}
