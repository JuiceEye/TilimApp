-- Create daily_tasks table
CREATE TABLE IF NOT EXISTS app.daily_tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    xp INTEGER NOT NULL,
    lesson_id INTEGER NOT NULL REFERENCES app.lessons(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user_daily_tasks table
CREATE TABLE IF NOT EXISTS app.user_daily_tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES auth.users(id),
    daily_task_id INTEGER NOT NULL REFERENCES app.daily_tasks(id),
    lesson_id INTEGER NOT NULL REFERENCES app.lessons(id),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    assigned_date TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE (user_id, daily_task_id, assigned_date)
);

-- Create index for faster queries
CREATE INDEX IF NOT EXISTS idx_user_daily_tasks_user_id ON app.user_daily_tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_user_daily_tasks_assigned_date ON app.user_daily_tasks(assigned_date);
CREATE INDEX IF NOT EXISTS idx_user_daily_tasks_lesson_id ON app.user_daily_tasks(lesson_id);

-- Insert some sample daily tasks
INSERT INTO app.daily_tasks (title, description, xp, lesson_id)
VALUES 
    ('Complete a basic lesson', 'Complete any basic level lesson to earn XP', 50, 1),
    ('Practice vocabulary', 'Complete a vocabulary lesson to improve your skills', 75, 2),
    ('Learn new grammar', 'Complete a grammar lesson to enhance your understanding', 100, 3),
    ('Review previous lessons', 'Go back and review a lesson you''ve already completed', 50, 4),
    ('Complete a challenge', 'Take on a challenging lesson to test your knowledge', 150, 5);