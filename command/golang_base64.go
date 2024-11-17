package command

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

// NewGolangBase64Decode    ...
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewGolangBase64Decode(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		Run:   doBase64Decode,
	}
	cmd.SetHelpTemplate(`./gf-run` + use + ` string
`)
	return cmd
}

// doBase64Decode ...
//
//	@param cmd
//	@param args
func doBase64Decode(cmd *cobra.Command, args []string) {
	value := args[0]
	decodedBytes, err := base64.StdEncoding.DecodeString(value)
	fmt.Println("value:", value)
	if err != nil {
		fmt.Println("解码错误:", err)
		return
	}
	fmt.Println("decoded：", string(decodedBytes))
}

// NewGolangBase64Encode    ...
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewGolangBase64Encode(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		Run:   doBase64Encode,
	}
	cmd.SetHelpTemplate(`./gf-run` + use + ` string
`)
	return cmd
}

// doBase64Decode ...
//
//	@param cmd
//	@param args
func doBase64Encode(cmd *cobra.Command, args []string) {
	value := args[0]
	encodedString := base64.StdEncoding.EncodeToString([]byte(value))
	fmt.Println("value:", value)
	fmt.Println("encoded", encodedString)
}
