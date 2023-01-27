package abstraction

import (
	"context"
	"embed"
	"os"
	"strings"
	"testing"

	"dagger.io/dagger"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/internal/assert"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/internal/util"
)

func TestAlpine(t *testing.T) {
	type testCase struct {
		opts AlpineOpts
	}

	cases := []testCase{
		{
			opts: AlpineOpts{
				Version: "3.17.1",
				Packages: []struct {
					Name    string
					Version string
				}{},
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	for _, tc := range cases {
		alpine, err := Alpine(ctx, client, tc.opts)

		if err != nil {
			t.Error(err)
			return
		}

		alpine = alpine.WithExec(util.ToCommand("cat /etc/alpine-release"))

		stdout, err := alpine.Stdout(ctx)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, strings.Trim(stdout, "\n"), tc.opts.Version)
	}
}

//go:embed data
var data embed.FS

//go:embed data/ansible.cfg
var ansibleCfg string

func TestAnsible(t *testing.T) {
	type testCase struct {
		opts AnsibleOps
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	project := client.Directory()
	project, err = util.CopyEmbedDir(data, project)

	if err != nil {
		t.Error(err)
		return
	}

	cases := []testCase{
		{
			opts: AnsibleOps{
				Version:    "7.1",
				Project:    project,
				AnsibleCfg: ansibleCfg,
			},
		},
	}

	for _, tc := range cases {
		ansible, err := Anisble(ctx, client, tc.opts)

		if err != nil {
			t.Error(err)
			return
		}

		ansible = ansible.WithExec(util.ToCommand("ansible --version"))

		stdout, err := ansible.Stdout(ctx)

		if err != nil {
			t.Error(err)
			return
		}

		assert.StringContains(t, stdout, "ansible python")
	}
}
