package command

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	httpServerPort int64
)

func NewHttpServer(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   doHttpServer,
	}
	cmd.PersistentFlags().Int64VarP(&httpServerPort, "port", "p", 8888, "http server port")
	cmd.SetHelpTemplate(`./got` + use + ` directory
`)
	return cmd
}

func doHttpServer(cmd *cobra.Command, args []string) {
	dir := "./"
	if len(args) > 0 {
		dir = args[0]
	}
	router := gin.Default()
	router.Static("/", dir)
	router.Run(fmt.Sprintf(":%d", httpServerPort))
}
