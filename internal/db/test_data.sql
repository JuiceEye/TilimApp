-- Test data for achievement system testing

-- Create a test user
-- INSERT INTO auth.users (id, username, password, email, image, registration_date)
-- VALUES (9999, 'test_user', '$2a$10$1qAz2wSx3eDc4rFv5tGb5t', 'test@example.com', 'default.png', NOW())
-- ON CONFLICT (id) DO NOTHING;

-- Create user progress for the test user
INSERT INTO app.user_progress (user_id, streak_days, xp_points, words_learned, lessons_done, last_lesson_completed_at, updated_at, last_streak_reset_date)
VALUES (9999, 2, 500, 100, 9, NOW() - INTERVAL '2 days', NOW(), NULL)
ON CONFLICT (user_id) DO UPDATE SET 
  streak_days = 2,
  xp_points = 500,
  words_learned = 100,
  lessons_done = 9,
  last_lesson_completed_at = NOW() - INTERVAL '2 days',
  updated_at = NOW();

-- Create some test lessons if they don't exist
INSERT INTO app.lessons (id, title, xp, section_id)
VALUES 
  (9001, 'Test Lesson 1', 50, 1),
  (9002, 'Test Lesson 2', 50, 1),
  (9003, 'Test Lesson 3', 50, 1),
  (9004, 'Test Lesson 4', 50, 1),
  (9005, 'Test Lesson 5', 50, 1),
  (9006, 'Test Lesson 6', 50, 1),
  (9007, 'Test Lesson 7', 50, 1),
  (9008, 'Test Lesson 8', 50, 1),
  (9009, 'Test Lesson 9', 50, 1),
  (9010, 'Test Lesson 10', 50, 1)
ON CONFLICT (id) DO NOTHING;

-- Create lesson completions for the test user
-- This sets up a user with:
-- 1. 2 days streak (needs 1 more day for LESSON_STREAK_3)
-- 2. 4 lessons completed today (needs 1 more for LESSONS_SINGLE_DAY_5)
-- 3. 9 total lessons completed (needs 1 more for LESSONS_TOTAL_10)

-- Clear existing completions for the test user
DELETE FROM app.lesson_completions WHERE user_id = 9999;

-- Day before yesterday: 2 lessons
INSERT INTO app.lesson_completions (user_id, lesson_id, date_completed)
VALUES 
  (9999, 9001, NOW() - INTERVAL '2 days'),
  (9999, 9002, NOW() - INTERVAL '2 days');

-- Yesterday: 3 lessons
INSERT INTO app.lesson_completions (user_id, lesson_id, date_completed)
VALUES 
  (9999, 9003, NOW() - INTERVAL '1 day'),
  (9999, 9004, NOW() - INTERVAL '1 day'),
  (9999, 9005, NOW() - INTERVAL '1 day');

-- Today: 4 lessons (one short of LESSONS_SINGLE_DAY_5)
INSERT INTO app.lesson_completions (user_id, lesson_id, date_completed)
VALUES 
  (9999, 9006, NOW()),
  (9999, 9007, NOW()),
  (9999, 9008, NOW()),
  (9999, 9009, NOW());

-- Make sure the user doesn't have any achievements yet
DELETE FROM app.user_achievements WHERE user_id = 9999;

-- Testing checklist:
-- 1. Streak achievement (LESSON_STREAK_3):
--    - Complete any lesson today to earn the achievement
--    - Verify the user has earned the LESSON_STREAK_3 achievement
--
-- 2. Daily lessons achievement (LESSONS_SINGLE_DAY_5):
--    - Complete one more lesson today to earn the achievement
--    - Verify the user has earned the LESSONS_SINGLE_DAY_5 achievement
--
-- 3. Total lessons achievement (LESSONS_TOTAL_10):
--    - Complete one more lesson (any lesson) to earn the achievement
--    - Verify the user has earned the LESSONS_TOTAL_10 achievement
--
-- 4. Verify XP is awarded for each achievement:
--    - LESSON_STREAK_3: 50 XP
--    - LESSONS_SINGLE_DAY_5: 75 XP
--    - LESSONS_TOTAL_10: 25 XP
--
-- 5. Verify achievements are not awarded twice:
--    - Complete more lessons after earning the achievements
--    - Verify the user still has only one of each achievement
--
-- 6. Verify achievements are processed correctly in transactions:
--    - Simulate a transaction failure after achievement processing
--    - Verify the achievements are not granted if the transaction fails