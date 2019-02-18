package core

import (
	"log"
	"sync"

	"github.com/metrue/fx/pkg/docker"
)

var fxBaseImageRepo = []string{
	"metrue/fx-java-base",
	"metrue/fx-julia-base",
	"metrue/fx-python-base",
	"metrue/fx-node-base",
	"metrue/fx-d-base",
}

// IsFxDockerImageReady check if base docker image is pulled
func IsFxDockerImageReady(cli docker.IClient) bool {
	for _, name := range fxBaseImageRepo {
		images, err := cli.ListImagesWithReference(name)
		if err != nil {
			return false
		}

		if len(images) == 0 {
			return false
		}
	}
	return true
}

// PullBaseIamage pull fx base images
func PullBaseIamage(cli docker.IClient) {
	count := len(fxBaseImageRepo)

	var wg sync.WaitGroup
	wg.Add(count)

	for _, name := range fxBaseImageRepo {
		go func(image string) {
			err := cli.Pull(image)
			if err != nil {
				log.Fatalf("pull image %s met error %v", image, err)
			} else {
				log.Printf("pull image %s ok", image)
			}
			wg.Done()
		}(name)
	}
	wg.Wait()
}
