package cmd

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"testing"
)

func TestWorkloadValidate(t *testing.T) {
	peo := NewPExecOptions(genericclioptions.IOStreams{})
	peo.offset = 1
	peo.args = []string{"pexec", "deploy", "nginx", "uname -a"}
	err := peo.Validate()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("pass common scence.")
	}

	// not enough args
	peo.args = []string{"pexec", "deploy", "nginx"}
	err = peo.Validate()
	if err != nil {
		t.Log("pass NotEnoughArgs")
	} else {
		t.Fatal("Failed to Parse not enough args scence.")
	}
}
