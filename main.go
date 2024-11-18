package main

import (
	"github.com/crackeer/go-cmd/command"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	rootCmd.AddCommand(command.NewCSV2JSONCommand("csv2json [file path]", "解析csv文件到json", ""))
	rootCmd.AddCommand(command.NewNowTime("ctime", "获取当前 time", ""))
	rootCmd.AddCommand(command.NewConvertTime("ts", "转换时间戳", ""))
	rootCmd.AddCommand(command.NewGolangMD5("go-md5", "计算md5", ""))
	rootCmd.AddCommand(command.NewGolangBase64Decode("base64-decode", "base64解码", ""))
	rootCmd.AddCommand(command.NewGolangBase64Encode("base64-encode", "base64编码", ""))
	rootCmd.AddCommand(command.NewHttpPost("http-post", "HttpPost请求", ""))
	rootCmd.AddCommand(command.NewQrcode("qrcode", "生成二维码", ""))
	rootCmd.AddCommand(command.NewRegexExtract("regex-extract", "正则提取", ""))
	rootCmd.AddCommand(command.NewHttpServer("http-server", "简单的httpServer", ""))
	rootCmd.AddCommand(command.NewIP("ip", "获取内网ip", ""))
	rootCmd.Execute()
}
