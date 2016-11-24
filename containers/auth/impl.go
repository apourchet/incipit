package main

import (
	"net/http"

	"github.com/apourchet/incipit/lib/auth"
	protos "github.com/apourchet/incipit/protos/go"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type AuthService struct {
	client auth.AuthClient
	logger protos.LoggerClient
}

func NewAuthService(client auth.AuthClient, logger protos.LoggerClient) *AuthService {
	return &AuthService{client, logger}
}

func (service *AuthService) UserExists(ctx context.Context, in *protos.UserExistsReq) (*protos.UserExistsRes, error) {
	glog.Infof("AuthService: UserExists")
	res := &protos.UserExistsRes{}

	found, err := service.client.UserExists(in.Key)
	res.Found = found
	if err != nil {
		return res, grpc.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	glog.Infof("AuthService: -UserExists")
	return res, nil
}
func (service *AuthService) Register(ctx context.Context, req *protos.RegisterReq) (*protos.RegisterRes, error) {
	glog.Infof("AuthService: Register")
	res := &protos.RegisterRes{}

	ok, err := service.client.Register(req.Key, req.Pass)
	if err != nil {
		glog.Errorf("Failed to register user: %v", err)
		return res, grpc.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	if !ok {
		glog.Errorf("Failed to register user: conflict")
		return res, grpc.Errorf(codes.AlreadyExists, codes.AlreadyExists.String())
	}

	glog.Infof("AuthService: -Register")
	return res, nil
}

func (service *AuthService) Login(ctx context.Context, req *protos.LoginReq) (*protos.LoginRes, error) {
	glog.Infof("AuthService: Login")
	res := &protos.LoginRes{}

	token, valid, err := service.client.Login(req.Key, req.Pass)
	if err != nil {
		glog.Errorf("Failed to validate token: %v", err)
		return res, grpc.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	if !valid {
		glog.Errorf("Invalid token")
		return res, grpc.Errorf(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	go service.logger.LogLogin(ctx, &protos.LogLoginReq{req.Key})

	res.Token = token
	glog.Infof("AuthService: -Login")
	return res, nil
}

func (service *AuthService) Logout(ctx context.Context, req *protos.LogoutReq) (*protos.LogoutRes, error) {
	glog.Infof("AuthService: Logout")
	res := &protos.LogoutRes{}

	token, err := auth.GetToken(ctx)
	if err != nil {
		glog.Errorf("Failed to retrieve token from headers: %v", err)
		return res, grpc.Errorf(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	err = service.client.Logout(token)
	if err != nil {
		glog.Errorf("Failed to invalidate token: %v", err)
		return res, grpc.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	glog.Infof("AuthService: -Logout")
	return res, nil
}

func (service *AuthService) Deregister(ctx context.Context, req *protos.DeregisterReq) (*protos.DeregisterRes, error) {
	glog.Infof("AuthService: Deregister")
	res := &protos.DeregisterRes{}

	token, err := auth.GetToken(ctx)
	if err != nil {
		glog.Errorf("Failed to retrieve token from headers: %v", err)
		return res, grpc.Errorf(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	_, valid, err := service.client.Validate(token)
	if err != nil {
		glog.Errorf("Failed to validate token: %v", err)
		return res, grpc.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	if !valid {
		glog.Errorf("Invalid Token")
		return res, grpc.Errorf(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	err = service.client.Deregister(token)
	if err != nil {
		glog.Errorf("Failed to deregister user: %v", err)
		return res, grpc.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	glog.Infof("AuthService: -Deregister")
	return res, nil
}
