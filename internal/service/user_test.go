package service

import (
	"context"
	"errors"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/model/entity"
	"robot_scheduler/internal/testutil/mocks"
	"testing"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func init() {
	// Initialize logger for tests
	logger.Logger = zap.NewNop()
}

func TestUserService_CreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserDAO := mocks.NewMockUserDAO(ctrl)
	service := NewUserService(mockUserDAO)

	ctx := context.Background()
	req := &dto.UserCreateRequest{
		UserName: "testuser",
		Password: "password123",
		Role:     entity.RoleUser,
	}
	desKey := "12345678"

	// Mock expectations
	mockUserDAO.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *entity.UserInfo) error {
			user.ID = 1
			return nil
		})

	// Execute
	resp, err := service.CreateUser(ctx, req, desKey)

	// Assert
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if resp.UserName != req.UserName {
		t.Errorf("Expected username %s, got %s", req.UserName, resp.UserName)
	}
}

func TestUserService_CreateUser_EncryptionFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserDAO := mocks.NewMockUserDAO(ctrl)
	service := NewUserService(mockUserDAO)

	ctx := context.Background()
	req := &dto.UserCreateRequest{
		UserName: "testuser",
		Password: "password123",
		Role:     entity.RoleUser,
	}
	invalidKey := "short" // Invalid DES key

	// Execute
	_, err := service.CreateUser(ctx, req, invalidKey)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid encryption key")
	}
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserDAO := mocks.NewMockUserDAO(ctrl)
	service := NewUserService(mockUserDAO)

	ctx := context.Background()
	userID := uint(1)

	expectedUser := &entity.UserInfo{
		UserName: "testuser",
		Role:     entity.RoleUser,
	}
	expectedUser.ID = userID

	// Mock expectations
	mockUserDAO.EXPECT().
		FindByID(ctx, userID).
		Return(expectedUser, nil)

	// Execute
	resp, err := service.GetUserByID(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if resp.ID != userID {
		t.Errorf("Expected ID %d, got %d", userID, resp.ID)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserDAO := mocks.NewMockUserDAO(ctrl)
	service := NewUserService(mockUserDAO)

	ctx := context.Background()
	userID := uint(999)

	// Mock expectations
	mockUserDAO.EXPECT().
		FindByID(ctx, userID).
		Return(nil, nil)

	// Execute
	resp, err := service.GetUserByID(ctx, userID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp != nil {
		t.Error("Expected nil response for non-existent user")
	}
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserDAO := mocks.NewMockUserDAO(ctrl)
	service := NewUserService(mockUserDAO)

	ctx := context.Background()
	userID := uint(1)

	// Mock expectations
	mockUserDAO.EXPECT().
		Delete(ctx, userID).
		Return(nil)

	// Execute
	err := service.DeleteUser(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
}

func TestUserService_DeleteUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserDAO := mocks.NewMockUserDAO(ctrl)
	service := NewUserService(mockUserDAO)

	ctx := context.Background()
	userID := uint(999)

	// Mock expectations
	mockUserDAO.EXPECT().
		Delete(ctx, userID).
		Return(errors.New("user not found"))

	// Execute
	err := service.DeleteUser(ctx, userID)

	// Assert
	if err == nil {
		t.Error("Expected error when deleting non-existent user")
	}
}
