package cmd

import (
	"github.com/spf13/cobra"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
	"os"
	"liulishuo/somechat/server/chatapp/conf"
	"liulishuo/somechat/server/chatapp/app"
	"liulishuo/somechat/core/data"
)

var chatRunConfFilePath string

var ChatRunCommand = &cobra.Command{
	Use:   "chat",
	Short: "Startup chat server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Log().Println(logs.Infof("read conf file, path= %s", chatRunConfFilePath).Trace())
		// read conf
		if confReadErr := conf.Read(chatRunConfFilePath); confReadErr != nil {
			log.Log().Println(logs.Errorf("read conf file failed, %v", confReadErr).Trace())
			os.Exit(1)
		}
		log.Log().Println(logs.Infof("read conf file succ. \n %v", conf.Conf))
		// build postgres
		databaseInitErr := data.InitPostgres(conf.Conf.Postgres.Url, conf.Conf.Postgres.MaxIdle, conf.Conf.Postgres.MaxOpen)
		if databaseInitErr != nil {
			log.Log().Println(logs.Errorf("database init failed, error : %v ", databaseInitErr))
			log.Log().Println(logs.Infof("close postgres, error : %v", data.Postgres().Close()))
			os.Exit(1)
			return
		}
		defer data.Postgres().Close()
		// TODO build redis

		// build chat app
		app.StartUp()
	},
}

func init()  {
	ChatRunCommand.PersistentFlags().StringVarP(&chatRunConfFilePath, "conf", "c", "", "the conf file path")
}