package core

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
	statsService "v2ray.com/core/app/stats/command"
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
		// defaultClient = &Client{}
		// conf := config.Configure()
		// ag, err := agent.NewConfig(
		// 	agent.SetupCluster(
		// 		conf.Server.AdvertiseAddr,
		// 		conf.Server.BindAddr,
		// 		conf.JoinClusterAddrs...,
		// 	),
		// 	agent.SetupUserEventHandler(defaultClient.UserEventHandler),
		// 	agent.SetupDataDir(conf.DataDir),
		// 	agent.SetupNodeName(conf.Name),
		// ).NewAgent()
		// if err != nil {
		// 	panic(err)
		// }
		// defaultClient.Agent = ag
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

func dial() (*grpc.ClientConn, error) {
	addr := config.Configure().V2HandlerConfig.Addr
	if addr == "" {
		addr = DEFAULT_API_ADDR
	}
	fmt.Println("dial.addr: ", addr)
	return grpc.Dial(addr, grpc.WithInsecure())
}

func tag() string {
	tag := config.Configure().V2HandlerConfig.Tag
	if tag == "" {
		tag = DEFAULT_TAG
	}
	return tag
}

func AddUser(u *User) (*User, error) {
	conn, err := dial()
	if err != nil {
		return nil, err
	}
	client := command.NewHandlerServiceClient(conn)

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

	if err != nil {
		return nil, err
	}
	return u, nil
}

func DelUser(email string) error {
	conn, err := dial()
	if err != nil {
		return err
	}
	client := command.NewHandlerServiceClient(conn)
	_, err = client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag:       tag(),
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{Email: email}),
	})
	return err
}

type Dosage struct {
	Email string
	Value string
}

func QueryStats() ([]Dosage, error) {
	conn, err := dial()
	if err != nil {
		return nil, err
	}

	client := statsService.NewStatsServiceClient(conn)
	r := &statsService.QueryStatsRequest{}
	resp, err := client.QueryStats(context.Background(), r)
	if err != nil {
		return nil, err
	}

	dos := make([]Dosage, len(resp.Stat))
	for i, s := range resp.Stat {
		dos[i] = Dosage{
			Email: s.GetName(),
			Value: Beautify(s.GetValue()),
		}
	}
	return dos, nil
}

func GetStats(email string) (*Dosage, error) {
	conn, err := dial()
	if err != nil {
		return nil, err
	}

	client := statsService.NewStatsServiceClient(conn)
	r := &statsService.GetStatsRequest{Name: fmt.Sprintf(`user>>>%v>>>traffic>>>uplink`, email)}
	resp, err := client.GetStats(context.Background(), r)

	return &Dosage{
		Email: email,
		Value: Beautify(resp.Stat.GetValue()),
	}, err
}

func Beautify(v int64) string {
	if v == 0 {
		return ""
	}
	if v < 1024 {
		return fmt.Sprintf("%v b", v)
	}
	sh := float64(v / 1024)
	if sh < 1024 {
		return fmt.Sprintf("%v kb", sh)
	}
	sh = sh / 1024
	if sh < 1024 {
		return fmt.Sprintf("%v mb", sh)
	}
	sh = sh / 1024
	return fmt.Sprintf("%v gb", sh)
}
