package auth

import (
	"backend-poc/backoffice/domain"
	"errors"
	"fmt"
	"time"

	"github.com/dlclark/regexp2"
	"golang.org/x/crypto/bcrypt"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type User struct {
	domain.Entity
	Email                     string     `json:"email" gorm:"type:varchar(255);not null;index:idx_users_email"`
	Role                      string     `json:"role" gorm:"type:varchar(45);not null"`
	Name                      string     `json:"name" gorm:"type:varchar(500);not null"`
	PasswordHash              string     `json:"-" gorm:"type:varchar(60);not null"`
	Photo                     *string    `json:"photo" gorm:"type:varchar(500)"`
	EmailVerified             bool       `json:"email_verified"`
	HasToChangePassword       bool       `json:"has_to_change_password"`
	Document                  *string    `json:"document" gorm:"type:varchar(255)"`
	Phone                     *string    `json:"phone" gorm:"type:varchar(45)"`
	PhoneVerified             bool       `json:"phone_verified" gorm:"type:boolean;default:false"`
	EnabledTOTP               bool       `json:"enabled_totp" gorm:"type:boolean;default:false"`
	SecretTOTP                string     `json:"secret_totp" gorm:"type:varchar(100)"`
	Country                   *string    `json:"country" gorm:"type:varchar(50)"`
	State                     *string    `json:"state" gorm:"type:varchar(50)"`
	City                      *string    `json:"city" gorm:"type:varchar(50)"`
	Neighborhood              *string    `json:"neighborhood" gorm:"type:varchar(50)"`
	Street                    *string    `json:"street" gorm:"type:varchar(50)"`
	Number                    *string    `json:"number" gorm:"type:varchar(50)"`
	Complement                *string    `json:"complement" gorm:"type:varchar(255)"`
	ZipCode                   *string    `json:"zip_code" gorm:"type:varchar(15)"`
	Gender                    *Gender    `json:"gender" gorm:"type:varchar(10)"`
	BirthDate                 *time.Time `json:"birth_date" gorm:"type:date"`
	Nationality               *string    `json:"nationality" gorm:"type:varchar(50)"`
	Currency                  *string    `json:"currency" gorm:"type:varchar(50)"`
	CustomWithdrawLimitAmount int        `json:"custom_withdraw_limit_amount" gorm:"not null;default:0"`
	CustomWithdrawLimitTimes  int        `json:"custom_withdraw_limit_times" gorm:"not null;default:0"`

	WithdrawLimitId *uint `json:"withdraw_limit_id" gorm:"foreignKey:WithdrawLimitId;references:Id"`

	Roles []Role `json:"-" gorm:"many2many:user_role"`
}

func (u *User) SetPassword(plainTextPassword, regex string, cost int) error {
	ok, err := regexp2.MustCompile(regex, regexp2.None).MatchString(plainTextPassword)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(fmt.Sprintf("password must match regex %s", regex))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), cost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(plainTextPassword string) bool {
	if u.PasswordHash == "" {
		return false
	}
	if plainTextPassword == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainTextPassword)) == nil
}

func (u *User) HasPassword() bool {
	return u.PasswordHash != ""
}

func (u *User) verifyPhone() {
	u.PhoneVerified = true

	u.UpdatedAt = time.Now()
}

// TODO veirificar implementacao dos novos campos na hora da criacao do usuario.
func NewUser(name, email string, phone, document *string, defaultWithdrawLimitId *uint) *User {
	return &User{
		Name:            name,
		Email:           email,
		Phone:           phone,
		Document:        document,
		Role:            "user",
		WithdrawLimitId: defaultWithdrawLimitId,
	}
}
