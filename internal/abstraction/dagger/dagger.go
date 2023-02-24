package dagger

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/internal/abstraction/debian"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/internal/abstraction/docker"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/internal/util"
)

const (
	daggerRepo = "https://github.com/dagger/dagger.git"

	EngineContainerName = "dagger-engine.ci"
)

type Opts struct {
	Image *dagger.Container

	Version string
}

// it supports only tags
func Dagger(client *dagger.Client, opts Opts) (*dagger.Container, error) {
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

	version := opts.Version

	if version == "" {
		version = "0.3.12"
	}

	cliBinPath := "/usr/local/bin/dagger"

	dck, err := docker.Docker(client, docker.Opts{
		Image:   base,
		Version: "20.10.23",
	})

	if err != nil {
		return nil, err
	}

	dag := dck.WithMountedFile(cliBinPath, Binary(client, version)).
		// Point the SDKs to use the dev engine via these env vars
		WithEnvVariable("_EXPERIMENTAL_DAGGER_CLI_BIN", cliBinPath).
		WithEnvVariable("_EXPERIMENTAL_DAGGER_RUNNER_HOST", "docker-container://"+EngineContainerName)

	return dag, nil
}

// Binary returns a compiled dagger binary
func Binary(c *dagger.Client, version string) *dagger.File {
	return PlatformBinary(c, "", "", "", version)
}

func PlatformBinary(c *dagger.Client, goos, goarch, goarm, version string) *dagger.File {
	base := daggerGoBase(c, version)
	if goos != "" {
		base = base.WithEnvVariable("GOOS", goos)
	}
	if goarch != "" {
		base = base.WithEnvVariable("GOARCH", goarch)
	}
	if goarm != "" {
		base = base.WithEnvVariable("GOARM", goarm)
	}
	return base.
		WithExec([]string{"go", "build", "-o", "./bin/dagger", "-ldflags", "-s -w", "./cmd/dagger"}).
		File("./bin/dagger")
}

func daggerGoCodeOnly(c *dagger.Client, version string) *dagger.Directory {
	daggerBranch := fmt.Sprintf("v%s", version)
	daggerRepo := c.Git(daggerRepo, dagger.GitOpts{KeepGitDir: true}).Branch(daggerBranch).Tree()

	return daggerRepo
}

func daggerGoBase(c *dagger.Client, version string) *dagger.Container {
	repo := daggerGoCodeOnly(c, version)

	// Create a directory containing only `go.{mod,sum}` files.
	goMods := c.Directory()
	for _, f := range []string{"go.mod", "go.sum", "sdk/go/go.mod", "sdk/go/go.sum"} {
		goMods = goMods.WithFile(f, repo.File(f))
	}

	return c.Container().
		From("golang:1.20.0-alpine").
		// gcc is needed to run go test -race https://github.com/golang/go/issues/9918 (???)
		WithExec(util.ToCommand("apk add build-base")).
		WithEnvVariable("CGO_ENABLED", "0").
		// adding the git CLI to inject vcs info
		// into the go binaries
		WithExec([]string{"apk", "add", "git"}).
		WithWorkdir("/app").
		// run `go mod download` with only go.mod files (re-run only if mod files have changed)
		WithMountedDirectory("/app", goMods).
		WithExec([]string{"go", "mod", "download"}).
		// run `go build` with all source
		WithMountedDirectory("/app", repo)
	// // include a cache for go build
	// WithMountedCache("/root/.cache/go-build", c.CacheVolume("go-build"))
}
