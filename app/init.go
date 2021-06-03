package app

import (
	"github.com/v-mars/frame/db"
	"github.com/v-mars/sys/app/config"
	"github.com/v-mars/sys/app/logger"
	"github.com/v-mars/sys/app/models"
	"log"
)


func Run(configPath string)  {
	// 加载配置
	err := config.LoadConfig(configPath)
	if err != nil { log.Fatal(err) }
	conf := config.GetConf()
	logger.Initial(conf.Server.LogLevel,conf.Server.LogPath)
	initDB(&conf)

}

func initDB(config *config.Config){
	db.InitDB(config.Gorm.DBType, config.MySQL.DSN(), models.DB)
	//migrate.Migrate()
}




