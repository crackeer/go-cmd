package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/crackeer/go-cmd/util"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// NewHttpPost
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewHttpPost(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(2),
		Run:   doHttpPost,
	}
	cmd.SetHelpTemplate(`./gf-run` + use + ` url /file/to/body.json /file/to/header.json
`)
	return cmd
}

func doHttpPost(cmd *cobra.Command, args []string) {
	valueURL := args[0]
	valueBody := args[1]

	client := resty.New()
	request := client.R().EnableTrace()
	headers := map[string]string{}
	if len(args) >= 3 {
		if err := util.ReadFileAs(args[2], &headers); err != nil {
			fmt.Println("读取header文件错误:", err)
			return
		}
	}
	headers["Content-Type"] = "application/json"
	request = request.SetHeaders(headers)

	bytes, err := os.ReadFile(valueBody)
	if err != nil {
		fmt.Println("读取body文件错误:", err)
		return
	}

	var bodyData interface{}
	if err := json.Unmarshal(bytes, &bodyData); err != nil {
		fmt.Println("解析body文件错误:", err)
		return
	}

	response, err := request.SetBody(bodyData).Post(valueURL)
	if err != nil {
		fmt.Println("请求错误:", err)
		return
	}

	result := response.Body()

	fmt.Println("response:")
	fmt.Println(string(result))
}
