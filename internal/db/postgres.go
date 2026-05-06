package db

import (
	"context"
	"time"
	"todo_list/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(databaseURL string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Database{pool: pool}, nil
}

// User operations
func (db *Database) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		CreatedAt: time.Now(),
	}

	_, err := db.pool.Exec(ctx,
		"INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)",
		user.ID, email, passwordHash)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := db.pool.QueryRow(ctx,
		"SELECT id, email, created_at FROM users WHERE email = $1",
		email).
		Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (db *Database) GetUserPassword(ctx context.Context, userID uuid.UUID) (string, error) {
	var passwordHash string
	err := db.pool.QueryRow(ctx,
		"SELECT password_hash FROM users WHERE id = $1",
		userID).
		Scan(&passwordHash)

	if err != nil {
		return "", err
	}

	return passwordHash, nil
}

// Tenant operations
func (db *Database) CreateTenant(ctx context.Context, userID uuid.UUID, name string) (*models.Tenant, error) {
	tenant := &models.Tenant{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
	}

	_, err := db.pool.Exec(ctx,
		"INSERT INTO tenants (id, user_id, name) VALUES ($1, $2, $3)",
		tenant.ID, userID, name)

	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (db *Database) GetTenantByUserID(ctx context.Context, userID uuid.UUID) (*models.Tenant, error) {
	tenant := &models.Tenant{}
	err := db.pool.QueryRow(ctx,
		"SELECT id, user_id, name, created_at FROM tenants WHERE user_id = $1",
		userID).
		Scan(&tenant.ID, &tenant.UserID, &tenant.Name, &tenant.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return tenant, nil
}

// Task operations
func (db *Database) GetTasks(ctx context.Context, tenantID uuid.UUID) ([]models.Task, error) {
	rows, err := db.pool.Query(ctx,
		"SELECT id, tenant_id, title, completed, created_at FROM tasks WHERE tenant_id = $1 ORDER BY created_at DESC",
		tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.TenantID, &task.Title, &task.Completed, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (db *Database) CreateTask(ctx context.Context, tenantID uuid.UUID, title string) (*models.Task, error) {
	task := &models.Task{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	_, err := db.pool.Exec(ctx,
		"INSERT INTO tasks (id, tenant_id, title, completed) VALUES ($1, $2, $3, $4)",
		task.ID, tenantID, title, false)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (db *Database) UpdateTask(ctx context.Context, taskID uuid.UUID, tenantID uuid.UUID, completed bool) (*models.Task, error) {
	task := &models.Task{}
	err := db.pool.QueryRow(ctx,
		"UPDATE tasks SET completed = $1 WHERE id = $2 AND tenant_id = $3 RETURNING id, tenant_id, title, completed, created_at",
		completed, taskID, tenantID).
		Scan(&task.ID, &task.TenantID, &task.Title, &task.Completed, &task.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return task, nil
}

func (db *Database) DeleteTask(ctx context.Context, taskID uuid.UUID, tenantID uuid.UUID) error {
	commandTag, err := db.pool.Exec(ctx,
		"DELETE FROM tasks WHERE id = $1 AND tenant_id = $2",
		taskID, tenantID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) Close() {
	db.pool.Close()
}
