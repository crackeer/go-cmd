package command

import (
	"fmt"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

func NewQrcode(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		Run:   showQRCode,
	}
	cmd.SetHelpTemplate(`./got qrcode string
`)
	return cmd
}

func showQRCode(cmd *cobra.Command, args []string) {
	qrcode, err := qrcode.New(args[0], qrcode.Medium)
	if err != nil {
		panic(err)
	}

	// 将二维码转换为ASCII字符
	ascii := qrcode.ToSmallString(false)

	// 输出二维码
	fmt.Println(ascii)
}
