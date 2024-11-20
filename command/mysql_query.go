package command

import (
	"encoding/json"
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
func NewMySQLQueryCommand(use, short, long string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   doMySQLQuery,
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.SetUsageTemplate(`got mysql-query table  
`)

	return cmd
}

func doMySQLQuery(cmd *cobra.Command, args []string) {
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
	if len(args) >= 2 {
		if err := db.Table(args[0]).Where(args[1]).Find(&list).Error; err != nil {
			fmt.Println("query error: ", err)
			return
		}
	} else {
		if err := db.Table(args[0]).Find(&list).Error; err != nil {
			fmt.Println("query error: ", err)
			return
		}
	}
	bytes, _ := json.Marshal(list)
	fmt.Println(string(bytes))
}
