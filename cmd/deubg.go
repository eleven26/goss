package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var debugExamples = Examples(`
		# 调试 get 命令，参数是 a.txt
		go run main.go debug Get a.txt`)

var debugCmd = &cobra.Command{
	Use:     "debug <cmd> <...args>",
	Short:   "调试命令",
	Long:    "调试指定命令",
	Example: debugExamples,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		debugArgs := args[1:]
		inputs := make([]reflect.Value, len(debugArgs))

		for i := range debugArgs {
			inputs[i] = reflect.ValueOf(debugArgs[i])
		}

		methodName := title(args[0])
		ref := reflect.ValueOf(app.Storage.Storage())
		method := ref.MethodByName(methodName)
		if !method.IsValid() {
			color.Red(fmt.Sprintf("%s 不存在方法 %s\n", clientName(), methodName))
			os.Exit(1)
		}

		values := method.Call(inputs)

		for _, val := range values {
			spew.Dump(val.Interface())
		}
	},
}

func clientName() string {
	return reflect.TypeOf(app.Storage.Storage()).Elem().String()
}

func title(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}
