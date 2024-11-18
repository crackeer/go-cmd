package command

import (
	"fmt"
	"strconv"
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
	cmd.SetHelpTemplate(`./got now
`)
	return cmd
}

func getNowTime(cmd *cobra.Command, args []string) {
	nowTime := time.Now()
	fmt.Println("timestamp:", nowTime.Unix())
	fmt.Println("date:", nowTime.Format("2006-01-02 15:04:05"))
}

func NewConvertTime(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		Run:   convertTime,
	}
	cmd.SetHelpTemplate(`./got timestamp
`)
	return cmd
}

func convertTime(cmd *cobra.Command, args []string) {
	timestamp := args[0]
	value, _ := strconv.ParseInt(timestamp, 10, 64)
	fmt.Println(time.Unix(value, 0).Format("2006-01-02 15:04:05"))
}
