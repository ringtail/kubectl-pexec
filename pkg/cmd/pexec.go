package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"fmt"
	"errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	pexecHelp = `
	# Do batch execution in all pods of workloads 
	%s pexec deployment nginx cat /etc/nginx/nginx.conf
`
)

func NewPExecCommand(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewPExecOptions(streams)

	cmd := &cobra.Command{
		Use:          "pexec [deployment(deploy)/daemonset(ds)/statefulset(ss)] [command]",
		Short:        "Do batch execution in all pods of workloads",
		Example:      fmt.Sprintf(pexecHelp, "kubectl"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

type PExecOptions struct {
	configFlags  *genericclioptions.ConfigFlags
	args         []string
	workloadType string
	genericclioptions.IOStreams
}

func (peo *PExecOptions) Complete(c *cobra.Command, args []string) (err error) {
	peo.args = args
	return nil
}

func (peo *PExecOptions) Validate() (err error) {
	args := peo.args

	if len(args) == 0 {
		return errors.New("NoneValidArgs")
	}

	workloadType := args[0]

	switch workloadType {
	case "deployment", "deploy":
		// change workloadType to Deployment
		peo.workloadType = "Deployment"
	case "statefulset", "ss":
		// change workloadType to statefulSet
		peo.workloadType = "StatefulSet"
	case "daemonset", "ds":
		// change workloadType to DaemonSet
		peo.workloadType = "DaemonSet"
	default:
		return errors.New("InvalidWorkloadType")
	}
	return nil
}

func (peo *PExecOptions) Run() (err error) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err)
	}

	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return nil
}

func NewPExecOptions(streams genericclioptions.IOStreams) *PExecOptions {
	return &PExecOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}
