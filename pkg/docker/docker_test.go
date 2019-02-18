package docker_test

import (
	"testing"

	. "github.com/metrue/fx/pkg/docker"
)

func TestNew(t *testing.T) {
	cli := New()
	if cli == nil {
		t.Fatal(cli)
	}

	ok := cli.IsRunning()
	if !ok {
		t.Fatal(ok)
	}

	ref := "metrue/snap"
	_, err := cli.ListImagesWithReference(ref)
	if err != nil {
		t.Fatal(err)
	}
}
