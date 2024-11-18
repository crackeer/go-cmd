package command

import (
	"fmt"

	"github.com/crackeer/go-cmd/util"
	"github.com/spf13/cobra"
)

func NewIP(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
		Run:   showIP,
	}
	cmd.SetHelpTemplate(`./got ip
`)
	return cmd
}

func showIP(cmd *cobra.Command, args []string) {
	ip := util.GetInnerIP()
	for _, v := range ip {
		fmt.Println(v)
	}
}
