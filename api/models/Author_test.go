package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "mypassword"

	// Тестуємо хешування та перевірку пароля
	hashedPassword, err := Hash(password)
	assert.NoError(t, err)

	err = VerifyPassword(string(hashedPassword), password)
	assert.NoError(t, err)

	// Тестуємо невірний пароль
	wrongPassword := "wrongpassword"
	err = VerifyPassword(string(hashedPassword), wrongPassword)
	assert.Error(t, err)
}

func TestPrepareAuthor(t *testing.T) {
	author := &Author{
		Nickname: "  john_doe ",
		Email:    "  john@example.com  ",
	}

	author.Prepare()

	assert.Equal(t, "john_doe", author.Nickname)
	assert.Equal(t, "john@example.com", author.Email)
}

func TestValidate(t *testing.T) {
	// Створюємо тестовий об'єкт автора
	author := &Author{
		Nickname: "testuser",
		Password: "testpassword",
		Email:    "test@example.com",
	}

	// Тестуємо дію "update"
	err := author.Validate("update")
	assert.Nil(t, err, "Expected validation to pass for 'update' action")

	// Тестуємо дію "login"
	err = author.Validate("login")
	assert.Nil(t, err, "Expected validation to pass for 'login' action")

	// Тестуємо сценарій з порожнім полем Nickname
	author.Nickname = ""
	err = author.Validate("update")
	assert.NotNil(t, err, "Expected validation to fail for 'update' action with empty Nickname")
	assert.EqualError(t, err, "required Nickname", "Expected error message 'required Nickname'")

	// Тестуємо сценарій update з порожнім полем Password
	author.Nickname = "testuser"
	author.Password = ""
	err = author.Validate("update")
	assert.NotNil(t, err, "Expected validation to fail for 'update' action with empty Password")
	assert.EqualError(t, err, "required Password", "Expected error message 'required Password'")

	// Тестуємо сценарій update з порожнім полем Password
	author.Nickname = "testuser"
	author.Password = ""
	err = author.Validate("login")
	assert.NotNil(t, err, "Expected validation to fail for 'update' action with empty Password")
	assert.EqualError(t, err, "required Password", "Expected error message 'required Password'")

	// Тестуємо сценарій update з порожнім полем Email
	author.Password = "testpassword"
	author.Email = ""
	err = author.Validate("update")
	assert.NotNil(t, err, "Expected validation to fail for 'update' action with empty Email")
	assert.EqualError(t, err, "required Email", "Expected error message 'required Email'")

	// Тестуємо сценарій login з порожнім полем Email
	author.Password = "testpassword"
	author.Email = ""
	err = author.Validate("login")
	assert.NotNil(t, err, "Expected validation to fail for 'update' action with empty Email")
	assert.EqualError(t, err, "required Email", "Expected error message 'required Email'")

	// Тестуємо сценарій update з некоректним Email
	author.Email = "invalidemail"
	err = author.Validate("update")
	assert.NotNil(t, err, "Expected validation to fail for 'update' action with invalid Email")
	assert.EqualError(t, err, "invalid Email", "Expected error message 'invalid Email'")

	// Тестуємо сценарій login з некоректним Email
	author.Email = "invalidemail"
	err = author.Validate("login")
	assert.NotNil(t, err, "Expected validation to fail for 'update' action with invalid Email")
	assert.EqualError(t, err, "invalid Email", "Expected error message 'invalid Email'")

	// Тестуємо дію за замовчуванням (інші дії) з порожнім полем Nickname
	author.Email = "test@example.com"
	author.Nickname = ""
	author.Password = "testpassword"
	err = author.Validate("unknown_action")
	assert.NotNil(t, err, "Expected validation to fail for unknown action with empty Nickname")
	assert.EqualError(t, err, "required Nickname", "Expected error message 'required Nickname'")

	// Тестуємо дію за замовчуванням (інші дії) з порожнім полем Password
	author.Nickname = "testuser"
	author.Password = ""
	err = author.Validate("unknown_action")
	assert.NotNil(t, err, "Expected validation to fail for unknown action with empty Password")
	assert.EqualError(t, err, "required Password", "Expected error message 'required Password'")

	// Тестуємо дію за замовчуванням (інші дії) з порожнім полем Email
	author.Password = "testpassword"
	author.Email = ""
	err = author.Validate("unknown_action")
	assert.NotNil(t, err, "Expected validation to fail for unknown action with empty Email")
	assert.EqualError(t, err, "required Email", "Expected error message 'required Email'")

	// Тестуємо дію за замовчуванням (інші дії) з некоректним Email
	author.Email = "invalidemail"
	err = author.Validate("unknown_action")
	assert.NotNil(t, err, "Expected validation to fail for unknown action with invalid Email")
	assert.EqualError(t, err, "invalid Email", "Expected error message 'invalid Email'")
}

// func TestMain(m *testing.M) {
// 	var err error
// 	err = godotenv.Load(os.ExpandEnv("../../.env"))
// 	if err != nil {
// 		log.Fatalf("Error getting env %v\n", err)
// 	}
// 	Database()

// 	os.Exit(m.Run())

// }

func TestSaveAuthors(t *testing.T) {

}
