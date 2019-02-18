package docker

import (
	"bufio"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	client "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jhoonb/archivex"

	"context"
	"os"
	"time"
)

// IClient docker client interface
type IClient interface {
	ListImagesWithReference(ref string) ([]string, error)
	Pull(image string) error
}

// Client client of docker
type Client struct {
	*client.Client
}

type dockerInfo struct {
	Stream string `json:"stream"`
}

// New a docker client
func New() *Client {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return &Client{cli}
}

// IsRunning if docker is working
func (d *Client) IsRunning() bool {
	ctx := context.Background()
	_, err := d.Info(ctx)
	return err == nil
}

// ListImagesWithReference check if host has images with names
func (d *Client) ListImagesWithReference(ref string) ([]string, error) {
	filter := filters.NewArgs()
	filter.Add("reference", ref)
	list, err := d.ImageList(
		context.Background(),
		types.ImageListOptions{Filters: filter},
	)
	if err != nil {
		return []string{}, err
	}

	images := []string{}
	for _, image := range list {
		// RepoTags is like this metrue/fx-rust-base:latest
		images = append(images, image.RepoTags...)
	}
	return images, nil
}

// Build builds a docker image from the image directory
func (d *Client) Build(name string, dir string) error {
	tar := new(archivex.TarFile)
	err := tar.Create(dir)
	if err != nil {
		return err
	}
	err = tar.AddAll(dir, false)
	if err != nil {
		return err
	}
	err = tar.Close()
	if err != nil {
		return err
	}

	dockerBuildContext, buildContextErr := os.Open(dir + ".tar")
	if buildContextErr != nil {
		return buildContextErr
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{name},
		Labels:     map[string]string{"belong-to": "fx"},
	}
	buildResponse, buildErr := d.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if buildErr != nil {
		return buildErr
	}
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)
	for scanner.Scan() {
		var info dockerInfo
		err := json.Unmarshal(scanner.Bytes(), &info)
		if err != nil {
			return err
		}
	}

	return nil
}

// Pull image from hub.docker.com
func (d *Client) Pull(name string) error {
	ctx := context.Background()
	_, err := d.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Deploy spins up a new container
func (d *Client) Deploy(name string, dir string, port string) (*container.ContainerCreateCreatedBody, error) {
	ctx := context.Background()
	imageName := name
	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"3000/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		},
	}

	resp, err := d.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		return nil, err
	}

	if err = d.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	return &resp, err
}

// Stop interrupts a running container
func (d *Client) Stop(containerID string) (err error) {
	timeout := time.Duration(1) * time.Second
	err = d.ContainerStop(context.Background(), containerID, &timeout)
	return err
}

// Remove interrupts and remove a running container
func (d *Client) Remove(containerID string) (err error) {
	return d.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
}

// RemoveImage remove docker image by imageID
func (d *Client) RemoveImage(imageID string) error {
	_, err := d.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true})
	return err
}
