package main

import (
	"github.com/crackeer/go-cmd/command"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	rootCmd.AddCommand(command.NewCSV2JSONCommand("csv2json [file path]", "解析csv文件到json", ""))
	rootCmd.AddCommand(command.NewNowTime("time", "current time", ""))
	rootCmd.AddCommand(command.NewGolangMD5("go-md5", "计算md5", ""))
	rootCmd.AddCommand(command.NewGolangBase64Decode("go-base64-decode", "base64解码", ""))
	rootCmd.AddCommand(command.NewGolangBase64Encode("go-base64-encode", "base64编码", ""))
	rootCmd.AddCommand(command.NewHttpPost("http-post", "HttpPost请求", ""))

	rootCmd.Execute()
}
