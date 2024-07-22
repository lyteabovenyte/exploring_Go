package server

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/lyteabovenyte/exploring_go/grpc/proto"
)

type API struct {
	pb.UnimplementedQOTDServer // contains all the methods that our server must implement

	addr   string
	quotes map[string][]string

	mu         sync.Mutex
	grpcServer *grpc.Server
}

func New(addr string) (*API, error) {
	var opts []grpc.ServerOption

	a := &API{
		addr: addr,
		quotes: map[string][]string{
			"mo seif": {
				"donya sar e jang dare.",
			},
			"mohammad hossein seyfi": {
				"Chamassen, is a Lifestyle, not just a State",
				"ina ham avvalesh adam bodan, kam kam injori shodan",
			},
			"mo seyf": {
				"LinkedIn o daryab",
			},
		},
		grpcServer: grpc.NewServer(opts...),
	}
	a.grpcServer.RegisterService(&pb.QOTD_ServiceDesc, a)

	return a, nil
}

func (a *API) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	lis, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(lis)
}

func (a *API) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.grpcServer.Stop()
}

func (a *API) GetQOTD(ctx context.Context, req *pb.GetReq) (*pb.GetResp, error) {
	var (
		author string
		quotes []string
	)

	if req.Author == "" {
		for author, quotes = range a.quotes {
			break
		}
	} else {
		author = req.Author
		var ok bool
		quotes, ok = a.quotes[req.Author]
		if !ok {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("author %q not found", req.Author))
		}
	}

	return &pb.GetResp{
		Author: author,
		Quote:  quotes[rand.Intn(len(quotes))],
	}, nil
}
