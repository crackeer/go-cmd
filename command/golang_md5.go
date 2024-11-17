package command

import (
	"fmt"

	"github.com/crackeer/go-cmd/util"

	"github.com/spf13/cobra"
)

// NewGolangMD5 NewUpdateMySQLValuet   ...
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewGolangMD5(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		Run:   doGolangMD5,
	}
	cmd.SetHelpTemplate(`./gf-run` + use + ` string
`)
	return cmd
}

// doDeleteMySQL doUpdateKeyUrlHost
//
//	@param cmd
//	@param args
func doGolangMD5(cmd *cobra.Command, args []string) {
	value := args[0]
	md5Value := util.MD5(value)
	fmt.Println("value:", value)
	fmt.Println("md5", md5Value)
}
