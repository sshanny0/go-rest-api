-- Migration: Create users and todos tables
-- Version: 001
-- Description: Initial schema for Todo REST API

-- ============================================
-- Create users table
-- ============================================

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Add comments for documentation
COMMENT ON TABLE users IS 'User accounts for authentication';
COMMENT ON COLUMN users.id IS 'Primary key, auto-increment';
COMMENT ON COLUMN users.username IS 'Unique username for login';
COMMENT ON COLUMN users.email IS 'Unique email address';
COMMENT ON COLUMN users.password IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.full_name IS 'User''s full name';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp (NULL = not deleted)';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- ============================================
-- Create todos table
-- ============================================

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',
    due_date DATE,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Foreign key constraint
    CONSTRAINT fk_todos_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- Add comments for documentation
COMMENT ON TABLE todos IS 'Todo items belonging to users';
COMMENT ON COLUMN todos.id IS 'Primary key, auto-increment';
COMMENT ON COLUMN todos.title IS 'Todo title (max 200 chars)';
COMMENT ON COLUMN todos.description IS 'Optional detailed description';
COMMENT ON COLUMN todos.status IS 'Status: pending, in_progress, completed';
COMMENT ON COLUMN todos.priority IS 'Priority: low, medium, high';
COMMENT ON COLUMN todos.due_date IS 'Optional due date';
COMMENT ON COLUMN todos.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN todos.deleted_at IS 'Soft delete timestamp (NULL = not deleted)';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);
CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
CREATE INDEX IF NOT EXISTS idx_todos_deleted_at ON todos(deleted_at);
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);

-- Composite index for common query (user's todos by status)
CREATE INDEX IF NOT EXISTS idx_todos_user_status ON todos(user_id, status);

-- ============================================
-- Add constraints for data validation
-- ============================================

-- Check constraint for status values
ALTER TABLE todos
    ADD CONSTRAINT check_todos_status
    CHECK (status IN ('pending', 'in_progress', 'completed'));

-- Check constraint for priority values
ALTER TABLE todos
    ADD CONSTRAINT check_todos_priority
    CHECK (priority IN ('low', 'medium', 'high'));

-- ============================================
-- Create updated_at trigger function
-- ============================================

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for users table
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for todos table
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;
CREATE TRIGGER update_todos_updated_at
    BEFORE UPDATE ON todos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- Migration complete
-- ============================================

-- Display success message
DO $$
BEGIN
    RAISE NOTICE 'Migration 001 completed successfully';
    RAISE NOTICE 'Created tables: users, todos';
    RAISE NOTICE 'Created indexes for performance';
    RAISE NOTICE 'Created triggers for updated_at';
END $$;
