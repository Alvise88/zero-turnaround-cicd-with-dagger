package cmd

import (
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/cmd/calc/cmd/div"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/cmd/calc/cmd/mul"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/cmd/calc/cmd/pow"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/cmd/calc/cmd/sub"
	"github.com/alvise88/zero-turnaround-cicd-with-dagger/cmd/calc/cmd/sum"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "calc",
		Short: "compute operations",
		Long:  `A Fast and Flexible calculator built with love by MailUp`,
		PreRun: func(cmd *cobra.Command, args []string) {

		},
	}

	rootCmd.AddCommand(sum.Sum())
	rootCmd.AddCommand(sub.Sub())
	rootCmd.AddCommand(mul.Mul())
	rootCmd.AddCommand(div.Div())
	rootCmd.AddCommand(pow.Pow())

	return rootCmd
}
