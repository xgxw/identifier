package identifier

import (
	"context"
	"time"

	"github.com/everywan/identifier/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultDialTimeout      = 10
	defaultKeepAliveTime    = 600
	defaultKeepAliveTimeout = 20
)

// SnowflakeClientOps snowflake client 配置
type SnowflakeClientOps struct {
	// Address 服务的地址，IP和端口
	Address string `json:"address" yaml:"address" mapstructure:"address"`

	// DialTimeout 连接超时时间，单位秒
	DialTimeout int64 `json:"dial_timeout" yaml:"dial_timeout" mapstructure:"dial_timeout"`

	// KeepAliveTime 连接保活周期，单位秒
	KeepAliveTime int64 `json:"keep_alive_time" yaml:"keep_alive_time" mapstructure:"keep_alive_time"`

	// KeepAliveTimeout 发送保活心跳包的超时时间，单位秒
	KeepAliveTimeout int64 `json:"keep_alive_timeout" yaml:"keep_alive_timeout" mapstructure:"keep_alive_timeout"`
}

func (s *SnowflakeClientOps) loadDefault() {
	if s.DialTimeout == 0 {
		s.DialTimeout = defaultDialTimeout
	}
	if s.KeepAliveTime == 0 {
		s.KeepAliveTime = defaultKeepAliveTime
	}
	if s.KeepAliveTimeout == 0 {
		s.KeepAliveTimeout = defaultKeepAliveTimeout
	}
}
func (s *SnowflakeClientOps) buildDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(s.KeepAliveTime) * time.Second,
			Timeout:             time.Duration(s.KeepAliveTimeout) * time.Second,
			PermitWithoutStream: true,
		}),
	}
}

// ISnowflakeClient Client 接口, 添加Close方法
type ISnowflakeClient interface {
	pb.SnowflakeClient
	Close() error
}

// SnowflakeClient 实际返回的 client 结构体
type SnowflakeClient struct {
	pb.SnowflakeClient
	conn *grpc.ClientConn
}

// NewSnowflakeClient 生成Snowflake grpc客户端
func NewSnowflakeClient(opts *SnowflakeClientOps) (client *SnowflakeClient, err error) {
	opts.loadDefault()
	client = new(SnowflakeClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.DialTimeout)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, opts.Address, opts.buildDialOptions()...)
	if err != nil {
		return client, err
	}
	client = &SnowflakeClient{
		SnowflakeClient: pb.NewSnowflakeClient(conn),
		conn:            conn,
	}
	return client, nil
}

// Close 关闭链接
func (s *SnowflakeClient) Close() error {
	return s.conn.Close()
}
