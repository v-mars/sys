package server

import (
	"fmt"
	"github.com/v-mars/sys/app"
	"github.com/v-mars/sys/app/config"
	"github.com/v-mars/sys/app/router"
	"os"

	"github.com/spf13/cobra"
)

var (
	//h          bool
	c        = &config.Configs

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of AirCloud",
		Long:  `This is AirCloud`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("airCloud version is v1.0")
		},
	}

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "mars system初始化数据库",
		Long:  `This is mars system`,
		Run: func(cmd *cobra.Command, args []string) {
			//migrate.Migrate()
		},
	}

	rootCmd = &cobra.Command{
		Use:   "",
		Short: "Welcome Use Mars system.",
		Long:  `Mars system e.g.`,
		Example: `## 启动命令 ./app -p 5000 -c ./configs/config.toml -f ./logs`,
		Run: func(cmd *cobra.Command, args []string) {
			app.Run(c.Server.ConfigFile)
			addr := c.Server.IP+":"+c.Server.Port
			router.InitRouter(c.Server.Mode, addr)
		},
	}
)

func init() {

	DefaultConfig := ""
	DefaultIP := "0.0.0.0"
	DefaultPort := "5000"
	DefaultMode := ""		   // release, debug, test
	DefaultLevel := ""
	DefaultLog := "./logs"
	rootCmd.Flags().StringVarP(&c.Server.ConfigFile,"config","c",DefaultConfig,"config file, example: ./configs/config.toml")
	rootCmd.Flags().StringVarP(&c.Server.IP, "ip", "i", DefaultIP, "服务IP")
	rootCmd.Flags().StringVarP(&c.Server.Port, "port", "p", DefaultPort, "服务启动的端口: 3000")
	//rootCmd.Flags().IntVarP(&c.Server.Port, "times", "t", 1, "times to echo the input")
	rootCmd.Flags().StringVarP(&c.Server.Mode, "mode", "m", DefaultMode, "启动模式(release, debug, test e.g)")
	rootCmd.Flags().StringVarP(&c.Server.LogPath, "log", "f", DefaultLog, "日志目录(/data/logs e.g)")
	rootCmd.Flags().StringVarP(&c.Server.LogLevel, "level", "l", DefaultLevel, "日志级别(DEBUG, INFO, WARNING e.g)")
	//cmd.AddFlags(rootCmd)
	rootCmd.AddCommand(versionCmd, migrateCmd)
	if c.Server.ConfigFile == "" {
		c.Server.ConfigFile = "./configs/config.toml"
		//fmt.Println("请使用\"-c\"指定配置文件")
		//os.Exit(-1)
	}
}



func Run() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}