package postgres

import (
	"fmt"
	"net"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	log "github.com/sirupsen/logrus"
)

var hostName = "localhost"

const defaultPostgresPort = "5432/tcp"

type Container struct {
	resource *dockertest.Resource
}

func NewContainer(config *Config, connectFn func() error) (*Container, error) {
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
			Repository: "postgres",
			Tag:        "13-alpine",
			Env: []string{
				"POSTGRES_USER=" + config.Username,
				"POSTGRES_PASSWORD=" + config.Password,
				"POSTGRES_DB=images",
				"listen_addresses = '*'",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				defaultPostgresPort: {{
					HostIP:   hostName,
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
	addr := fmt.Sprintf("%s:%s", hostName, resource.GetPort(defaultPostgresPort))
	if err := pool.Retry(func() error {
		config.URL = fmt.Sprintf("postgres://%s:%s@%s/images?sslmode=disable", config.Username, config.Password, addr)
		return connectFn()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return container, nil
}

func (c *Container) Purge() error {
	return c.resource.Close()
}

func ApplyMigrate(databaseUrl string) {
	mig, err := migrate.New(
		"file://../../../migration",
		databaseUrl)
	if err != nil {
		log.Errorf("failed to create migrations instance: %v", err)
		return

	}

	if err := mig.Up(); err != nil {
		log.Errorf("could not exec migration: %v", err)
		return
	}
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
