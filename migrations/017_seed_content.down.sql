-- Remove all seeded content (reverse of 017_seed_content.up.sql)

DELETE FROM workout_exercises;
DELETE FROM workouts;
DELETE FROM program_enrollments;
DELETE FROM programs;
DELETE FROM exercises;
DELETE FROM rehab_sessions;
DELETE FROM rehab_courses;
DELETE FROM meals;
DELETE FROM meal_plans;
