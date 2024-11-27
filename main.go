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
	rootCmd.AddCommand(command.NewUnzip("unzip", "解压zip文件", ""))
	rootCmd.AddCommand(command.NewZip("zip", "压缩文件", ""))
	rootCmd.AddCommand(command.NewDownload("download", "批量下载文件", ""))
	rootCmd.AddCommand(command.NewJsonExtractCmd("json-extract", "json提取", ""))
	rootCmd.AddCommand(command.NewMySQLQueryCommand("mysql-query", "MySQL查询", ""))
	rootCmd.AddCommand(command.NewMySQLInsertCommand("mysql-insert", "MySQL插入", ""))
	rootCmd.AddCommand(command.NewSSHExecute("ssh-exec", "远程执行", ""))
	rootCmd.AddCommand(command.NewEncryptFile("encrypt-file", "文件内容加密", ""))
	rootCmd.AddCommand(command.NewDecryptFile("decrypt-file", "文件内容解密", ""))
	rootCmd.Execute()
}
