package service

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "go-user-api/internal/db"
    "go-user-api/internal/logger"
    "go-user-api/internal/models"

    "go.uber.org/zap"
)

// ErrUserNotFound is returned when a requested user does not exist.
var ErrUserNotFound = errors.New("user not found")

// UserService handles user business logic using SQLC
type UserService struct {
    queries *db.Queries
}

// NewUserService creates a new user service
func NewUserService(sqlDB *sql.DB) *UserService {
    return &UserService{
        queries: db.New(sqlDB),
    }
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name string, dobString string) (*models.User, error) {
    dob, err := time.Parse("2006-01-02", dobString)
    if err != nil {
        logger.Error("Invalid date format", zap.String("dob", dobString))
        return nil, err
    }

    row, err := s.queries.CreateUser(context.Background(), db.CreateUserParams{
        Name: name,
        Dob:  dob,
    })

    if err != nil {
        logger.Error("Failed to create user", zap.Error(err))
        return nil, err
    }

    user := &models.User{
        ID:   row.ID,
        Name: row.Name,
        DOB:  row.Dob,
    }

    logger.Info("User created successfully", zap.Int32("id", row.ID), zap.String("name", row.Name))
    return user, nil
}

// GetUser retrieves a user by ID and calculates age
func (s *UserService) GetUser(id int32) (*models.User, error) {
    row, err := s.queries.GetUser(context.Background(), id)

    if err != nil {
        if err == sql.ErrNoRows {
            logger.Warn("User not found", zap.Int32("id", id))
            return nil, ErrUserNotFound
        }
        logger.Error("Failed to get user", zap.Int32("id", id), zap.Error(err))
        return nil, err
    }

    user := &models.User{
        ID:   row.ID,
        Name: row.Name,
        DOB:  row.Dob,
        Age:  s.CalculateAge(row.Dob),
    }

    return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id int32, name string, dobString string) (*models.User, error) {
    dob, err := time.Parse("2006-01-02", dobString)
    if err != nil {
        logger.Error("Invalid date format", zap.String("dob", dobString))
        return nil, err
    }

    row, err := s.queries.UpdateUser(context.Background(), db.UpdateUserParams{
        ID:   id,
        Name: name,
        Dob:  dob,
    })

    if err != nil {
        if err == sql.ErrNoRows {
            logger.Warn("User not found for update", zap.Int32("id", id))
            return nil, ErrUserNotFound
        }
        logger.Error("Failed to update user", zap.Int32("id", id), zap.Error(err))
        return nil, err
    }

    user := &models.User{
        ID:   row.ID,
        Name: row.Name,
        DOB:  row.Dob,
    }

    logger.Info("User updated successfully", zap.Int32("id", id))
    return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id int32) error {
    err := s.queries.DeleteUser(context.Background(), id)
    if err != nil {
        logger.Error("Failed to delete user", zap.Int32("id", id), zap.Error(err))
        return err
    }

    logger.Info("User deleted successfully", zap.Int32("id", id))
    return nil
}

// GetAllUsers retrieves all users with calculated ages
func (s *UserService) GetAllUsers() ([]models.User, error) {
    rows, err := s.queries.GetAllUsers(context.Background())
    if err != nil {
        logger.Error("Failed to get all users", zap.Error(err))
        return nil, err
    }

    var users []models.User
    for _, row := range rows {
        users = append(users, models.User{
            ID:   row.ID,
            Name: row.Name,
            DOB:  row.Dob,
            Age:  s.CalculateAge(row.Dob),
        })
    }

    logger.Info("Retrieved all users", zap.Int("count", len(users)))
    return users, nil
}

// CalculateAge calculates age from date of birth
func (s *UserService) CalculateAge(dob time.Time) int {
    today := time.Now()
    age := today.Year() - dob.Year()

    if int(today.Month()) < int(dob.Month()) ||
        (int(today.Month()) == int(dob.Month()) && today.Day() < dob.Day()) {
        age--
    }

    return age
}
