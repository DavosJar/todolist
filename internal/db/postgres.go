package db

import (
	"context"
	"time"
	"todo_list/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New(databaseURL string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate crea/actualiza tablas automáticamente
func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(&models.User{}, &models.Tenant{}, &models.Task{})
}

// User operations
func (d *Database) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	if err := d.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := d.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Database) GetUserPassword(ctx context.Context, userID string) (string, error) {
	var user models.User
	if err := d.DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.PasswordHash, nil
}

// Tenant operations
func (d *Database) CreateTenant(ctx context.Context, userID string, name string) (*models.Tenant, error) {
	tenant := &models.Tenant{
		UserID: userID,
		Name:   name,
	}
	if err := d.DB.WithContext(ctx).Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

func (d *Database) GetTenantByUserID(ctx context.Context, userID string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := d.DB.WithContext(ctx).Where("user_id = ?", userID).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// Task operations
func (d *Database) CreateTask(ctx context.Context, tenantID, title string) (*models.Task, error) {
	task := &models.Task{
		TenantID: tenantID,
		Title:    title,
	}
	if err := d.DB.WithContext(ctx).Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (d *Database) GetTasks(ctx context.Context, tenantID string) ([]models.Task, error) {
	var tasks []models.Task
	if err := d.DB.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (d *Database) UpdateTask(ctx context.Context, taskID, tenantID string, completed bool) (*models.Task, error) {
	var task models.Task
	if err := d.DB.WithContext(ctx).Where("id = ? AND tenant_id = ?", taskID, tenantID).First(&task).Error; err != nil {
		return nil, err
	}
	task.Completed = completed
	if err := d.DB.WithContext(ctx).Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (d *Database) DeleteTask(ctx context.Context, taskID, tenantID string) error {
	return d.DB.WithContext(ctx).Where("id = ? AND tenant_id = ?", taskID, tenantID).Delete(&models.Task{}).Error
}
