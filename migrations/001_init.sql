-- 001_init.sql
-- Initial database schema for coding-winner bot

-- Users table
CREATE TABLE IF NOT EXISTS users (
    discord_id VARCHAR(20) PRIMARY KEY,
    atcoder_username VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Contest notifications configuration per server
CREATE TABLE IF NOT EXISTS contest_notifications (
    id SERIAL PRIMARY KEY,
    server_id VARCHAR(20) NOT NULL,
    channel_id VARCHAR(20) NOT NULL,
    reminder_dm BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(server_id)
);

-- Submissions from AtCoder
CREATE TABLE IF NOT EXISTS submissions (
    id BIGINT PRIMARY KEY,
    user_id VARCHAR(20) REFERENCES users(discord_id) ON DELETE CASCADE,
    problem_id VARCHAR(50) NOT NULL,
    contest_id VARCHAR(50),
    result VARCHAR(20),
    point FLOAT,
    language VARCHAR(50),
    submitted_at TIMESTAMP NOT NULL,
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Problems from AtCoder
CREATE TABLE IF NOT EXISTS problems (
    problem_id VARCHAR(50) PRIMARY KEY,
    contest_id VARCHAR(50),
    title VARCHAR(200),
    difficulty INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Daily problem configuration per server
CREATE TABLE IF NOT EXISTS daily_problem_config (
    server_id VARCHAR(20) PRIMARY KEY,
    channel_id VARCHAR(20) NOT NULL,
    difficulty_min INT DEFAULT 400,
    difficulty_max INT DEFAULT 800,
    post_time TIME DEFAULT '09:00:00',
    enabled BOOLEAN DEFAULT true
);

-- Virtual contests
CREATE TABLE IF NOT EXISTS virtual_contests (
    id SERIAL PRIMARY KEY,
    server_id VARCHAR(20) NOT NULL,
    channel_id VARCHAR(20) NOT NULL,
    created_by VARCHAR(20) REFERENCES users(discord_id) ON DELETE SET NULL,
    title VARCHAR(200) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    duration_minutes INT NOT NULL,
    problem_ids TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Virtual contest submissions
CREATE TABLE IF NOT EXISTS virtual_contest_submissions (
    id SERIAL PRIMARY KEY,
    contest_id INT REFERENCES virtual_contests(id) ON DELETE CASCADE,
    user_id VARCHAR(20) REFERENCES users(discord_id) ON DELETE CASCADE,
    problem_id VARCHAR(50),
    submitted_at TIMESTAMP NOT NULL,
    result VARCHAR(20),
    point FLOAT,
    UNIQUE(contest_id, user_id, problem_id)
);

-- Contest notified messages (for reaction tracking)
CREATE TABLE IF NOT EXISTS contest_notified_messages (
    id SERIAL PRIMARY KEY,
    server_id VARCHAR(20) NOT NULL,
    channel_id VARCHAR(20) NOT NULL,
    message_id VARCHAR(20) NOT NULL UNIQUE,
    contest_id VARCHAR(50) NOT NULL,
    contest_start_time TIMESTAMP NOT NULL,
    notified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Weekly report configuration
CREATE TABLE IF NOT EXISTS weekly_report_config (
    server_id VARCHAR(20) PRIMARY KEY,
    channel_id VARCHAR(20) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    post_day INT DEFAULT 1, -- 1=Monday
    post_time TIME DEFAULT '09:00:00'
);
