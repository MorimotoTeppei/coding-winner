-- 002_indexes.sql
-- Add indexes for performance optimization

-- Submissions indexes
CREATE INDEX IF NOT EXISTS idx_submissions_user_time
ON submissions(user_id, submitted_at DESC);

CREATE INDEX IF NOT EXISTS idx_submissions_problem
ON submissions(problem_id);

CREATE INDEX IF NOT EXISTS idx_submissions_result
ON submissions(result);

CREATE INDEX IF NOT EXISTS idx_submissions_synced
ON submissions(synced_at);

-- Problems indexes
CREATE INDEX IF NOT EXISTS idx_problems_difficulty
ON problems(difficulty);

CREATE INDEX IF NOT EXISTS idx_problems_contest
ON problems(contest_id);

-- Virtual contests indexes
CREATE INDEX IF NOT EXISTS idx_virtual_contests_server
ON virtual_contests(server_id);

CREATE INDEX IF NOT EXISTS idx_virtual_contests_start_time
ON virtual_contests(start_time);

-- Virtual contest submissions indexes
CREATE INDEX IF NOT EXISTS idx_virt_sub_contest
ON virtual_contest_submissions(contest_id);

CREATE INDEX IF NOT EXISTS idx_virt_sub_user
ON virtual_contest_submissions(user_id);

-- Contest notified messages indexes
CREATE INDEX IF NOT EXISTS idx_contest_notified_server
ON contest_notified_messages(server_id);

CREATE INDEX IF NOT EXISTS idx_contest_notified_start_time
ON contest_notified_messages(contest_start_time);
