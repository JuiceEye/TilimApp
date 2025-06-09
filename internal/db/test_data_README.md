# Achievement System Test Data

This directory contains SQL scripts to set up test data for testing the achievement system functionality.

## Test Data Overview

The `test_data.sql` script creates:

1. A test user with ID 9999
2. User progress with:
   - 2 days streak (needs 1 more day for LESSON_STREAK_3 achievement)
   - 9 total lessons completed (needs 1 more for LESSONS_TOTAL_10 achievement)
3. Test lessons with IDs 9001-9010
4. Lesson completions for the test user:
   - 2 lessons completed 2 days ago
   - 3 lessons completed yesterday
   - 4 lessons completed today (needs 1 more for LESSONS_SINGLE_DAY_5 achievement)

## How to Use

1. Run the SQL script to set up the test data:
   ```bash
   psql -U your_username -d your_database -f internal/db/test_data.sql
   ```

2. Log in as the test user (ID: 9999, username: test_user) or use the API to perform actions as this user.

3. Complete one more lesson to trigger the achievements:
   - This will trigger all three achievements at once:
     - LESSON_STREAK_3 (3-day streak)
     - LESSONS_SINGLE_DAY_5 (5 lessons in one day)
     - LESSONS_TOTAL_10 (10 total lessons)

4. Verify the achievements were granted by checking the database:
   ```sql
   SELECT a.code, a.name, a.description, ua.earned_at 
   FROM app.user_achievements ua
   JOIN app.achievements a ON ua.achievement_id = a.id
   WHERE ua.user_id = 9999;
   ```

5. Verify the XP was awarded by checking the user's XP:
   ```sql
   SELECT xp_points FROM app.user_progress WHERE user_id = 9999;
   ```
   The user should have received:
   - 50 XP for LESSON_STREAK_3
   - 75 XP for LESSONS_SINGLE_DAY_5
   - 25 XP for LESSONS_TOTAL_10
   - Total: 150 XP (plus the original 500 XP = 650 XP)

## Testing Checklist

1. **Streak Achievement (LESSON_STREAK_3)**
   - Complete any lesson today to earn the achievement
   - Verify the user has earned the LESSON_STREAK_3 achievement

2. **Daily Lessons Achievement (LESSONS_SINGLE_DAY_5)**
   - Complete one more lesson today to earn the achievement
   - Verify the user has earned the LESSONS_SINGLE_DAY_5 achievement

3. **Total Lessons Achievement (LESSONS_TOTAL_10)**
   - Complete one more lesson (any lesson) to earn the achievement
   - Verify the user has earned the LESSONS_TOTAL_10 achievement

4. **Verify XP is Awarded for Each Achievement**
   - LESSON_STREAK_3: 50 XP
   - LESSONS_SINGLE_DAY_5: 75 XP
   - LESSONS_TOTAL_10: 25 XP

5. **Verify Achievements are Not Awarded Twice**
   - Complete more lessons after earning the achievements
   - Verify the user still has only one of each achievement

6. **Verify Achievements are Processed Correctly in Transactions**
   - Simulate a transaction failure after achievement processing
   - Verify the achievements are not granted if the transaction fails

## Reset Test Data

To reset the test data and run the tests again:

```sql
-- Delete user achievements
DELETE FROM app.user_achievements WHERE user_id = 9999;

-- Reset user progress
UPDATE app.user_progress 
SET streak_days = 2, 
    xp_points = 500, 
    lessons_done = 9, 
    last_lesson_completed_at = NOW() - INTERVAL '2 days'
WHERE user_id = 9999;

-- Delete lesson completions
DELETE FROM app.lesson_completions WHERE user_id = 9999;

-- Re-insert lesson completions
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
```