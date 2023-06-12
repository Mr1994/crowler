package cmd

import (
	"api_client/modeltools/dbtools"
	"api_client/modeltools/generate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createModelCommand)
	createModelCommand.Flags().String("table_name", "", "A help for foo")
}

var createModelCommand = &cobra.Command{

	Use:   "createModel",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		tableName, _ := cmd.Flags().GetString("table_name")
		//初始化数据库
		dbtools.Init()
		//generate.Genertate() //生成所有表信息
		generate.Genertate(tableName) //生成指定表信息，可变参数可传入多个表名
	},
}
