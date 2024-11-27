package command

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/crackeer/go-cmd/util"
	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
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
	//router.Static("/", dir)
	router.NoRoute(createStaticHandler(dir))
	errChan := make(chan error, 1)
	go func() {
		errChan <- router.Run(fmt.Sprintf(":%d", httpServerPort))
	}()
	fmt.Println("")
	fmt.Println("--------> Please scan the QR code to visit the server")
	for _, v := range getURLs(httpServerPort) {
		fmt.Println("visit", v)
		showURLQRCode(v)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalChan:
		fmt.Println("exit")
	case err := <-errChan:
		if err != nil {
			fmt.Println(err)
		}
	}
}

var tplExample = pongo2.Must(pongo2.FromString(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Document</title>
</head>
<style type="text/css">
    .directory {
	   color: blue;
    }
	.file {
		text-decoration: none;
		color: #3b3bd9
	}
	.item {
		margin-bottom: 6px;
	}
</style>
<body>
{% for item in list %}
	<div class="item">
		{% if item.is_dir %}
			<a href="/{{ item.path }}" class="directory">{{ item.name }}</a>
		{% else %}
			<a href="/{{ item.path }}" class="file">{{ item.name }}</a>
		{% endif %}
		
	</div>
{% endfor %}
</body>`))

func createStaticHandler(dir string) gin.HandlerFunc {
	fileServer := http.StripPrefix("", http.FileServer(http.Dir(dir)))

	return func(ctx *gin.Context) {
		file := strings.TrimLeft(ctx.Request.URL.Path, "/")

		fullPath := filepath.Join(dir, file)

		if value, err := os.Stat(fullPath); err == nil && !value.IsDir() {
			fileServer.ServeHTTP(ctx.Writer, ctx.Request)
			return
		}

		list := []map[string]interface{}{}
		fileList, err := os.ReadDir(fullPath)
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		sort.Slice(fileList, func(i, j int) bool {
			return fileList[i].IsDir() && !fileList[j].IsDir()
		})
		for _, entry := range fileList {
			item := map[string]interface{}{
				"name":   entry.Name(),
				"is_dir": entry.IsDir(),
				"path":   filepath.Join(file, entry.Name()),
			}
			list = append(list, item)
		}

		err = tplExample.ExecuteWriter(pongo2.Context{"list": list}, ctx.Writer)
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getURLs(port int64) []string {
	innerIPs := util.GetInnerIP()
	retData := []string{}
	for _, v := range innerIPs {
		url := fmt.Sprintf("http://%s:%d", v, port)
		retData = append(retData, url)
	}
	return retData
}

func showURLQRCode(urlString string) {
	qrcode, err := qrcode.New(urlString, qrcode.Low)
	if err == nil {
		ascii := qrcode.ToSmallString(false)
		fmt.Println(ascii)
	}
}
