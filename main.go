package main

import (
	"github.com/crackeer/go-cmd/cmd/convert"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	rootCmd.AddCommand(convert.NewCSV2JSONCommand("csv2json [file path]", "解析csv文件到json", ""))

	rootCmd.Execute()
}
