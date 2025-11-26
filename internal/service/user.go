package service

import (
	"errors"

	"cloudsaver/internal/api/response"
	"cloudsaver/internal/model"
	"cloudsaver/internal/pkg/utils"
	"cloudsaver/internal/repository"

	"github.com/google/uuid"
)

// UserService 用户服务
type UserService struct {
	userRepo       *repository.UserRepository
	settingRepo    *repository.SettingRepository
	jwtSecret      string
	jwtExpireHours int
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repository.UserRepository, settingRepo *repository.SettingRepository, jwtSecret string, jwtExpireHours int) *UserService {
	return &UserService{
		userRepo:       userRepo,
		settingRepo:    settingRepo,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

// Register 用户注册
func (s *UserService) Register(username, password string) (*response.UserLoginResponse, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.FindByUsername(username)
	if existingUser != nil {
		return nil, errors.New("用户名已被使用")
	}

	// 判断是否是第一个用户（自动成为管理员）
	userCount, err := s.userRepo.Count()
	if err != nil {
		return nil, errors.New("查询用户数量失败")
	}

	role := 0 // 默认普通用户
	if userCount == 0 {
		role = 1 // 第一个用户为管理员
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		UUID:     uuid.New().String(),
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 生成Token
	token, err := utils.GenerateToken(user.UUID, user.Role, s.jwtSecret, s.jwtExpireHours)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	return &response.UserLoginResponse{
		User: response.UserInfo{
			UUID:     user.UUID,
			Username: user.Username,
			Role:     user.Role,
		},
		Token: token,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*response.UserLoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成Token
	token, err := utils.GenerateToken(user.UUID, user.Role, s.jwtSecret, s.jwtExpireHours)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	return &response.UserLoginResponse{
		User: response.UserInfo{
			UUID:     user.UUID,
			Username: user.Username,
			Role:     user.Role,
		},
		Token: token,
	}, nil
}
