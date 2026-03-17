BEGIN;

-- Revert achievements to English
UPDATE achievements SET name = 'First Workout', description = 'Complete your first workout' WHERE slug = 'first_workout';
UPDATE achievements SET name = '10 Workouts', description = 'Complete 10 workouts' WHERE slug = '10_workouts';
UPDATE achievements SET name = '50 Workouts', description = 'Complete 50 workouts' WHERE slug = '50_workouts';
UPDATE achievements SET name = 'Week Streak', description = 'Train 7 days in a row' WHERE slug = 'week_streak';
UPDATE achievements SET name = 'Month Streak', description = 'Train 30 days in a row' WHERE slug = 'month_streak';
UPDATE achievements SET name = 'First Rehab', description = 'Complete a rehabilitation course' WHERE slug = 'first_rehab';
UPDATE achievements SET name = 'Pain Reducer', description = 'Reduce pain level by 3+ points during a course' WHERE slug = 'pain_reducer';
UPDATE achievements SET name = 'Food Tracker', description = 'Log food for 7 days straight' WHERE slug = 'food_tracker';
UPDATE achievements SET name = 'Weight Tracker', description = 'Log weight 10 times' WHERE slug = 'weight_tracker';
UPDATE achievements SET name = 'Early Bird', description = 'Complete 5 morning workouts' WHERE slug = 'early_bird';

-- Remove added exercises
DELETE FROM exercises WHERE name IN (
  'Растяжка квадрицепсов стоя',
  'Растяжка задней поверхности бедра',
  'Кошка-корова',
  'Мостик (ягодичный)',
  'Планка на локтях',
  'Выпады на месте',
  'Тяга верхнего блока',
  'Жим ногами в тренажёре',
  'Разведение гантелей лёжа',
  'Гиперэкстензия',
  'Вращения плечами',
  'Подъём на носки стоя'
);

COMMIT;
