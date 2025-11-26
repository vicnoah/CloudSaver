package repository

import (
	"cloudsaver/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUUID 根据UUID查找用户
func (r *UserRepository) FindByUUID(uuid string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(uuid string) error {
	return r.db.Where("uuid = ?", uuid).Delete(&model.User{}).Error
}

// Count 统计用户数量
func (r *UserRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
