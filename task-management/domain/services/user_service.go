package services

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *requests.RegisterRequest) (*responses.UserWithTokenResponse, *errutils.Error)
	Login(ctx context.Context, req *requests.LoginRequest) (*responses.UserWithTokenResponse, *errutils.Error)
	FindUserByEmail(ctx context.Context, email string) (*responses.UserResponse, *errutils.Error)
	Search(ctx context.Context, req *requests.SearchUserParams, searcherUserId string) (*responses.ListUserResponse, *errutils.Error)
	SetupFirstUser(ctx context.Context, req *requests.RegisterRequest) (*responses.UserWithTokenResponse, *errutils.Error)
}

type userServiceImpl struct {
	config            *config.Config
	userRepo          repositories.UserRepository
	globalSettingRepo repositories.GlobalSettingRepository
}

func NewUserService(
	config *config.Config,
	userRepo repositories.UserRepository,
	globalSettingRepo repositories.GlobalSettingRepository,
) UserService {
	return &userServiceImpl{
		config:            config,
		userRepo:          userRepo,
		globalSettingRepo: globalSettingRepo,
	}
}

func (u *userServiceImpl) generateJWT(user *models.User, expireAt time.Time) (string, *errutils.Error) {
	claims := models.UserCustomClaims{
		ID:          user.ID.Hex(),
		FullName:    user.FullName,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.Hex(),
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(u.config.JWT.AccessTokenSecret))
	if err != nil {
		return "", errutils.NewError(exceptions.ErrInternalError, errutils.InternalError)
	}

	return tokenString, nil
}

func (u *userServiceImpl) Register(ctx context.Context, req *requests.RegisterRequest) (*responses.UserWithTokenResponse, *errutils.Error) {
	// Check if email already exists
	existsUser, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError)
	}
	if existsUser != nil {
		return nil, errutils.NewError(exceptions.ErrUserAlreadyExists, errutils.BadRequest)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}
	req.Password = string(hashedPassword)

	// Generate profile url
	fullName := strings.Trim(req.FullName, " ")
	nameParts := strings.Split(fullName, " ")
	var profileUrl = "https://ui-avatars.com/api/?name="
	if len(nameParts) == 1 {
		profileUrl += nameParts[0]
	} else {
		profileUrl += nameParts[0] + "+" + nameParts[1]
	}

	user := &repositories.CreateUserRequest{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		DisplayName:  fullName,
		ProfileUrl:   profileUrl,
	}
	createdUser, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError)
	}

	// Generate JWT token
	expireAt := time.Now().Add(time.Hour * 120)

	token, tokenErr := u.generateJWT(createdUser, expireAt)
	if tokenErr != nil {
		return nil, tokenErr
	}

	res := &responses.UserWithTokenResponse{
		UserResponse: responses.UserResponse{
			ID:          createdUser.ID.Hex(),
			Email:       createdUser.Email,
			FullName:    createdUser.FullName,
			DisplayName: createdUser.DisplayName,
			ProfileUrl:  createdUser.ProfileUrl,
			CreatedAt:   createdUser.CreatedAt,
			UpdatedAt:   createdUser.UpdatedAt,
		},
		Token:         token,
		TokenExpireAt: expireAt,
	}
	return res, nil
}

func (u *userServiceImpl) Login(ctx context.Context, req *requests.LoginRequest) (*responses.UserWithTokenResponse, *errutils.Error) {
	// Find user by email
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError)
	}

	if user == nil {
		return nil, errutils.NewError(exceptions.ErrInvalidCredentials, errutils.Unauthorized)
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInvalidCredentials, errutils.Unauthorized)
	}

	// Generate JWT token
	expireAt := time.Now().Add(time.Hour * 120)

	token, tokenErr := u.generateJWT(user, expireAt)
	if tokenErr != nil {
		return nil, tokenErr
	}

	res := &responses.UserWithTokenResponse{
		UserResponse: responses.UserResponse{
			ID:          user.ID.Hex(),
			Email:       user.Email,
			FullName:    user.FullName,
			ProfileUrl:  user.ProfileUrl,
			DisplayName: user.DisplayName,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		Token:         token,
		TokenExpireAt: expireAt,
	}
	return res, nil
}

func (u *userServiceImpl) FindUserByEmail(ctx context.Context, email string) (*responses.UserResponse, *errutils.Error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError)
	}

	if user == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.NotFound)
	}

	res := &responses.UserResponse{
		ID:          user.ID.Hex(),
		Email:       user.Email,
		FullName:    user.FullName,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return res, nil
}

func validateSearchUserPaginationRequestSortBy(sortBy string) bool {
	switch sortBy {
	case constant.UserFieldEmail, constant.UserFieldFullName, constant.UserFieldDisplayName:
		return true
	}
	return false
}

func (u *userServiceImpl) Search(ctx context.Context, req *requests.SearchUserParams, searcherUserId string) (*responses.ListUserResponse, *errutils.Error) {
	if req.PaginationRequest != nil {
		if req.PaginationRequest.Page <= 0 {
			req.PaginationRequest.Page = 1
		}
		if req.PaginationRequest.PageSize <= 0 {
			req.PaginationRequest.PageSize = 100
		}
		if req.PaginationRequest.SortBy == "" || !validateSearchUserPaginationRequestSortBy(req.PaginationRequest.SortBy) {
			req.PaginationRequest.SortBy = constant.UserFieldEmail
		}
		if req.PaginationRequest.Order == "" {
			req.PaginationRequest.Order = constant.ASC
		}
	} else {
		req.PaginationRequest = &requests.PaginationRequest{
			Page:     1,
			PageSize: 100,
			SortBy:   constant.UserFieldEmail,
			Order:    constant.ASC,
		}
	}

	users, totalUser, err := u.userRepo.Search(ctx, &repositories.SearchUserRequest{
		Keyword:           req.Keyword,
		PaginationRequest: repositories.PaginationRequest{Page: req.PaginationRequest.Page, PageSize: req.PaginationRequest.PageSize, SortBy: req.PaginationRequest.SortBy, Order: req.PaginationRequest.Order},
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalError).WithDebugMessage(err.Error())
	}

	res := &responses.ListUserResponse{
		Users: make([]responses.UserResponse, 0),
		PaginationResponse: responses.PaginationResponse{
			Page:      req.PaginationRequest.Page,
			PageSize:  req.PaginationRequest.PageSize,
			TotalPage: int(math.Ceil(float64(totalUser) / float64(req.PaginationRequest.PageSize))),
			TotalItem: int(totalUser),
		},
	}

	for _, user := range users {
		if user.ID.Hex() == searcherUserId {
			res.PaginationResponse.TotalItem--
			continue
		}

		res.Users = append(res.Users, responses.UserResponse{
			ID:          user.ID.Hex(),
			Email:       user.Email,
			FullName:    user.FullName,
			DisplayName: user.DisplayName,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}

	return res, nil
}

func (u *userServiceImpl) SetupFirstUser(ctx context.Context, req *requests.RegisterRequest) (*responses.UserWithTokenResponse, *errutils.Error) {
	// Check is setup
	isSetupOwner, err := u.globalSettingRepo.GetByKey(ctx, constant.GlobalSettingKeyIsSetupOwner)
	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	if isSetupOwner == nil {
		err := u.globalSettingRepo.Set(ctx, &models.GlobalSetting{
			Key:   constant.GlobalSettingKeyIsSetupOwner,
			Type:  models.GlobalSettingTypeBool,
			Value: false,
		})

		if err != nil {
			return nil, errutils.NewError(err, errutils.InternalServerError)
		}
	}

	if isSetupOwner.Value.(bool) {
		return nil, errutils.NewError(exceptions.ErrOwnerAlreadySetup, errutils.BadRequest)
	}

	newUser, regErr := u.Register(ctx, req)
	if regErr != nil {
		return nil, regErr
	}

	// Set is setup owner
	err = u.globalSettingRepo.Set(ctx, &models.GlobalSetting{
		Key:   constant.GlobalSettingKeyIsSetupOwner,
		Type:  models.GlobalSettingTypeBool,
		Value: true,
	})

	if err != nil {
		return nil, errutils.NewError(err, errutils.InternalServerError)
	}

	return newUser, nil
}
