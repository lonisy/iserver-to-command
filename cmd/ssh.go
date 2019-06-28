/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh [user@host]",
	Short: "Quickly connect to the server and save server.",
	Long:  `Quickly connect to the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(0)
		}
		a := strings.Split(strings.Join(args, ""), "@")
		if len(a) != 2 {
			cmd.Help()
			os.Exit(0)
		}
		SshServer.User = a[0]
		SshServer.Host = a[1]
		initSqliteDb()
		if hasServer() == false {
			saveSshServer()
		}
		toServer()
	},
}
func init() {
	rootCmd.AddCommand(sshCmd)
}

func hasServer() bool {
	sql := "select id,username,alias,port,host,password,description,used_count from servers where username = '" + SshServer.User + "' and host = '" + SshServer.Host + "' limit 1"
	rows, err := DbDriver.Query(sql)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}
