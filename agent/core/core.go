package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/noaway/v2agent/config"
	"github.com/twinj/uuid"
	"google.golang.org/grpc"
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
	dialOnce = sync.Once{}
	connOnce *grpc.ClientConn
)

type Client struct {
}

func dial() (*grpc.ClientConn, error) {
	var err error
	dialOnce.Do(func() {
		addr := config.Configure().V2HandlerConfig.Addr
		if addr == "" {
			addr = DEFAULT_API_ADDR
		}
		connOnce, err = grpc.Dial(addr, grpc.WithInsecure())
	})
	return connOnce, err
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

	if u.Email == "" {
		return nil, fmt.Errorf("%v", "email is empty")
	}

	if u.UUID == "" {
		u.UUID = uuid.NewV4().String()
	}

	if u.AlterId == 0 {
		u.AlterId = 10
	}

	client := command.NewHandlerServiceClient(conn)
	if resp, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: tag(),
		Operation: serial.ToTypedMessage(
			&command.AddUserOperation{
				User: &protocol.User{
					Email: u.Email,
					Account: serial.ToTypedMessage(&vmess.Account{
						Id:      u.UUID,
						AlterId: u.AlterId,
					}),
				},
			}),
	}); err != nil {
		return nil, err
	} else {
		fmt.Println(resp)
	}

	return u, err
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

	if resp == nil {
		return nil, fmt.Errorf("%v", "stat is nil")
	}
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
