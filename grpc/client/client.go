package client

import (
	"context"
	pb "github.com//lyteabovenyte/exploring_go/grpc/proto"
	"time"

	"google.golang.org/grpc"
)

type Client struct {
	client pb.QOTDClient
	conn   *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}
	return &Client{
		client: pb.NewQOTDClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) QOTD(ctx context.Context, wantAuthor string) (author, quote string, err error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
	}
	resp, err := c.client.GetQOTD(ctx, &pb.GetReq{author: wantAuthor})
	if err != nil {
		return "", "", err
	}
	return resp.author, resp.quote, nil
}
