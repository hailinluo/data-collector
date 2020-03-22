package cmd

import (
	"bufio"
	"fmt"
	"github.com/hailinluo/data-collector/config"
	"github.com/hailinluo/data-collector/logger"
	"github.com/hailinluo/data-collector/storage"
	"github.com/hailinluo/data-collector/task"
	"github.com/hailinluo/data-collector/task/fund"
	"github.com/hailinluo/data-collector/task/fundcompany"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

var root = &cobra.Command{
	Use:   "data-fundcompany",
	Short: "It's an app for data collecting.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runServer(); err != nil {
			os.Exit(1)
		}
	},
}

var (
	configPtah string // 配置文件路径
	logPath    string // 日志文件存放路径
)

func init() {
	root.PersistentFlags().StringVarP(&configPtah, "config_path", "c", "./config/config.yaml", "config file path")
	root.PersistentFlags().StringVarP(&logPath, "log_path", "l", "./log/", "log file path")
}

func runServer() error {
	logger.Errorf("server start...")

	// 初始化配置
	err := config.InitConfig(configPtah)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	if config.Server.LogType == "file" {
		fileObj, err := os.Open(logPath)
		if err != nil {
			panic(fmt.Sprintf("os open error: %v", err))
		}
		defer fileObj.Close()
		writer := bufio.NewWriter(fileObj)
		logger.Init(logger.WithOutput(writer))
	}

	// 初始化 MySQL 数据库
	dbCloser, err := storage.InitDB(config.Server.DbUri)
	if err != nil {
		logger.Errorf("init db failed. err: %s", err)
		panic(err)
	}
	defer dbCloser.Close()

	// 初始化任务中心
	hub := task.InitTaskHub()
	defer hub.Close()
	// 添加任务
	fundcompany := fundcompany.NewFcCollector(
		fundcompany.WithSpec(config.Server.Tasks["fundcompany"]["spec"]),
		fundcompany.WithHomePage(config.Server.Tasks["fundcompany"]["home-page"]),
		fundcompany.WithResUrl(config.Server.Tasks["fundcompany"]["resource-url"]),
	)
	_, err = hub.AddTask(fundcompany)
	if err != nil {
		logger.Errorf("add task failed. err: %s", err)
	}

	fund := fund.NewFundCollector(
		fund.WithSpec(config.Server.Tasks["fundcompany"]["spec"]),
	)
	_, err = hub.AddTask(fund)
	if err != nil {
		logger.Errorf("add task failed. err: %s", err)
	}

	exit := make(chan bool)
	<-exit
	return nil
}
