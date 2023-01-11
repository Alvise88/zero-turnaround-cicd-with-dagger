package sum

import (
	"fmt"
	"strconv"

	"github.com/alvise88/zero-turnaround-cicd-with-dagger/pkg/calc"
	"github.com/spf13/cobra"
)

const (
	FIRST  = 0
	SECOND = 1
)

func Sum() *cobra.Command {
	sumCmd := &cobra.Command{
		Use:   "sum first second",
		Short: "addition operation",
		Long:  `addition operation`,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(2)(cmd, args); err != nil {
				return err
			}

			if _, err := strconv.Atoi(args[0]); err != nil {
				return err
			}

			if _, err := strconv.Atoi(args[1]); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			first, _ := strconv.Atoi(args[FIRST])
			second, _ := strconv.Atoi(args[SECOND])

			_, err := fmt.Printf("%d", calc.Sum(first, second))

			return err
		},
	}

	return sumCmd
}
