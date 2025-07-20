package auth

import (
	"DistanceTrackerServer/models"
	"DistanceTrackerServer/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"unicode"
)

var (
	validateEmailFunc    = ValidateEmail
	validatePasswordFunc = ValidatePassword
	validateFirstName    = ValidateFirstName
)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Basic email validation
	emailAddress, err := mail.ParseAddress(email)
	if err != nil || emailAddress.Address != email {
		return fmt.Errorf("email address %s failed validation", email)
	}
	return nil
}

func ValidatePassword(password string, confirmPassword string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if confirmPassword == "" {
		return fmt.Errorf("confirm password is required")
	}
	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasNumber = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		} else if unicode.IsSpace(char) {
			return fmt.Errorf("password cannot contain spaces")
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

func ValidateFirstName(name string) error {
	if name == "" {
		return fmt.Errorf("first name is required")
	}
	expression := regexp.MustCompile(`[:;()\[\]{}|\\/]+`)
	if expression.MatchString(name) {
		return fmt.Errorf("first name contains dangerous characters: %s", name)
	}
	return nil
}

func ValidateRegistration(user models.UserRegister) error {
	if err := validateEmailFunc(user.Email); err != nil {
		return err
	}
	if err := validateFirstName(user.FirstName); err != nil {
		return err
	}
	if err := validatePasswordFunc(user.Password, user.ConfirmPassword); err != nil {
		return err
	}
	return nil
}

func Register(ctx *gin.Context, user models.UserRegister) (int, error) {
	//sugar, err := utils.SugarFromContext(ctx)
	//if err != nil {
	//	return fmt.Errorf("failed to retrieve logger from context: %s", err)
	//}
	password := ([]byte)(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %s", err)
	}

	dbConn, err := utils.DBConnFromContext(ctx)
	if err != nil {
		return 0, err
	}
	res, err := dbConn.Exec("INSERT INTO users(id, email, name, password) VALUES(NULL, ?, ?, ?)",
		user.Email, user.FirstName, hashedPassword)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return 0, fmt.Errorf("email already exists")
		}
		return 0, fmt.Errorf("failed to insert registration into db")
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id from db insert: %s", err)
	}

	tokenString, err := createToken(user.Email)
	if err != nil {
		return 0, fmt.Errorf("failed to create token: %w", err)
	}
	ctx.SetCookie("token", tokenString, 3600, "/", "", false, true)

	return int(userID), nil
}
