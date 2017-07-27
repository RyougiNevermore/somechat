package cmd

import (
	"github.com/spf13/cobra"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
	"os"
	"liulishuo/somechat/server/webapp/conf"
	"liulishuo/somechat/server/webapp/app"
	"liulishuo/somechat/core/data"
)

var runConfFilePath string

var WebRunCommand = &cobra.Command{
	Use:   "web",
	Short: "Startup web server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Log().Println(logs.Infof("read conf file, path= %s", runConfFilePath).Trace())
		// read conf
		if confReadErr := conf.Read(runConfFilePath); confReadErr != nil {
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

		// build web app
		app.StartUp()
	},
}

func init()  {
	WebRunCommand.PersistentFlags().StringVarP(&runConfFilePath, "conf", "c", "", "the conf file path")
}