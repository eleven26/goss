package cmd

import (
	"fmt"
	"log"

	"github.com/eleven26/goss/utils"
	"github.com/spf13/cobra"
)

type ListCmdArgs struct {
	Prefix string
}

func parseListCmdArgs(cmd *cobra.Command, args []string) ListCmdArgs {
	var prefix string
	if len(args) > 0 {
		prefix = args[0]
	} else {
		prefix = ""
	}

	return ListCmdArgs{
		Prefix: prefix,
	}
}

var listExamples = Examples(`
		# 列出 test 目录下的文件
		go run main.go list test`)

var listCmd = &cobra.Command{
	Use:     "list <prefix>",
	Short:   "列出指定目录下的文件",
	Long:    "列出指定目录下的文件",
	Example: listExamples,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listCmdArgs := parseListCmdArgs(cmd, args)

		files, err := app.Storage.Files(listCmdArgs.Prefix)
		if err != nil {
			log.Fatalf("list err: %#v", err)
		}

		fmt.Printf("%9s %19s %s\n", "Size", "LastModified", "Key")
		for _, file := range files {
			fmt.Printf("%9s %s %s\n",
				utils.Size(file.Size()),
				file.LastModified().Format("2006-01-02 15:04:05"),
				file.Key())
		}

		// spew.Dump(files)
	},
}
