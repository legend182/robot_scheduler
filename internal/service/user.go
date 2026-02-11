package service

import (
	"context"
	"errors"

	dao "robot_scheduler/internal/dao/interfaces"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/model/entity"
	"robot_scheduler/internal/utils"

	"go.uber.org/zap"
)

// UserService 用户服务
type UserService struct {
	userDAO dao.UserDAO
}

func NewUserService(userDAO dao.UserDAO) *UserService {
	return &UserService{
		userDAO: userDAO,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *dto.UserCreateRequest, desKey string) (*dto.UserResponse, error) {
	logger.Info("creating user in service", zap.String("username", req.UserName))

	// 检查用户名是否已存在
	// existingUser, err := s.userDAO.FindByUserName(ctx, req.UserName)
	// if err != nil {
	// 	logger.Error("failed to check username existence", zap.Error(err), zap.String("username", req.UserName))
	// 	return nil, err
	// }

	// if existingUser != nil {
	// 	logger.Warn("username already exists", zap.String("username", req.UserName))
	// 	return nil, errors.New("username already exists")
	// }

	// 创建前加密密码
	encryptedPassword, err := utils.DESEncrypt(req.Password, desKey)
	if err != nil {
		logger.Error("failed to encrypt user password", zap.Error(err), zap.String("username", req.UserName))
		return nil, errors.New("加密用户密码失败")
	}

	// 创建用户实体
	user := &entity.UserInfo{
		UserName:  req.UserName,
		Password:  encryptedPassword,
		Role:      req.Role,
		IsLocked:  0,
		ExtraInfo: req.ExtraInfo,
	}

	// 保存到数据库
	if err := s.userDAO.Create(ctx, user); err != nil {
		logger.Error("failed to create user in service", zap.Error(err), zap.String("username", req.UserName))
		return nil, err
	}

	logger.Info("user created successfully in service", zap.String("username", req.UserName), zap.Uint("id", user.ID))
	return dto.NewUserResponseFromEntity(user), nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, id uint, req *dto.UserUpdateRequest, desKey string) error {
	logger.Info("updating user in service", zap.Uint("id", id))

	// 获取用户
	user, err := s.userDAO.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find user for update", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if user == nil {
		logger.Warn("user not found for update", zap.Uint("id", id))
		return errors.New("user not found")
	}

	// 更新字段
	if req.Password != nil {
		encryptedPassword, err := utils.DESEncrypt(*req.Password, desKey)
		if err != nil {
			logger.Error("failed to encrypt user password", zap.Error(err), zap.Uint("id", id))
			return errors.New("加密用户密码失败")
		}
		user.Password = encryptedPassword
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsLocked != nil {
		user.IsLocked = *req.IsLocked
	}
	if req.ExtraInfo != nil {
		user.ExtraInfo = req.ExtraInfo
	}

	// 保存更新
	if err := s.userDAO.Update(ctx, user); err != nil {
		logger.Error("failed to update user in service", zap.Error(err), zap.Uint("id", id))
		return err
	}

	logger.Info("user updated successfully in service", zap.Uint("id", id))
	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	logger.Info("deleting user in service", zap.Uint("id", id))
	return s.userDAO.Delete(ctx, id)
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	logger.Debug("getting user by id in service", zap.Uint("id", id))
	user, err := s.userDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return dto.NewUserResponseFromEntity(user), nil
}

// ListUsers 分页获取用户列表
func (s *UserService) ListUsers(ctx context.Context, req dto.PageRequest) (*dto.UserListResponse, error) {
	logger.Debug("listing users in service with pagination", zap.Int("page", req.Page), zap.Int("pageSize", req.PageSize))

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	users, total, err := s.userDAO.FindPage(ctx, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	pages := 0
	if req.PageSize > 0 {
		pages = int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}

	page := dto.PageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Pages:    pages,
	}

	return dto.NewUserListResponseFromEntities(users, page), nil
}

// FindByUserName 根据用户名查找用户
func (s *UserService) FindByUserName(ctx context.Context, userName string) (*entity.UserInfo, error) {
	logger.Debug("finding user by username in service", zap.String("username", userName))
	return s.userDAO.FindByUserName(ctx, userName)
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, password string, desKey string, jwtSecret string, jwtExpireHours int) (*dto.LoginResponse, error) {
	logger.Info("user login attempt", zap.String("username", username))

	// 查找用户
	user, err := s.userDAO.FindByUserName(ctx, username)
	if err != nil {
		logger.Error("failed to find user for login", zap.Error(err), zap.String("username", username))
		return nil, errors.New("登录失败")
	}

	if user == nil {
		logger.Warn("user not found for login", zap.String("username", username))
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户是否锁定
	if user.IsLocked == 1 {
		logger.Warn("locked user attempted login", zap.String("username", username))
		return nil, errors.New("用户已被锁定")
	}

	// DES解密数据库中的密码
	decryptedPassword, err := utils.DESDecrypt(user.Password, desKey)
	if err != nil {
		logger.Error("failed to decrypt password", zap.Error(err), zap.String("username", username))
		return nil, errors.New("登录失败")
	}

	// 比较密码（使用常量时间比较避免时序攻击）
	if !constantTimeCompare(password, decryptedPassword) {
		logger.Warn("invalid password for login", zap.String("username", username))
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.UserName, user.Role, jwtSecret, jwtExpireHours)
	if err != nil {
		logger.Error("failed to generate token", zap.Error(err), zap.String("username", username))
		return nil, errors.New("登录失败")
	}

	// 构建响应
	response := &dto.LoginResponse{
		Token: token,
		User:  dto.NewUserResponseFromEntity(user),
	}

	logger.Info("user login successful", zap.String("username", username), zap.Uint("user_id", user.ID))
	return response, nil
}

// Logout 用户退出
func (s *UserService) Logout(ctx context.Context) error {
	logger.Info("user logout")
	// 简单实现，客户端删除token即可
	return nil
}

// InitSuperAdmin 初始化超级管理员用户（如果不存在）
func (s *UserService) InitSuperAdmin(ctx context.Context, desKey string) error {
	const superAdminUsername = "superAdmin"
	const superAdminPassword = "superAdmin"

	logger.Info("checking superAdmin user", zap.String("username", superAdminUsername))

	// 检查用户是否已存在
	existingUser, err := s.userDAO.FindByUserName(ctx, superAdminUsername)
	if err != nil {
		logger.Error("failed to check superAdmin existence", zap.Error(err))
		return err
	}

	if existingUser != nil {
		logger.Info("superAdmin user already exists", zap.String("username", superAdminUsername))
		return nil
	}

	// 加密密码
	encryptedPassword, err := utils.DESEncrypt(superAdminPassword, desKey)
	if err != nil {
		logger.Error("failed to encrypt superAdmin password", zap.Error(err))
		return err
	}

	// 创建超级管理员用户
	user := &entity.UserInfo{
		UserName:  superAdminUsername,
		Password:  encryptedPassword,
		Role:      entity.RoleAdministrator,
		IsLocked:  0,
		ExtraInfo: nil,
	}

	// 保存到数据库
	if err := s.userDAO.Create(ctx, user); err != nil {
		logger.Error("failed to create superAdmin user", zap.Error(err))
		return err
	}

	logger.Info("superAdmin user created successfully",
		zap.String("username", superAdminUsername),
		zap.Uint("id", user.ID),
	)
	return nil
}

// constantTimeCompare 常量时间比较字符串，避免时序攻击
func constantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	result := 0
	for i := 0; i < len(a); i++ {
		result |= int(a[i]) ^ int(b[i])
	}
	return result == 0
}
