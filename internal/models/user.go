package models

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ogidi/church-media/ent"
	"github.com/ogidi/church-media/ent/contactprofile"
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
		user.And(user.Or(user.EmailEQ(dto.EmailUsername), user.UsernameEQ(dto.EmailUsername)))).WithContactProfile().First(ctx)
	if err != nil {
		return nil, EmailUsernameInvalidError
	}

	if userExist == nil {
		return nil, EmailUsernameInvalidError
	}
	if userExist.State == user.StatePENDING || userExist.State == user.StateFRESH {
		return nil, AccountNotVerifiedError
	}

	passwordCorrect, err := m.PasswordCheck(ctx, userExist.Password, dto.Password)
	if err != nil {
		if errors.Is(err, InvalidCredentialsError) {
			return nil, InvalidCredentialsError
		} else {
			return nil, err
		}
	}
	if !passwordCorrect {
		return nil, InvalidCredentialsError
	}
	return userExist, nil
}

func (m *UserModel) PasswordCheck(ctx context.Context, hashed string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, InvalidCredentialsError
		} else {
			return false, err
		}
	}
	return true, nil
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
	if userData == nil {
		return nil, UserNotFoundError
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

func (m *UserModel) UpdateProfile(ctx context.Context, dto *UpdateProfileRequest, id int) error {
	userExist, err := m.UserExistById(ctx, id)
	if err != nil {
		return err
	}
	if !userExist {
		return UserNotFoundError
	}

	tx, err := m.DB.Tx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// user update
	userUpdate := tx.User.UpdateOneID(id).
		SetDepartment(dto.Department).
		SetUpdatedAt(time.Now())
	if dto.NewPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		userUpdate = userUpdate.SetPassword(string(hashed))
	}
	_, err = userUpdate.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return ConstraintError
		}
	}

	contactUpdate := tx.ContactProfile.Update().
		Where(contactprofile.HasUserWith(user.IDEQ(id))).
		SetFirstName(dto.FirstName).
		SetSurname(dto.Surname).
		SetPhoneNumber(dto.PhoneNumber).
		SetDateOfBirth(dto.DateOfBirth).
		SetGender(contactprofile.Gender(dto.Gender)).
		SetOccupation(dto.Occupation).
		SetMaritalStatus(dto.MaritalStatus).
		SetAddress(dto.Address)

	if dto.ProfilePicture != "" {
		contactUpdate = contactUpdate.SetProfilePicture(dto.ProfilePicture)
	}
	_, err = contactUpdate.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return UserNotFoundError
		}
		return err
	}
	return nil
}

func (m *UserModel) GetAllInvitedAdmins(ctx context.Context, page, pageSize int, filters Filters, sort Sort) ([]*ent.User, Pagination, error) {
	query := m.DB.User.Query().
		Where(user.RoleIn("admin", "content_admin", "secretary", "pastor", "deacon")) // Only show admin-type roles

	// Apply filters
	if filters.Role != "" {
		query = query.Where(user.RoleEQ(user.Role(filters.Role)))
	}
	if filters.Status != "" {
		query = query.Where(user.StateEQ(user.State(filters.Status)))
	}
	if filters.Search != "" {
		query = query.Where(
			user.Or(
				user.EmailContains(filters.Search),
				user.HasContactProfileWith(
					contactprofile.Or(
						contactprofile.FirstNameContains(filters.Search),
						contactprofile.SurnameContains(filters.Search),
					),
				),
			),
		)
	}

	// Apply sorting
	switch sort.Field {
	case "name":
		direction := sql.OrderAsc()
		if sort.Order == "desc" {
			direction = sql.OrderDesc()
		}
		query = query.Order(
			user.ByContactProfileField(contactprofile.FieldFirstName, direction),
			user.ByContactProfileField(contactprofile.FieldSurname, direction),
		)
	case "email":
		direction := sql.OrderAsc()
		if sort.Order == "desc" {
			direction = sql.OrderDesc()
		}
		query = query.Order(user.ByEmail(direction))
	case "date":
		direction := sql.OrderAsc()
		if sort.Order == "desc" {
			direction = sql.OrderDesc()
		}
		query = query.Order(user.ByCreatedAt(direction))
	default:
		// Default sorting by created_at desc
		query = query.Order(user.ByCreatedAt(sql.OrderDesc()))
	}

	// Count total items for pagination
	total, err := query.Count(ctx)
	if err != nil {
		return nil, Pagination{}, err
	}

	// Calculate pagination
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	pagination := Pagination{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		TotalItems:  total,
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Include contact profile
	admins, err := query.WithContactProfile().All(ctx)
	if err != nil {
		return nil, Pagination{}, err
	}

	return admins, pagination, nil
}

//func (m *UserModel) UpdateInvitedAdmin(ctx context.Context, id int, dto *InviteDto) error {
//	token, err := m.GenerateToken(dto.Email, dto.ExpiresAt)
//	if err != nil {
//		return err
//	}
//	updatedUser, err = m.DB.User.UpdateOneID(id).
//		SetRegistrationToken(token).
//		SetRole(dto.Role).Save(ctx)
//	if err != nil {
//		return err
//	}
//	_, err = m.DB.ContactProfile.Update().Where(usercon)
//	return nil
//}

func (m *UserModel) UpdateLastLogin(ctx context.Context, id int) error {
	_, err := m.DB.User.UpdateOneID(id).SetLastLogin(time.Now()).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UpdateUserRole(ctx context.Context, id int, role string) error {
	_, err := m.DB.User.UpdateOneID(id).SetRole(user.Role(role)).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) DeleteUser(ctx context.Context, id int) error {
	err := m.DB.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) LastUpdated(ctx context.Context, id int) error {
	_, err := m.DB.User.UpdateOneID(id).SetUpdatedAt(time.Now()).Save(ctx)
	if err != nil {
		return err
	}
	return nil
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

func (m *UserModel) UpdateResetToken(ctx context.Context, email string) (string, error) {
	code, err := m.GenerateToken(email, 3)
	err = m.DB.User.Update().Where(user.EmailEQ(email)).SetResetToken(code).Exec(ctx)
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
