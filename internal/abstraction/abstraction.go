package abstraction

import (
	"context"
	"fmt"

	"dagger.io/dagger"

	"github.com/alvise88/zero-turnaround-cicd-with-dagger/internal/util"
)

// utility
func toPackages(pkgs []string) []struct {
	Name    string
	Version string
} {
	packages := []struct {
		Name    string
		Version string
	}{}

	for _, pkg := range pkgs {
		packages = append(packages, struct {
			Name    string
			Version string
		}{Name: pkg})
	}

	return packages
}

// based on dagger-cue AlpineOpts definition
type AlpineOpts struct {
	Version  string
	Packages []struct {
		Name    string
		Version string
	}
}

// based on dagger-cue AnsibleOpts definition
type AnsibleOps struct {
	Version    string
	Project    *dagger.Directory
	AnsibleCfg string
}

func Alpine(ctx context.Context, client *dagger.Client, opts AlpineOpts) (*dagger.Container, error) {
	alpine := client.Container().From(fmt.Sprintf("alpine:%s", opts.Version)).
		WithExec(util.ToCommand("apk update"))

	for _, pckg := range opts.Packages {
		alpine = alpine.WithExec(util.ToCommand(fmt.Sprintf("apk --no-cache add -U %s", pckg.Name)))
	}

	return alpine, nil
}

var defaultAnsibleConfigLocation = "/etc/ansible"

func Anisble(ctx context.Context, client *dagger.Client, opts AnsibleOps) (*dagger.Container, error) {
	ansible, err := Alpine(ctx, client, AlpineOpts{
		Version: "3.17.1",
		Packages: toPackages([]string{
			"bash",
			"gcc",
			"grep",
			"openssl-dev",
			"py3-pip",
			"python3-dev",
			"gpgme-dev",
			"libc-dev",
			"rust",
			"cargo",
		}),
	})

	if err != nil {
		return nil, err
	}

	ansible = ansible.WithDirectory(defaultAnsibleConfigLocation, opts.Project, dagger.ContainerWithDirectoryOpts{}).
		WithEnvVariable("ANSIBLE_CONFIG_LOCATION", defaultAnsibleConfigLocation).
		WithExec(util.ToCommand("mkdir -p /etc/ansible")).
		WithNewFile(fmt.Sprintf("%s/%s", defaultAnsibleConfigLocation, "ansible.cfg"), dagger.ContainerWithNewFileOpts{
			Contents: opts.AnsibleCfg,
		}).
		WithExec(util.ToCommand(fmt.Sprintf("pip3 install ansible==%s", opts.Version)))

	return ansible, nil
}
