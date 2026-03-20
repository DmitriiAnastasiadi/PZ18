package grpcserver

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tech-ip-sem2/proto/authpb"
)

// AuthGRPCServer implements authpb.AuthServiceServer for token verification.
type AuthGRPCServer struct {
	authpb.UnimplementedAuthServiceServer
}

func (s *AuthGRPCServer) Verify(ctx context.Context, req *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	if req == nil || strings.TrimSpace(req.Token) == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if req.Token != "demo-token" {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &authpb.VerifyResponse{Valid: true, Subject: "student"}, nil
}

// HTTP-compatible helper for local endpoint compatibility.
func VerifyToken(token string) (bool, string) {
	if token == "demo-token" {
		return true, "student"
	}
	return false, ""
}

// VerifyHandler is optional HTTP fallback for /v1/auth/verify.
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"valid":false,"error":"unauthorized"}`))
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	valid, subject := VerifyToken(token)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"valid":false,"error":"unauthorized"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"valid":true,"subject":"` + subject + `"}`))
}
