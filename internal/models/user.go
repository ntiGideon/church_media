package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/user"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type UserModel struct {
	DB *ent.Client
}

func (m *UserModel) Login(ctx context.Context, dto *LoginDto) (*ent.User, error) {
	userExist, err := m.DB.User.Query().Where(
		user.And(user.Or(user.EmailEQ(dto.EmailUsername), user.UsernameEQ(dto.EmailUsername)),
			user.StateEQ(user.StateVERIFIED))).First(ctx)
	if err != nil {
		return nil, EmailUsernameInvalidError
	}

	if userExist == nil {
		return nil, EmailUsernameInvalidError
	}
	if userExist.State == user.StatePENDING || userExist.State == user.StateFRESH {
		return nil, AccountNotVerifiedError
	}

	err = bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(dto.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, InvalidCredentialsError
		} else {
			return nil, err
		}
	}
	return userExist, nil
}

func (m *UserModel) Register(ctx context.Context, dto *RegisterDto) error {
	usernameExist, err := m.UsernameExists(ctx, dto.Username)
	if err != nil {
		return err
	}
	if usernameExist {
		return UsernameExistsError
	}

	tokenExist, err := m.TokenExists(ctx, dto.RegistrationToken)
	if err != nil {
		return err
	}
	if !tokenExist {
		return RegistrationTokenInvalidError
	}

	getUser, err := m.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}
	//if getUser.RegistrationToken != &dto.RegistrationToken {
	//	return EmailExistsError
	//}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return BcryptError
	}

	_, err = m.DB.User.Update().Where(user.EmailEQ(getUser.Email)).
		SetState(user.StateFRESH).
		SetUsername(dto.Username).
		SetPassword(string(hashedPassword)).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UsernameExists(ctx context.Context, username string) (bool, error) {
	usernameExists, err := m.DB.User.Query().Where(user.UsernameEQ(username)).Exist(ctx)
	return usernameExists, err
}

func (m *UserModel) EmailExists(ctx context.Context, email string) (bool, error) {
	emailExists, err := m.DB.User.Query().Where(user.EmailEQ(email)).Exist(ctx)
	return emailExists, err
}

func (m *UserModel) GetUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	userEmail, err := m.DB.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return userEmail, nil
}

func (m *UserModel) GetUserById(ctx context.Context, id int) (*ent.User, error) {
	userData, err := m.DB.User.Query().Where(user.IDEQ(id)).WithContactProfile().Only(ctx)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (m *UserModel) TokenExists(ctx context.Context, token string) (bool, error) {
	tokenExist, err := m.DB.User.Query().Where(user.RegistrationTokenEQ(token)).Exist(ctx)
	return tokenExist, err
}

func (m *UserModel) UserExistById(ctx context.Context, id int) (bool, error) {
	userExists, _ := m.DB.User.Query().Where(user.IDEQ(id)).Exist(ctx)
	if userExists {
		return true, nil
	}
	return false, nil
}

func (m *UserModel) InviteUser(ctx context.Context, dto *InviteDto) (string, string, error) {
	emailExist, err := m.EmailExists(ctx, dto.Email)
	if err != nil {
		return "", "", err
	}
	if emailExist {
		return "", "", EmailExistsError
	}

	tx, err := m.DB.Tx(ctx)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	token, err := m.GenerateToken(dto.Email, dto.ExpiresAt)
	if err != nil {
		return "", "", err
	}

	userInvite, err := tx.User.Create().
		SetState(user.StatePENDING).
		SetEmail(dto.Email).
		SetRole(dto.Role).
		SetRegistrationToken(token).
		Save(ctx)

	contact, err := tx.ContactProfile.Create().
		SetFirstName(dto.Firstname).
		SetSurname(dto.Lastname).
		SetUser(userInvite).
		Save(ctx)

	if err != nil {
		return "", "", err
	}

	return *userInvite.RegistrationToken, contact.FirstName, nil
}

func (m *UserModel) GenerateToken(email string, hours int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := &jwt.RegisteredClaims{
		Issuer:    "church-admin-app",
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *UserModel) VerifyToken(tokenString, expectedEmail string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return false, TokenValidationError
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Subject != expectedEmail {
			return false, TokenMissMatchError
		}

		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return false, TokenExpiredError
		}

		return true, nil
	}

	return false, fmt.Errorf("invalid token claims")
}

func (m *UserModel) UpdateVerificationToken(ctx context.Context, email string) (string, error) {
	code, err := m.GenerateToken(email, 3)
	err = m.DB.User.Update().Where(user.EmailEQ(email)).SetVerifyToken(code).Exec(ctx)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (m *UserModel) ActivateAccount(ctx context.Context, code string) error {
	codeExist, err := m.DB.User.Query().Where(user.And(user.VerifyToken(code), user.StateEQ(user.StateFRESH))).Exist(ctx)
	if err != nil {
		return err
	}
	if !codeExist {
		return WrongVerficationTokenError
	}
	err = m.DB.User.Update().Where(user.VerifyToken(code)).SetState(user.StateVERIFIED).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
