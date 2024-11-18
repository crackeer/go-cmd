package command

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	joinString *string
	unique     *bool
)

// NewRegexExtract
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewRegexExtract(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(2),
		Run:   doRegexExtract,
	}
	joinString = cmd.PersistentFlags().StringP("join", "j", "\n", "输出链接")
	unique = cmd.PersistentFlags().BoolP("unique", "u", false, "去重")
	cmd.SetHelpTemplate(`./got regex_extract file.txt pattern
`)

	return cmd
}

func doRegexExtract(cmd *cobra.Command, args []string) {
	bytes, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("读取文件错误:", err)
		return
	}
	text := string(bytes)
	expression := regexp.MustCompile(args[1])
	matches := expression.FindAllString(text, -1)

	if *unique {
		mapData := map[string]interface{}{}
		for i := 0; i < len(matches); i++ {
			mapData[matches[i]] = true
		}
		uniqueData := []string{}
		for k := range mapData {
			uniqueData = append(uniqueData, k)
		}
		fmt.Println(strings.Join(uniqueData, *joinString))
		return
	}
	fmt.Println(strings.Join(matches, *joinString))
}
