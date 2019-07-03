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
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)


// toCmd represents the to command
var toCmd = &cobra.Command{
	Use:   "to [Id|Alias]",
	Short: "Quickly connect to the server.",
	Long:  `Quickly connect to the server.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			showServers()
			os.Exit(0)
		}
		Id, _ := strconv.ParseUint(args[0], 10, 32)
		if Id > 0 {
			SshServer.Id = uint32(Id)
		} else {
			SshServer.Alias = args[0]
		}

		if GetServer() == true {
			updataServer()
			toServer()
		} else {
			fmt.Println("No host available!")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(toCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	toCmd.Flags().BoolVarP(&copySshId,"copy","i",false,"222")
	viper.BindPFlag("copy", toCmd.Flags().Lookup("copy"))
	viper.SetDefault("copy",false)
	//toCmd.PersistentFlags().Bool("copy", true, "Use Viper for configuration")

}

func GetServer() bool {
	SshTmpServer = SshServer

	var sql string
	if SshServer.Id > 0 {
		sql = "select id,username,alias,port,host,password,description,used_count from servers where id='" + fmt.Sprint(SshServer.Id) + "' or alias like '%" + fmt.Sprint(SshServer.Id) + "%' limit 1"
	} else {
		sql = "select id,username,alias,port,host,password,description,used_count from servers where alias like '%" + SshServer.Alias + "%' limit 1"
	}
	//fmt.Println(sql)
	rows, err := DbDriver.Query(sql)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&SshServer.Id, &SshServer.User, &SshServer.Alias, &SshServer.Port, &SshServer.Host, &SshServer.Password, &SshServer.Description, &SshServer.Count)
		checkErr(err)

		if SshTmpServer.User != "" && SshTmpServer.User != SshServer.User {
			SshServer.User = SshTmpServer.User
		}

		if SshTmpServer.Port > 0 && SshTmpServer.Port != SshServer.Port {
			SshServer.Port = SshTmpServer.Port
		}

		if SshTmpServer.Alias != "" && SshTmpServer.Alias != SshServer.Alias {
			SshServer.Alias = SshTmpServer.Alias
		}

		if SshTmpServer.Host != "" && SshTmpServer.Host != SshServer.Host {
			SshServer.Host = SshTmpServer.Host
		}

		if SshTmpServer.Password != "" && SshTmpServer.Password != SshServer.Password {
			SshServer.Password = SshTmpServer.Password
		}

		if SshTmpServer.Description != "" && SshTmpServer.Description != SshServer.Description {
			SshServer.Description = SshTmpServer.Description
		}

		//fmt.Println(SshServer)
		return true
	}

	return false
}

func showServers() {
	sql := "select id,username,alias,port,host,password,description,used_count,tags from servers order by tags, used_count desc"
	//fmt.Println(sql)
	rows, err := DbDriver.Query(sql)
	defer rows.Close()
	checkErr(err)

	ServerList := make([]Server, 0)
	for rows.Next() {
		var CurServer Server

		err = rows.Scan(&CurServer.Id, &CurServer.User, &CurServer.Alias, &CurServer.Port, &CurServer.Host, &CurServer.Password, &CurServer.Description, &CurServer.Count,&CurServer.Tags)
		if len(CurServer.Password) > 6 {
			CurServer.Password = CurServer.Password[0:6] + "***"
		}
		checkErr(err)
		ServerList = append(ServerList, CurServer)
	}

	if len(ServerList) == 0 {
		fmt.Println("No valid server record")
		return
	}
	table.Output(ServerList)
	//table.OutputA(ServerList)
	s := table.Table(ServerList)
	_ = s
	return
}

func updataServer() {
	stmt2, err := DbDriver.Prepare("update servers set username=?,alias=?,port=?,host=?,password=?,description=?,tags=?,used_count=used_count+1 where id=?")
	checkErr(err)
	res, err := stmt2.Exec(SshServer.User, SshServer.Alias, SshServer.Port, SshServer.Host, SshServer.Password, SshServer.Description, SshServer.Tags, SshServer.Id)
	checkErr(err)
	defer stmt2.Close()
	affect, err := res.RowsAffected()
	checkErr(err)
	_ = affect
}

