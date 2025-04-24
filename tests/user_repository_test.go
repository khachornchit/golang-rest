package repository_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang-rest/internal/core/domain"
	"golang-rest/internal/core/ports"
	"testing"
)

type MockUserRepository struct {
	mock.Mock
	Users map[string]domain.User
}

func NewMockUserRepo() ports.UserRepositoryInterface {
	return &MockUserRepository{Users: make(map[string]domain.User)}
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	m.Users[user.Email] = *user
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	var userList []domain.User
	for _, u := range m.Users {
		u.Password = ""
		userList = append(userList, u)
	}
	return userList, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	user, ok := m.Users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	user.Password = ""
	return &user, args.Error(1)
}

func (m *MockUserRepository) GetUserLoginByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	user, ok := m.Users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return &user, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id string) (*domain.User, error) {
	args := m.Called(id)
	for _, u := range m.Users {
		if u.ID.Hex() == id {
			u.Password = ""
			return &u, args.Error(1)
		}
	}
	return nil, errors.New("not found")
}

func (m *MockUserRepository) UpdateUserByID(id string, updates bson.M) (*domain.User, error) {
	args := m.Called(id, updates)
	for k, u := range m.Users {
		if u.ID.Hex() == id {
			if name, ok := updates["name"].(string); ok {
				u.Name = name
			}
			if email, ok := updates["email"].(string); ok {
				u.Email = email
			}
			m.Users[k] = u
			return &u, args.Error(1)
		}
	}
	return nil, errors.New("not found")
}

func (m *MockUserRepository) DeleteUserByID(id string) error {
	args := m.Called(id)
	for k, u := range m.Users {
		if u.ID.Hex() == id {
			delete(m.Users, k)
			return args.Error(0)
		}
	}
	return errors.New("not found")
}

func TestUserRepoMock_CreateAndFetchUser(t *testing.T) {
	mockRepo := &MockUserRepository{Users: make(map[string]domain.User)}
	mockRepo.On("CreateUser", mock.Anything).Return(nil)
	mockRepo.On("GetUserByEmail", "alice@example.com").Return(nil, nil)

	user := &domain.User{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret",
	}
	err := mockRepo.CreateUser(user)
	assert.NoError(t, err)

	res, err := mockRepo.GetUserByEmail("alice@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", res.Name)
	assert.Empty(t, res.Password)
}

func TestUserRepoMock_UpdateAndDeleteUser(t *testing.T) {
	mockRepo := &MockUserRepository{Users: make(map[string]domain.User)}
	mockRepo.On("CreateUser", mock.Anything).Return(nil)
	mockRepo.On("UpdateUserByID", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("DeleteUserByID", mock.Anything).Return(nil)

	user := &domain.User{
		ID:       primitive.NewObjectID(),
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "pwd123",
	}
	mockRepo.CreateUser(user)

	updates := bson.M{"name": "Robert"}
	updated, err := mockRepo.UpdateUserByID(user.ID.Hex(), updates)
	assert.NoError(t, err)
	assert.Equal(t, "Robert", updated.Name)

	err = mockRepo.DeleteUserByID(user.ID.Hex())
	assert.NoError(t, err)
}

func TestUserRepoMock_Login(t *testing.T) {
	// Initialize the mock repository
	mockRepo := &MockUserRepository{
		Users: make(map[string]domain.User),
	}

	// Create a test user
	testUser := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "securepassword",
	}

	// Add the test user to the mock repository
	mockRepo.Users[testUser.Email] = testUser

	// Set up expectations
	mockRepo.On("GetUserLoginByEmail", testUser.Email).Return(&testUser, nil)

	// Call the method under test
	user, err := mockRepo.GetUserLoginByEmail(testUser.Email)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.Email, user.Email)
	assert.Equal(t, testUser.Password, user.Password)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUserRepoMock_GetAllUsers(t *testing.T) {
	// Initialize the mock repository
	mockRepo := &MockUserRepository{
		Users: make(map[string]domain.User),
	}

	// Create test users
	user1 := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password1",
	}
	user2 := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "password2",
	}

	// Add test users to the mock repository
	mockRepo.Users[user1.Email] = user1
	mockRepo.Users[user2.Email] = user2

	// Set up expectations
	mockRepo.On("GetAllUsers").Return([]domain.User{user1, user2}, nil)

	// Call the method under test
	users, err := mockRepo.GetAllUsers()

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	assert.Equal(t, "Bob", users[1].Name)
	assert.Empty(t, users[0].Password)
	assert.Empty(t, users[1].Password)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)
}
