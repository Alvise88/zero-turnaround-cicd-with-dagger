package mage

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/magefile/mage/mg"
)

var goVersion = "1.19.5"

// define build matrix
var oses = []string{"linux", "darwin"}
var arches = []string{"amd64", "arm64"}

type Calc mg.Namespace

func (calc Calc) Build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	// create empty directory to put build outputs
	outputs := client.Directory()

	// get `golang` image
	golang := client.Container().From(fmt.Sprintf("golang:%s-alpine", goVersion))

	// mount cloned repository into `golang` image
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

	for _, goos := range oses {
		for _, goarch := range arches {
			// create a directory for each os and arch
			path := fmt.Sprintf("out/%s/%s/", goos, goarch)

			// set GOARCH and GOOS in the build environment
			build := golang.WithEnvVariable("GOOS", goos)
			build = build.WithEnvVariable("GOARCH", goarch)

			// build application, (crosscompilation)
			build = build.WithExec([]string{"go", "build", "-o", path, "./cmd/calc"})

			// get reference to build output directory in container
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	// write build artifacts to host
	_, err = outputs.Export(ctx, ".")
	if err != nil {
		return err
	}

	return nil
}

func (calc Calc) Lint(ctx context.Context) error {
	fmt.Println("Linting with Dagger")

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	_, err = client.Container().
		From("golangci/golangci-lint:v1.48").
		WithMountedDirectory("/app", src).
		WithWorkdir("/app").
		WithExec([]string{"golangci-lint", "run", "-v", "--timeout", "5m"}, dagger.ContainerWithExecOpts{}).
		ExitCode(ctx)

	return err
}

func (calc Calc) Test(ctx context.Context) error {
	fmt.Println("Testing with Dagger")

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	// get `golang` image
	golang := client.Container().From(fmt.Sprintf("golang:%s-alpine", goVersion))

	// mount cloned repository into `golang` image
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

	_, err = golang.
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "test", "./..."}, dagger.ContainerWithExecOpts{}).
		ExitCode(ctx)

	return err
}

func platforms() []dagger.Platform {
	platforms := []dagger.Platform{}
	for _, goos := range oses {
		for _, goarch := range arches {
			platform := dagger.Platform(fmt.Sprintf("%s/%s", goos, goarch))

			platforms = append(platforms, platform)
		}
	}

	return platforms
}

func (calc Calc) Publish(ctx context.Context) error {
	fmt.Println("Publishing with Dagger")

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	imageRepo := "alvisevitturi/calc:latest"

	// get reference to the local project
	src := client.Host().Directory(".")

	// get `golang` image
	golang := client.Container().From(fmt.Sprintf("golang:%s-alpine", goVersion))

	// mount cloned repository into `golang` image
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

	platformVariants := make([]*dagger.Container, 0, len(platforms()))

	for _, goos := range oses {
		for _, goarch := range arches {
			platform := dagger.Platform(fmt.Sprintf("%s/%s", goos, goarch))

			// set GOARCH and GOOS in the build environment
			build := golang.WithEnvVariable("GOOS", goos)
			build = build.WithEnvVariable("GOARCH", goarch)

			// build application (crosscompilation)
			build = build.WithExec([]string{"go", "build", "-o", "/output/calc", "./cmd/calc"})

			// select the output directory
			outputDir := build.Directory("/output")

			// wrap the output directory in a new empty container marked
			// with the platform
			calc := client.
				Container(dagger.ContainerOpts{Platform: platform}).
				WithRootfs(outputDir)
			platformVariants = append(platformVariants, calc)
		}
	}

	// publishing the final image uses the same API as single-platform
	// images, but now additionally specify the `PlatformVariants`
	// option with the containers built before.
	imageDigest, err := client.
		Container().
		Publish(ctx, imageRepo, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		panic(err)
	}
	fmt.Println("published multi-platform image with digest", imageDigest)

	return nil
}
