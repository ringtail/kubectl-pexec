package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"fmt"
	"errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
	"github.com/ringtail/kubectl-pexec/pkg/util"
	"k8s.io/client-go/rest"
	"strings"
	"os"
	"path"
	"time"
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

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

type PExecOptions struct {
	configFlags  *genericclioptions.ConfigFlags
	restConfig   *rest.Config
	args         []string
	workloadType string
	offset       int
	genericclioptions.IOStreams
}

func (peo *PExecOptions) Complete(c *cobra.Command, args []string) (err error) {
	peo.args = args
	for index, _ := range args {
		if args[index] == "pexec" {
			peo.offset = index + 1
		}
	}
	return nil
}

func (peo *PExecOptions) Validate() (err error) {
	args := peo.args

	if len(args) < 3+peo.offset {
		return errors.New("NoneValidArgs")
	}

	workloadType := args[0+peo.offset]

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

	kubeconf := *peo.configFlags.KubeConfig

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	udc := path.Join(home, ".kube", "config")

	if _, err := os.Stat(udc); kubeconf == "" && err == nil {
		kubeconf = udc
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconf)
	if err != nil {
		panic(err)
	}

	peo.restConfig = config

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespace := peo.configFlags.Namespace
	if *namespace == "" {
		*namespace = "default"
	}

	pods, err := peo.GetPods(clientSet, namespace)

	if err != nil {
		return err
	}

	err = peo.Pexec(clientSet, namespace, pods)
	if err != nil {
		return err
	}
	return nil
}

func (peo *PExecOptions) GetPods(clientSet *kubernetes.Clientset, namespace *string) (pods []v1.Pod, err error) {
	workloadName := peo.args[1+peo.offset]

	matchLabels := make(map[string]string)
	switch peo.workloadType {
	case "Deployment":
		deploy, err := clientSet.AppsV1().Deployments(*namespace).Get(workloadName, metav1.GetOptions{})
		if err != nil {
			// handle error
			return pods, err
		}
		matchLabels = deploy.GetLabels()
	case "StatefulSet":
		statefulSet, err := clientSet.AppsV1().StatefulSets(*namespace).Get(workloadName, metav1.GetOptions{})
		if err != nil {
			// handle error
			return pods, err
		}
		matchLabels = statefulSet.GetLabels()
	case "DaemonSet":
		daemonSet, err := clientSet.AppsV1().DaemonSets(*namespace).Get(workloadName, metav1.GetOptions{})
		if err != nil {
			// handle error
			return pods, err
		}
		matchLabels = daemonSet.GetLabels()
	default:
		//
		return pods, errors.New("UnknownWorkloadType")
	}

	podList, err := clientSet.CoreV1().Pods(*namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(matchLabels)).String(),
	})

	if err != nil {
		// handle error
		return pods, err
	}

	return podList.Items, nil
}

func (peo *PExecOptions) Pexec(clientSet *kubernetes.Clientset, namespace *string, pods []v1.Pod) (err error) {
	now := time.Now()
	wg := &sync.WaitGroup{}
	for index, _ := range pods {
		wg.Add(1)
		go func(pod *v1.Pod, clientSet *kubernetes.Clientset, wg *sync.WaitGroup) {
			if err != nil {
				panic(err)
			}
			util.Execute(clientSet, namespace, peo.restConfig, pod.Name, strings.Join(peo.args[2+peo.offset:], " "), peo.IOStreams.In, peo.IOStreams.Out, peo.IOStreams.ErrOut)
			wg.Done()
		}(&pods[index], clientSet, wg)
	}
	wg.Wait()
	summary := fmt.Sprintf("All pods execution done in %.03fs\n", time.Now().Sub(now).Seconds())
	fmt.Printf("%c[1;0;32m%s%c[0m\n\n", 0x1B, summary, 0x1B)
	return nil
}

func NewPExecOptions(streams genericclioptions.IOStreams) *PExecOptions {
	return &PExecOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
		offset:      0,
	}
}
