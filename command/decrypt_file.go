package command

import (
	"fmt"

	"github.com/crackeer/go-cmd/util"
	"github.com/spf13/cobra"
)

// NewDecryptFile
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewDecryptFile(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   doDecryptFile,
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.PersistentFlags().StringVarP(key, "key", "k", "1234567812345678", "key")
	cmd.PersistentFlags().StringVarP(output, "output", "o", "", "output")
	cmd.SetHelpTemplate(`./gf-run` + use + ` path
`)

	return cmd
}

func doDecryptFile(cmd *cobra.Command, args []string) {
	if len(*output) < 1 {
		*output = args[0] + ".enc"
	}
	if err := util.AesDecryptFile(args[0], *output, []byte(*key), AesIV); err != nil {
		panic(err)
	}

	fmt.Println("decrypt success, file:", *output)
}
