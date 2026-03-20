package authclient

import (
	"context"
	"errors"
	"time"

	"tech-ip-sem2/proto/authpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrUnavailable     = errors.New("unavailable")
)

// Client is a gRPC client for auth verification.
type Client struct {
	conn   *grpc.ClientConn
	client authpb.AuthServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, client: authpb.NewAuthServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) Verify(ctx context.Context, token string) (string, bool, error) {
	if c.conn == nil || c.conn.GetState() == connectivity.Shutdown {
		return "", false, ErrUnavailable
	}

	resp, err := c.client.Verify(ctx, &authpb.VerifyRequest{Token: token})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.Unauthenticated {
			return "", false, ErrUnauthenticated
		}
		if st.Code() == codes.DeadlineExceeded || st.Code() == codes.Unavailable {
			return "", false, ErrUnavailable
		}
		return "", false, err
	}

	return resp.Subject, resp.Valid, nil
}
