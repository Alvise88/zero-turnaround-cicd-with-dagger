package docker

import (
	"os"
	"path/filepath"
	"runtime"

	"dagger.io/dagger"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/internal/abstraction/debian"
)

type Opts struct {
	Image *dagger.Container

	Version string
}

func Docker(client *dagger.Client, opts Opts) (*dagger.Container, error) {
	base := opts.Image

	if base == nil {
		debian, err := debian.Debian(client, debian.Opts{
			Packages: []struct {
				Name    string
				Version string
			}{
				{
					Name: "bash",
				},
				{
					Name: "curl",
				},
				{
					Name: "openssh-client",
				},
			},
		})

		if err != nil {
			return nil, err
		}

		base = debian
	}

	dck := advertiseDevEngine(client, base)

	return dck, nil
}

// HostDockerCredentials returns the host's ~/.docker dir if it exists, otherwise just an empty dir
func HostDockerDir(c *dagger.Client) *dagger.Directory {
	if runtime.GOOS != "linux" {
		// doesn't work on darwin, untested on windows
		return c.Directory()
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return c.Directory()
	}
	path := filepath.Join(home, ".docker")
	if _, err := os.Stat(path); err != nil {
		return c.Directory()
	}
	return c.Host().Directory(path)
}

func advertiseDevEngine(c *dagger.Client, ctr *dagger.Container) *dagger.Container {
	// the cli bin is statically linked, can just mount it in anywhere
	dockerCli := c.Container().From("docker:20.10.23-cli").File("/usr/local/bin/docker")

	return ctr.
		// Mount in the docker cli + socket, this will be used to connect to the dev engine
		// container
		WithUnixSocket("/var/run/docker.sock", c.Host().UnixSocket("/var/run/docker.sock")).
		WithMountedFile("/usr/bin/docker", dockerCli)
}
