package grpc

import (
	"context"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/domain/permission"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PermissionHandler struct {
	pb.UnimplementedPermissionServiceServer

	Repo permission.PermissionRepository
}

func NewPermissionHandler(repo permission.PermissionRepository) *PermissionHandler {
	return &PermissionHandler{Repo: repo}
}

func (h *PermissionHandler) CreatePermission(_ context.Context, input *pb.CreatePermissionRequest) (*pb.Permission, error) {
	permission := permission.NewPermission(input.Codename, input.Name)
	_, err := h.Repo.AddPermission(permission)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Permission{
		Id:       permission.ID,
		Codename: permission.Codename,
		Name:     permission.Name,
	}, nil
}

func (h *PermissionHandler) FindPermissionById(_ context.Context, input *pb.FindPermissionByIdRequest) (*pb.Permission, error) {
	permission, err := h.Repo.FindPermissionById(input.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.Permission{
		Id:       permission.ID,
		Codename: permission.Codename,
		Name:     permission.Name,
	}, nil
}

func (h *PermissionHandler) DeletePermission(_ context.Context, input *pb.DeletePermissionRequest) (*pb.PermissionEmpty, error) {
	err := h.Repo.DeletePermission(input.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.PermissionEmpty{}, nil
}

func (h *PermissionHandler) FindAllPermissions(_ context.Context, _ *pb.PermissionEmpty) (*pb.FindAllPermissionsResponse, error) {
	permissions, err := h.Repo.FindAllPermissions()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var pbPermissions []*pb.Permission
	for _, permission := range permissions {
		pbPermissions = append(pbPermissions, &pb.Permission{
			Id:       permission.ID,
			Codename: permission.Codename,
			Name:     permission.Name,
		})
	}
	return &pb.FindAllPermissionsResponse{Permissions: pbPermissions}, nil
}
