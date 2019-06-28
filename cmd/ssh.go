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

		// ssh: connect to host 0.0.0.12 port 22: No route to host
		//fmt.Println("ssh called")
		//fmt.Println(len(args))
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

		//fmt.Printf("%d****%s****\n", len(a), a[0])
		//fmt.Println("Print: " + strings.Join(args, " "))
		//fmt.Println("Print: " + strings.Join(args, "|"))
		//fmt.Println(args);
		//fmt.Println(SshServer)
		//initSqliteDb(saveSshServer)
		initSqliteDb()

		if hasServer() == false {
			saveSshServer()
		}
		toServer()
	},
}
// 参考
// https://www.cnblogs.com/borey/p/5715641.html
// https://github.com/emacski/redact/blob/2b1380f943e9e963758bfc3d6f6f71ac5ab01373/redact/cmd.go
func init() {
	rootCmd.AddCommand(sshCmd)
	//sshCmd.SetUsageTemplate("Usage:\n  cq vm stop [instance-id] [instance-id] ...\n")
	//sshCmd.PersistentFlags().Uint16VarP(&SshServer.Port, "port", "p", 22, "Port to connect to on the remote host.  This can be specified on a per-host basis in the configuration file.")
	//sshCmd.PersistentFlags().StringVarP(&SshServer.User,"/","","","remote Host")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func hasServer() bool {
	sql := "select id,username,alias,port,host,password,description,used_count from servers where username = '" + SshServer.User + "' and host = '" + SshServer.Host + "' limit 1"
	rows, err := DbDriver.Query(sql)
	checkErr(err)
	defer rows.Close() // 果然在多次操作 db 时会锁库
	for rows.Next() {
		return true
	}
	return false
}
