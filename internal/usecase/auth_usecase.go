package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"go-clean-arch-saas/internal/entity"
	"go-clean-arch-saas/internal/model"
	"go-clean-arch-saas/internal/model/converter"
	"go-clean-arch-saas/internal/repository"
	"go-clean-arch-saas/pkg/email"
	jwtPkg "go-clean-arch-saas/pkg/jwt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB                           *gorm.DB
	Log                          *logrus.Logger
	Validate                     *validator.Validate
	UserRepository               *repository.UserRepository
	OrganizationRepository       *repository.OrganizationRepository
	OrganizationMemberRepository *repository.OrganizationMemberRepository
	PlanRepository               *repository.PlanRepository
	SubscriptionRepository       *repository.SubscriptionRepository
	JWTService                   *jwtPkg.JWTService
	EmailService                 *email.EmailService
	BaseURL                      string
}

func NewAuthUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	userRepo *repository.UserRepository,
	orgRepo *repository.OrganizationRepository,
	orgMemberRepo *repository.OrganizationMemberRepository,
	planRepo *repository.PlanRepository,
	subRepo *repository.SubscriptionRepository,
	jwtService *jwtPkg.JWTService,
	emailService *email.EmailService,
	baseURL string,
) *AuthUseCase {
	return &AuthUseCase{
		DB:                           db,
		Log:                          logger,
		Validate:                     validate,
		UserRepository:               userRepo,
		OrganizationRepository:       orgRepo,
		OrganizationMemberRepository: orgMemberRepo,
		PlanRepository:               planRepo,
		SubscriptionRepository:       subRepo,
		JWTService:                   jwtService,
		EmailService:                 emailService,
		BaseURL:                      baseURL,
	}
}

func (u *AuthUseCase) Register(ctx context.Context, request *model.RegisterRequest) (*model.RegisterResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if email already exists
	count, err := u.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		u.Log.Warnf("Failed to count user by email: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if count > 0 {
		u.Log.Warnf("Email already exists: %s", request.Email)
		return nil, fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	// Create organization
	orgSlug := strings.ToLower(strings.ReplaceAll(request.OrganizationName, " ", "-"))
	orgID := uuid.New().String()

	organization := &entity.Organization{
		ID:        orgID,
		Name:      request.OrganizationName,
		Slug:      orgSlug,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	if err := u.OrganizationRepository.Create(tx, organization); err != nil {
		u.Log.Warnf("Failed to create organization: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Warnf("Failed to hash password: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate verification token
	verificationToken, err := generateVerificationToken()
	if err != nil {
		u.Log.Warnf("Failed to generate verification token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create user
	userID := uuid.New().String()
	user := &entity.User{
		ID:                userID,
		Name:              request.Name,
		Email:             request.Email,
		Password:          string(hashedPassword),
		SystemRole:        entity.SystemRoleUser, // Default to regular user
		EmailVerified:     false,
		VerificationToken: &verificationToken,
		OrganizationID:    orgID,
		CreatedAt:         time.Now().UnixMilli(),
		UpdatedAt:         time.Now().UnixMilli(),
	}

	if err := u.UserRepository.Create(tx, user); err != nil {
		u.Log.Warnf("Failed to create user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Add user as owner of organization
	orgMember := &entity.OrganizationMember{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           entity.OrgRoleOwner, // Use constant instead of hardcoded string
		JoinedAt:       time.Now().UnixMilli(),
	}

	if err := u.OrganizationMemberRepository.Create(tx, orgMember); err != nil {
		u.Log.Warnf("Failed to create organization member: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Get free plan
	freePlan := new(entity.Plan)
	if err := u.PlanRepository.FindBySlug(tx, freePlan, "free"); err != nil {
		u.Log.Warnf("Failed to find free plan: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create subscription with free plan
	subscription := &entity.Subscription{
		ID:                 uuid.New().String(),
		OrganizationID:     orgID,
		PlanID:             freePlan.ID,
		Status:             "active",
		CurrentPeriodStart: time.Now().UnixMilli(),
		CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0).UnixMilli(), // 1 month
		CreatedAt:          time.Now().UnixMilli(),
		UpdatedAt:          time.Now().UnixMilli(),
	}

	if err := u.SubscriptionRepository.Create(tx, subscription); err != nil {
		u.Log.Warnf("Failed to create subscription: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Send verification email (non-blocking, don't fail if email fails)
	go func() {
		if err := u.EmailService.SendVerificationEmail(user.Email, user.Name, verificationToken, u.BaseURL); err != nil {
			u.Log.Warnf("Failed to send verification email to %s: %+v", user.Email, err)
		} else {
			u.Log.Infof("Verification email sent to %s", user.Email)
		}
	}()

	return &model.RegisterResponse{
		User:         *converter.UserToResponse(user),
		Organization: *converter.OrganizationToResponse(organization),
	}, nil
}

func (u *AuthUseCase) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find user by email
	user := new(entity.User)
	if err := u.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		u.Log.Warnf("Failed to find user by email: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.Log.Warnf("Invalid password: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	// Generate access token (JWT)
	accessToken, err := u.JWTService.GenerateAccessToken(user.ID, user.Email, user.OrganizationID)
	if err != nil {
		u.Log.Warnf("Failed to generate access token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate refresh token (UUID)
	refreshToken := uuid.New().String()
	refreshTokenExpiresAt := time.Now().Add(u.JWTService.GetRefreshTokenExpiration()).UnixMilli()

	// Save refresh token to database
	user.RefreshToken = refreshToken
	user.RefreshTokenExpiresAt = refreshTokenExpiresAt

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed to update user refresh token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(u.JWTService.GetAccessTokenExpiration().Seconds()),
		TokenType:    "Bearer",
		User:         *converter.UserToResponse(user),
	}, nil
}

func (u *AuthUseCase) Refresh(ctx context.Context, request *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find user by refresh token
	user := new(entity.User)
	if err := u.UserRepository.FindByRefreshToken(tx, user, request.RefreshToken); err != nil {
		u.Log.Warnf("Failed to find user by refresh token: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	// Check if refresh token is expired
	if user.RefreshTokenExpiresAt < time.Now().UnixMilli() {
		u.Log.Warnf("Refresh token expired for user: %s", user.ID)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Refresh token expired")
	}

	// Generate new access token
	accessToken, err := u.JWTService.GenerateAccessToken(user.ID, user.Email, user.OrganizationID)
	if err != nil {
		u.Log.Warnf("Failed to generate access token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(u.JWTService.GetAccessTokenExpiration().Seconds()),
		TokenType:   "Bearer",
	}, nil
}

func (u *AuthUseCase) Logout(ctx context.Context, userID string) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Find user
	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, userID); err != nil {
		u.Log.Warnf("Failed to find user: %+v", err)
		return fiber.ErrNotFound
	}

	// Clear refresh token
	user.RefreshToken = ""
	user.RefreshTokenExpiresAt = 0

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed to update user: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *AuthUseCase) VerifyToken(ctx context.Context, token string) (*model.Auth, error) {
	// Remove "Bearer " prefix if exists
	token = strings.Replace(token, "Bearer ", "", 1)

	// Validate JWT token
	claims, err := u.JWTService.ValidateToken(token)
	if err != nil {
		u.Log.Warnf("Invalid JWT token: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	return &model.Auth{
		UserID:         claims.UserID,
		Email:          claims.Email,
		OrganizationID: claims.OrganizationID,
	}, nil
}

func (u *AuthUseCase) VerifyEmail(ctx context.Context, request *model.VerifyEmailRequest) (*model.VerifyEmailResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find user by verification token
	user := new(entity.User)
	if err := u.UserRepository.FindByVerificationToken(tx, user, request.Token); err != nil {
		u.Log.Warnf("Invalid verification token: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid or expired verification token")
	}

	// Check if already verified
	if user.EmailVerified {
		return &model.VerifyEmailResponse{
			Message: "Email already verified",
		}, nil
	}

	// Update user - mark email as verified
	now := time.Now().UnixMilli()
	user.EmailVerified = true
	user.EmailVerifiedAt = &now
	user.VerificationToken = nil // Clear token after verification

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed to update user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	u.Log.Infof("Email verified for user: %s (%s)", user.ID, user.Email)

	return &model.VerifyEmailResponse{
		Message: "Email verified successfully",
	}, nil
}

func (u *AuthUseCase) ResendVerification(ctx context.Context, request *model.ResendVerificationRequest) (*model.ResendVerificationResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find user by email
	user := new(entity.User)
	if err := u.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		// Don't reveal if email exists or not for security
		u.Log.Warnf("User not found for resend verification: %s", request.Email)
		return &model.ResendVerificationResponse{
			Message: "If the email exists, a verification link has been sent",
		}, nil
	}

	// Check if already verified
	if user.EmailVerified {
		return &model.ResendVerificationResponse{
			Message: "Email already verified",
		}, nil
	}

	// Generate new verification token
	verificationToken, err := generateVerificationToken()
	if err != nil {
		u.Log.Warnf("Failed to generate verification token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update verification token
	user.VerificationToken = &verificationToken
	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed to update user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Send verification email (non-blocking)
	go func() {
		if err := u.EmailService.SendVerificationEmail(user.Email, user.Name, verificationToken, u.BaseURL); err != nil {
			u.Log.Warnf("Failed to send verification email to %s: %+v", user.Email, err)
		} else {
			u.Log.Infof("Verification email resent to %s", user.Email)
		}
	}()

	return &model.ResendVerificationResponse{
		Message: "If the email exists, a verification link has been sent",
	}, nil
}

// generateVerificationToken generates a random verification token
func generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
