package util

import (
	"errors"
	"io"
	"strings"

	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func Execute(client kubernetes.Interface, namespace *string, config *restclient.Config, ignoreHostname bool, podName string, containerName string, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {

	var cmd []string
	if !ignoreHostname {
		cmd = []string{
			"sh",
			"-c",
			fmt.Sprintf("echo -n \"[%s] \"&&%s", podName, command),
		}
	} else {
		cmd = []string{
			"sh",
			"-c",
			command,
		}
	}

	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(*namespace).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	if containerName != "" {
		option.Container = containerName
	}

	if stdin == nil {
		option.Stdin = false
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})

	if err != nil {
		return err
	}

	return nil
}

func ParseLabels(selectLabels string) (map[string]string, error) {
	labels := make(map[string]string)
	sliceLabels := strings.Split(selectLabels, ",")
	for _, label := range sliceLabels {
		kv := strings.Split(label, "=")
		if len(kv) != 2 {
			return nil, errors.New("invalid labels format")
		}
		labels[kv[0]] = kv[1]
	}
	return labels, nil
}
