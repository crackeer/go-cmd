package convert

import (
	"encoding/json"
	"fmt"

	"github.com/crackeer/go-cmd/util"

	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

var (
	withHeader *bool
)

// NewCSV2JSONCommand
//  @param use
//  @param short
//  @param long
//  @return *cobra.Command
func NewCSV2JSONCommand(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		Run:   csv2json,
	}
	cmd.SetHelpTemplate(`./xtool csv2json your/csv/file/name.csv --header
`)
	withHeader = cmd.Flags().Bool("header", false, "with header")
	return cmd
}

func csv2json(cmd *cobra.Command, args []string) {
	fileName := args[0]

	if withHeader != nil && *withHeader {
		if list, err := util.ReadCSVWithHeader(fileName); err != nil {
			fmt.Println(chalk.Red, err.Error(), chalk.Reset)
		} else {
			bytes, _ := json.Marshal(list)
			fmt.Println(string(bytes))
		}
		return
	}
	if list, err := util.ReadCSV(fileName); err != nil {
		fmt.Println(chalk.Red, err.Error(), chalk.Reset)
	} else {
		bytes, _ := json.Marshal(list)
		fmt.Println(string(bytes))
	}
}
