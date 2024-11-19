package command

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/crackeer/go-cmd/util"
	"github.com/spf13/cobra"
)

// NewDownload
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewDownload(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   doDownload,
		Args:  cobra.MinimumNArgs(2),
	}
	cmd.SetHelpTemplate(`./got ` + use + ` url_text.txt dir
`)
	return cmd
}

func doDownload(cmd *cobra.Command, args []string) {
	bytes, err := os.ReadFile(args[0])
	if err != nil {
		panic(err)
	}
	urls := strings.Split(string(bytes), "\n")

	for _, urlString := range urls {
		if urlString == "" {
			continue
		}
		fmt.Println("")
		fmt.Println("Downloading: ", urlString)
		urlObject, err := url.Parse(urlString)
		if err != nil {
			fmt.Println("url parse error: ", err)
			continue
		}
		target := filepath.Join(args[1], urlObject.Path)

		if err := util.DownloadTo(urlString, target); err != nil {
			fmt.Println("download error: ", err)
			continue
		}
	}
}
