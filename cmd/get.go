package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type GetCmdArgs struct {
	Key    string
	Target string
}

func parseGetCmdArgs(cmd *cobra.Command, args []string) GetCmdArgs {
	key := args[0]

	target, _ := cmd.Flags().GetString("output")
	if target == "" {
		target = filepath.Base(key)
	}

	return GetCmdArgs{
		Key:    key,
		Target: target,
	}
}

var getExamples = Examples(`
		# 获取 key 为 test.txt 的文件，保存到当前文件夹的 test.txt 文件中
		go run main.go get test.txt`)

var getCmd = &cobra.Command{
	Use:     "get <key>",
	Short:   "获取指定文件",
	Long:    "获取指定文件，保存到当前目录",
	Example: getExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getCmdArgs := parseGetCmdArgs(cmd, args)

		err := app.Storage.GetToFile(getCmdArgs.Key, getCmdArgs.Target)
		if err != nil {
			log.Fatal(err)
		}

		color.Green(fmt.Sprintf("下载成功！保存路径：\"%s\"", getCmdArgs.Target))
	},
}

func init() {
	getCmd.PersistentFlags().StringP("output", "o", "", "保存到本地的路径")
}
