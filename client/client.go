package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/noaway/v2agent/agent"
	"github.com/noaway/v2agent/config"
	"github.com/twinj/uuid"
	"google.golang.org/grpc"
	"v2ray.com/core"
	"v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

const (
	DEFAULT_API_ADDR = "127.0.0.1:10086"
	DEFAULT_TAG      = "proxy"
)

var (
	once          sync.Once
	defaultClient *Client
)

func InitClient() {
	once.Do(func() {
		defaultClient = &Client{}
		conf := config.Configure()
		ag, err := agent.NewConfig(
			agent.SetupCluster(
				conf.Server.AdvertiseAddr,
				conf.Server.BindAddr,
				conf.JoinClusterAddrs...,
			),
			agent.SetupUserEventHandler(defaultClient.UserEventHandler),
			agent.SetupDataDir(conf.DataDir),
			agent.SetupNodeName(conf.Name),
		).NewAgent()
		if err != nil {
			panic(err)
		}
		defaultClient.Agent = ag
	})
}

type Client struct {
	*agent.Agent
}

func (c *Client) UserEventHandler(e agent.UserEvent) {

}

type User struct {
	Email    string
	UUID     string
	Level    uint32
	AlterId  uint32
	Security int32
}

func dial() (command.HandlerServiceClient, error) {
	addr := config.Configure().V2HandlerConfig.Addr
	if addr == "" {
		addr = DEFAULT_API_ADDR
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	return command.NewHandlerServiceClient(conn), err
}

func tag() string {
	tag := config.Configure().V2HandlerConfig.Tag
	if tag == "" {
		tag = DEFAULT_TAG
	}
	return tag
}

func AddUser(u *User) (*User, error) {
	client, err := dial()
	if err != nil {
		return nil, err
	}

	if u.Email == "" {
		return nil, fmt.Errorf("%v", "email is empty")
	}

	if u.UUID == "" {
		u.UUID = uuid.NewV4().String()
	}

	if u.AlterId == 0 {
		u.AlterId = 10
	}

	if u.Security == 0 {
		u.Security = 2
	}

	client.AddInbound(context.Background(), &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: tag(),
			ProxySettings: serial.ToTypedMessage(&command.AddUserOperation{
				User: &protocol.User{
					Level: u.Level,
					Email: u.Email,
					Account: serial.ToTypedMessage(&vmess.Account{
						Id:               u.UUID,
						AlterId:          u.AlterId,
						SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType(u.Security)},
					}),
				},
			}),
		},
	})

	// _, err = client.AddInbound(context.Background(), &command.AlterInboundRequest{
	// Tag: tag(),
	// Operation: serial.ToTypedMessage(&command.AddUserOperation{
	// 	User: &protocol.User{
	// 		Level: u.Level,
	// 		Email: u.Email,
	// 		Account: serial.ToTypedMessage(&vmess.Account{
	// 			Id:               u.UUID,
	// 			AlterId:          u.AlterId,
	// 			SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
	// 		}),
	// 	},
	// }),
	// })

	if err != nil {
		return nil, err
	}
	return u, nil
}

func DelUser(email string) error {
	client, err := dial()
	if err != nil {
		return err
	}
	_, err = client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag:       tag(),
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{Email: email}),
	})
	return err
}
