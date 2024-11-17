package command

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// NewCSV2JSONCommand
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewNowTime(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   getNowTime,
	}
	cmd.SetHelpTemplate(`./xtool now
`)
	return cmd
}

func getNowTime(cmd *cobra.Command, args []string) {
	nowTime := time.Now()
	fmt.Println("timestamp:", nowTime.Unix())
	fmt.Println("date:", nowTime.Format("2006-01-02 15:04:05"))

}
