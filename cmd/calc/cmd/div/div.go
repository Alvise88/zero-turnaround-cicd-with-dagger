package div

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

func Div() *cobra.Command {
	divCmd := &cobra.Command{
		Use:   "div first second",
		Short: "division operation",
		Long:  `division operation`,
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

			div, err := calc.Div(first, second)

			if err != nil {
				return err
			}

			_, err = fmt.Printf("%d", div)

			return err
		},
	}

	return divCmd
}
