-- Create achievements table
CREATE TABLE IF NOT EXISTS app.achievements (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    xp_reward INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user_achievements table
CREATE TABLE IF NOT EXISTS app.user_achievements (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES auth.users(id),
    achievement_id INTEGER NOT NULL REFERENCES app.achievements(id),
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (user_id, achievement_id)
);

-- Create index for faster queries
CREATE INDEX IF NOT EXISTS idx_user_achievements_user_id ON app.user_achievements(user_id);

-- Insert some sample achievements
INSERT INTO app.achievements (code, name, description, xp_reward)
VALUES 
    ('LESSON_STREAK_3', 'Learning Streak', 'Complete lessons on 3 consecutive days', 50),
    ('LESSON_STREAK_7', 'Weekly Warrior', 'Complete lessons on 7 consecutive days', 100),
    ('LESSON_STREAK_30', 'Monthly Master', 'Complete lessons on 30 consecutive days', 500),
    ('LESSONS_SINGLE_DAY_5', 'Daily Dedication', 'Complete 5 lessons in a single day', 75),
    ('LESSONS_TOTAL_10', 'Getting Started', 'Complete 10 lessons total', 25),
    ('LESSONS_TOTAL_50', 'Learning Machine', 'Complete 50 lessons total', 150),
    ('LESSONS_TOTAL_100', 'Knowledge Seeker', 'Complete 100 lessons total', 300);