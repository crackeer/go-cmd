package command

import (
	"fmt"
	"net/url"

	"github.com/crackeer/go-cmd/util"
	"github.com/spf13/cobra"
)

var matchType *string = new(string)

func NewJsonExtractCmd(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(2),
		Run:   doJsonExtract,
	}
	matchType = cmd.PersistentFlags().StringP("match", "m", "url", "输出链接")
	cmd.SetHelpTemplate(`./got json_extract --match=url file.json
`)
	return cmd
}

func doJsonExtract(cmd *cobra.Command, args []string) {
	var fileData interface{}
	if err := util.ReadFileAs(args[0], &fileData); err != nil {
		panic(err)
	}
	retData := util.Extract(fileData, isURL)
	for _, value := range retData {
		fmt.Println(value)
	}

}

func isURL(value string) bool {
	_, err := url.Parse(value)
	if err != nil {
		return false
	}
	return true
}
