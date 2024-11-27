package command

import (
	"fmt"

	"github.com/crackeer/go-cmd/util"
	"github.com/spf13/cobra"
)

var (
	key    *string = new(string)
	output *string = new(string)
)

const AesIV string = "1234567890123456"

// NewEncryptFile
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewEncryptFile(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   doEncryptFile,
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.PersistentFlags().StringVarP(key, "key", "k", "1234567812345678", "key")
	cmd.PersistentFlags().StringVarP(output, "output", "o", "", "output")
	cmd.SetHelpTemplate(`./gf-run` + use + ` path
`)

	return cmd
}

func doEncryptFile(cmd *cobra.Command, args []string) {
	if len(*output) < 1 {
		*output = args[0] + ".enc"
	}
	if err := util.AesEncryptFile(args[0], *output, []byte(*key), AesIV); err != nil {
		panic(err)
	}

	fmt.Println("encrypt success, file:", *output)
}
