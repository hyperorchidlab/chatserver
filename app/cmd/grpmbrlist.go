/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/kprc/chatserver/app/cmdcommon"
	"github.com/kprc/chatserver/app/cmdclient"
)

// grpmbrlistCmd represents the grpmbrlist command
var grpmbrlistCmd = &cobra.Command{
	Use:   "grpmbrlist",
	Short: "list group member ",
	Long: `list all group member in db or show group members with group id`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := cmdcommon.IsProcessStarted(); err != nil {
			log.Println(err)
			return
		}
		if len(args) > 1{
			log.Println("command error")
		}

		if len(args) == 0{
			cmdclient.StringOpCmdSend("", cmdcommon.CMD_LIST_GRPMBR, "")
		}else{
			cmdclient.StringOpCmdSend("", cmdcommon.CMD_LIST_GRPMBR, args[0])
		}

		//cmdclient.StringOpCmdSend("", cmdcommon.CMD_LIST_GRPMBR, keypassword)
	},
}

func init() {
	rootCmd.AddCommand(grpmbrlistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpmbrlistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpmbrlistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
