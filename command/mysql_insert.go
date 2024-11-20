package command

import (
	"fmt"

	"github.com/crackeer/go-cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewMySQLQueryCommand
//
//	@param use
//	@param short
//	@param long
//	@return *cobra.Command
func NewMySQLInsertCommand(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   doMySQLInsert,
		Args:  cobra.MinimumNArgs(2),
	}
	cmd.SetUsageTemplate(`got mysql-insert table /file/to/data.json
`)

	return cmd
}

func doMySQLInsert(cmd *cobra.Command, args []string) {
	viper.AutomaticEnv()
	database := viper.Get("DATABASE")
	if database == nil {
		fmt.Println("DATABASE is not set")
		return
	}
	db, err := util.GetGormDB(database.(string))
	if err != nil {
		fmt.Println("get gorm db error: ", err)
		return
	}
	list := []map[string]interface{}{}
	if err := util.ReadFileAs(args[1], &list); err != nil {
		fmt.Println("read file error: ", err)
		return
	}
	fmt.Println("start insert...")
	for _, v := range list {
		if err := db.Table(args[0]).Create(&v).Error; err != nil {
			fmt.Println("insert error: ", err)
			return
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println("")
	fmt.Println("insert success")
}
