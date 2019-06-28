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

		//fmt.Println("rootCmd called")

		//fmt.Println("Print: " + strings.Join(args, " "))
		if SshServer.Id > 0 {
			delSshServer()
			os.Exit(0)
		}
		cmd.Help()

		//os.Exit(0)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

type Server struct {
	Id          uint32
	Alias       string
	Port        uint16
	User        string
	Host        string
	Password    string
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
//fmt.Printf("%s", os.UserHomeDir())
//fmt.Printf("%s", os.TempDir())

const (
	DataSourceName = "iserver.sqlite"
	UserDir        = string(os.PathSeparator) + "user" + string(os.PathSeparator) // 注册用户
	DatabaseName   = "servers"
	TimeFormat     = "2006/1/2"
	BeignTime      = 1551369600
	SoLikeU        = 1525795200
	MangoName      = ""
	MangoAvatar    = ""
	XzkName        = ""
	XzkAvatar      = ""
	FilePath       = "/Users/leili/"
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
	rootCmd.PersistentFlags().StringVarP(&SshServer.User, "user", "u", "", "user.")
	rootCmd.PersistentFlags().StringVarP(&SshServer.Host, "host", "l", "", "host.")
	rootCmd.PersistentFlags().StringVarP(&SshServer.Password, "passwd", "P", "", "password.")
	rootCmd.PersistentFlags().StringVarP(&SshServer.Description, "description", "d", "", "Description.")
	rootCmd.PersistentFlags().Uint16VarP(&SshServer.Port, "port", "p", 0, "Port to connect to on the remote host.  This can be specified on a per-host basis in the configuration file.")

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iserver.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().Uint32VarP(&SshServer.Id, "del", "D", 0, "Delete host record")
	rootCmd.Flags().Uint32Var(&SshServer.Id, "delete", 0, "Delete host record")
	//rootCmd.Flags().Uint32Var(&SshServer.Id, "export", 0, "Delete host record")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".iserver" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".iserver")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// 初始化db
func initSqliteDb() {
	//db, err := sql.Open("sqlite3", DataSourceName)
	//checkErr(err)

	res, err := DbDriver.Query("select name from sqlite_master where name='" + DatabaseName + "'")
	checkErr(err)

	if !res.Next() {
		// 实际上一步可以省略
		sqlTable := `
        CREATE TABLE IF NOT EXISTS servers (
            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            username VARCHAR(64) NOT NULL DEFAULT '',
            alias VARCHAR(64) NOT NULL DEFAULT '',
            port INT(11) NOT NULL DEFAULT '22',
            host VARCHAR(255)  NOT NULL DEFAULT '',
            password VARCHAR(255)  NOT NULL DEFAULT '',
            description VARCHAR(255)  NOT NULL DEFAULT '',
            used_count INT(11) NOT NULL DEFAULT '0',
            created_at INT(11) NOT NULL DEFAULT '0',
            updated_at INT(11) NOT NULL DEFAULT '0'
        );`
		DbDriver.Exec(sqlTable)
	} else {
		res.Close()
	}
}

// https://blog.csdn.net/LOVETEDA/article/details/82690498
func saveSshServer() int64 {
	//db, err := sql.Open("sqlite3", DataSourceName)
	//checkErr(err)
	stmt, err := DbDriver.Prepare("insert into servers(username,alias,port,host,password,description,created_at,updated_at) values(?,?,?,?,?,?,?,?)")
	checkErr(err)

	// func (s *Stmt) Exec(args ...interface{}) (Result, error)
	// Exec使用提供的参数执行准备好的命令状态，返回Result类型的该状态执行结果的总结。
	// result类型为接口，有两个方法：
	// LastInsertId() (int64, error)
	// RowsAffected() (int64, error)
	
	if SshServer.Port == 0 {
		SshServer.Port = 22
	}
	res, err := stmt.Exec(SshServer.User, SshServer.Alias, SshServer.Port, SshServer.Host, SshServer.Password, SshServer.Description, time.Now().Unix(), time.Now().Unix())
	checkErr(err)
	if err != nil {
		//fmt.Println(time.Now())
		//fmt.Println(err.Error())
		stmt.Close()
	}
	// RowsAffected() (int64, error)
	id, err := res.LastInsertId()
	checkErr(err)
	// 打印出受影响的id号
	//DbDriver.Close()
	return id
}

func delSshServer() bool {

	if SshServer.Id > 0 {
		stmt3, err := DbDriver.Prepare("delete from servers where id=?")
		checkErr(err)
		res, err := stmt3.Exec(SshServer.Id)
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
