package impl

import (
	"context"
	"errors"

	dao "robot_scheduler/internal/dao/interfaces"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserDAOImpl struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) dao.UserDAO {
	return &UserDAOImpl{db: db}
}

func (d *UserDAOImpl) Create(ctx context.Context, user *entity.UserInfo) error {
	logger.Info("creating user", zap.String("username", user.UserName))

	if err := d.db.WithContext(ctx).Create(user).Error; err != nil {
		logger.Error("failed to create user", zap.Error(err), zap.Any("user", user))
		return err
	}

	logger.Info("user created successfully", zap.Uint("id", user.ID), zap.String("username", user.UserName))
	return nil
}

func (d *UserDAOImpl) Update(ctx context.Context, user *entity.UserInfo) error {
	logger.Info("updating user", zap.Uint("id", user.ID))

	result := d.db.WithContext(ctx).Save(user)
	if err := result.Error; err != nil {
		logger.Error("failed to update user", zap.Error(err), zap.Uint("id", user.ID))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("user not found for update", zap.Uint("id", user.ID))
		return errors.New("user not found")
	}

	logger.Info("user updated successfully", zap.Uint("id", user.ID))
	return nil
}

func (d *UserDAOImpl) Delete(ctx context.Context, id uint) error {
	logger.Info("deleting user", zap.Uint("id", id))

	result := d.db.WithContext(ctx).Delete(&entity.UserInfo{}, id)
	if err := result.Error; err != nil {
		logger.Error("failed to delete user", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("user not found for deletion", zap.Uint("id", id))
		return errors.New("user not found")
	}

	logger.Info("user deleted successfully", zap.Uint("id", id))
	return nil
}

func (d *UserDAOImpl) FindByID(ctx context.Context, id uint) (*entity.UserInfo, error) {
	logger.Debug("finding user by id", zap.Uint("id", id))

	var user entity.UserInfo
	err := d.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("user not found", zap.Uint("id", id))
			return nil, nil
		}
		logger.Error("failed to find user by id", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	logger.Debug("user found", zap.Uint("id", id))
	return &user, nil
}

func (d *UserDAOImpl) FindAll(ctx context.Context) ([]*entity.UserInfo, error) {
	logger.Debug("finding all users")

	var users []*entity.UserInfo
	err := d.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		logger.Error("failed to find all users", zap.Error(err))
		return nil, err
	}

	logger.Debug("found users", zap.Int("count", len(users)))
	return users, nil
}

// FindPage 分页查询用户
func (d *UserDAOImpl) FindPage(ctx context.Context, offset, limit int) ([]*entity.UserInfo, int64, error) {
	logger.Debug("finding users with pagination", zap.Int("offset", offset), zap.Int("limit", limit))

	var (
		users []*entity.UserInfo
		total int64
	)

	db := d.db.WithContext(ctx).Model(&entity.UserInfo{})

	if err := db.Count(&total).Error; err != nil {
		logger.Error("failed to count users for pagination", zap.Error(err))
		return nil, 0, err
	}

	if total == 0 {
		return []*entity.UserInfo{}, 0, nil
	}

	if err := db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		logger.Error("failed to find users with pagination", zap.Error(err))
		return nil, 0, err
	}

	logger.Debug("found users with pagination", zap.Int("count", len(users)), zap.Int64("total", total))
	return users, total, nil
}

func (d *UserDAOImpl) FindByUserName(ctx context.Context, userName string) (*entity.UserInfo, error) {
	logger.Debug("finding user by username", zap.String("username", userName))

	var user entity.UserInfo
	err := d.db.WithContext(ctx).Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("user not found by username", zap.String("username", userName))
			return nil, nil
		}
		logger.Error("failed to find user by username", zap.Error(err), zap.String("username", userName))
		return nil, err
	}

	logger.Debug("user found by username", zap.String("username", userName))
	return &user, nil
}