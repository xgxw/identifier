package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/everywan/identifier/internal/controllers"
	"github.com/everywan/identifier/pb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var snowflakeCmd = &cobra.Command{
	Use:   "snowflake",
	Short: "Generate Unique ID by Snowflake",
	Long:  "Generate Unique ID by Snowflake. 通过调用grpc服务获取唯一ID. 目前使用使用配置文件的方式设置workerID, 后续若需多节点自动化部署, 需在节点间通信, 相互同步时间, 协商workerid值",
	Run: func(cmd *cobra.Command, args []string) {
		opts := loadApplocationOps()
		boot, err := newBootstrap(opts)
		if err != nil {
			fmt.Printf("bootstrap err, err=%v", err)
		}
		defer boot.Teardown()

		sfCtrl := controllers.NewSnowflakeController(boot.Logger, boot.SfSvc)
		grpcSvc := grpc.NewServer()
		pb.RegisterSnowflakeServer(grpcSvc, sfCtrl)
		go func() {
			lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", boot.Opts.Snowflake.Port))
			boot.Logger.Info("snowflake grpc server  start, will listening port: ", boot.Opts.Snowflake.Port)
			err := grpcSvc.Serve(lis)
			if err != nil {
				panic(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		grpcSvc.GracefulStop()
	},
}

type (
	// SnowflakeOps is ...
	SnowflakeOps struct {
		Port     uint  `mapstructure:"port" yaml:"port"`
		WorkerID int64 `mapstructure:"worker_id" yaml:"worker_id"`
	}
)

func init() {
	rootCmd.AddCommand(snowflakeCmd)
}
