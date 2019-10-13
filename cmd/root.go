/*
Copyright © 2019 lonisy@163.com

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
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iserver",
	Short: "Iserver tool",
	Long:  `Iserver is a tool for managing ssh servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if SshServer.Id > 0 {
			delSshServer()
			os.Exit(0)
		}
		cmd.Help()
	},
}

type Server struct {
	Id          uint32
	Alias       string
	Port        uint16
	User        string
	Host        string
	Password    string
	Tags        string
	Description string
	Count       uint32
}
type Callback func() int64

var SshServer = Server{}
var SshTmpServer = Server{}
var importServer bool
var exportServer bool
var copySshId bool

var DbDriver *sql.DB

const (
	DataSourceName = ".iserver.sqlite"
	DatabaseName   = "servers"
	UserDir        = string(os.PathSeparator) + "user" + string(os.PathSeparator) // 注册用户
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DBInit() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Could not find local user folder. Error: %v\n", err)
	}
	DbDriver, err = sql.Open("sqlite3", userHomeDir+string(os.PathSeparator)+DataSourceName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	DBInit()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&SshServer.Alias, "alias", "a", "", "host alias.")
	rootCmd.PersistentFlags().StringVarP(&SshServer.Tags, "tags", "t", "", "host tags.")
	rootCmd.PersistentFlags().StringVarP(&SshServer.User, "user", "u", "", "user.")
	rootCmd.PersistentFlags().StringVarP(&SshServer.Host, "host", "l", "", "host.")
	rootCmd.PersistentFlags().StringVarP(&SshServer.Password, "passwd", "P", "", "password.")
	rootCmd.PersistentFlags().StringVarP(&SshServer.Description, "description", "d", "", "Description.")
	rootCmd.PersistentFlags().Uint16VarP(&SshServer.Port, "port", "p", 0, "Port to connect to on the remote host.  This can be specified on a per-host basis in the configuration file.")

	rootCmd.Flags().Uint32Var(&SshServer.Id, "delete", 0, "Delete host record")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".iserver")
	}

	viper.AutomaticEnv()
	// read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initSqliteDb() {
	res, err := DbDriver.Query("select name from sqlite_master where name='" + DatabaseName + "'")
	checkErr(err)
	defer res.Close()
	if !res.Next() {
		sqlTable := `
        CREATE TABLE IF NOT EXISTS servers (
            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            username VARCHAR(64) NOT NULL DEFAULT '',
            alias VARCHAR(64) NOT NULL DEFAULT '',
            port INT(11) NOT NULL DEFAULT '22',
            host VARCHAR(255)  NOT NULL DEFAULT '',
            password VARCHAR(255)  NOT NULL DEFAULT '',
            description VARCHAR(255)  NOT NULL DEFAULT '',
            tags VARCHAR(64) NOT NULL DEFAULT '',
            used_count INT(11) NOT NULL DEFAULT '0',
            created_at INT(11) NOT NULL DEFAULT '0',
            updated_at INT(11) NOT NULL DEFAULT '0'
        );`
		DbDriver.Exec(sqlTable)
	}
}

func saveSshServer() int64 {
	stmt, err := DbDriver.Prepare("insert into servers(username,alias,port,host,password,description,tags,created_at,updated_at) values(?,?,?,?,?,?,?,?,?)")
	checkErr(err)
	if SshServer.Port == 0 {
		SshServer.Port = 22
	}
	res, err := stmt.Exec(SshServer.User, SshServer.Alias, SshServer.Port, SshServer.Host, SshServer.Password, SshServer.Description,SshServer.Tags, time.Now().Unix(), time.Now().Unix())
	checkErr(err)
	defer stmt.Close()
	//if err != nil {
	//	stmt.Close()
	//}
	// RowsAffected() (int64, error)
	id, err := res.LastInsertId()
	checkErr(err)
	//DbDriver.Close()
	return id
}

func delSshServer() bool {

	if SshServer.Id > 0 {
		stmt3, err := DbDriver.Prepare("delete from servers where id=?")
		checkErr(err)
		res, err := stmt3.Exec(SshServer.Id)
		defer stmt3.Close()
		checkErr(err)
		affect2, err := res.RowsAffected()
		checkErr(err)
		if affect2 == 1 {
			return true
		}
		fmt.Println("No valid server record")
	}
	return false
}

func toServer() {
	i := viper.GetBool("copy")
	var command string
	if i {
		command = fmt.Sprintf("ssh-copy-id -p %d %s@%s", SshServer.Port, SshServer.User, SshServer.Host)
	} else {
		command = fmt.Sprintf("ssh -p %d %s@%s", SshServer.Port, SshServer.User, SshServer.Host)
	}
	fmt.Println("Commands: " + command)
	fmt.Println("Password: " + SshServer.Password)
	fmt.Println("Description: " + SshServer.Description)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
