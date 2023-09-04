package s3

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

const defaultHostName = "localhost"

type Container struct {
	resource *dockertest.Resource
	addr     string
}

func NewContainer() (*Container, error) {
	hostPort, err := getFreePort()
	if err != nil {
		return nil, fmt.Errorf("could not get free hostPort: %w", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "minio/minio",
			Tag:        "latest",
			Cmd: []string{
				"server",
				"/data",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				"9000/tcp": {{
					HostIP:   defaultHostName,
					HostPort: strconv.Itoa(hostPort),
				}},
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
	if err != nil {
		return nil, fmt.Errorf("could not create a container: %w", err)
	}

	container := &Container{
		resource: resource,
	}
	addr := fmt.Sprintf("%s:%s", defaultHostName, resource.GetPort("9000/tcp"))
	uri := fmt.Sprintf("http://%s/minio/health/live", addr)
	container.addr = addr
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = 5 * time.Second
	if err = pool.Retry(func() error {
		if _, err := retryClient.Get(uri); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return container, nil
}

func (c *Container) Purge() error {
	return c.resource.Close()
}

func (c *Container) GetAddr() string {
	return c.addr
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
