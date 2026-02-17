BEGIN;

-- Subscription plans
INSERT INTO subscription_plans (slug, name, description, price_kzt, duration_days) VALUES
    ('monthly', 'Месячная подписка', 'Доступ ко всем модулям на 30 дней', 5000, 30),
    ('quarterly', 'Квартальная подписка', 'Доступ ко всем модулям на 90 дней', 12000, 90),
    ('yearly', 'Годовая подписка', 'Доступ ко всем модулям на 365 дней', 40000, 365);

-- Modules
INSERT INTO modules (slug, name, description, icon, requires_subscription, sort_order) VALUES
    ('lfk', 'ЛФК', 'Лечебная физическая культура — упражнения при проблемах с позвоночником и суставами', '🏥', TRUE, 1),
    ('trainings', 'Тренировки', 'Тренировки по группам мышц с видео-уроками', '💪', TRUE, 2),
    ('food', 'Питание', 'Правильное питание и рецепты', '🥗', TRUE, 3);

-- ЛФК categories
INSERT INTO module_categories (module_id, slug, name, description, icon, sort_order) VALUES
    ((SELECT id FROM modules WHERE slug = 'lfk'), 'hernia', 'Грыжа', 'Упражнения при грыже позвоночника', '🔴', 1),
    ((SELECT id FROM modules WHERE slug = 'lfk'), 'protrusion', 'Протрузии', 'Упражнения при протрузиях', '🟠', 2),
    ((SELECT id FROM modules WHERE slug = 'lfk'), 'scoliosis', 'Сколиоз', 'Упражнения при сколиозе', '🟡', 3),
    ((SELECT id FROM modules WHERE slug = 'lfk'), 'kyphosis', 'Кифоз', 'Упражнения при кифозе', '🟢', 4),
    ((SELECT id FROM modules WHERE slug = 'lfk'), 'lordosis', 'Лордоз', 'Упражнения при лордозе', '🔵', 5);

-- Training categories
INSERT INTO module_categories (module_id, slug, name, description, icon, sort_order) VALUES
    ((SELECT id FROM modules WHERE slug = 'trainings'), 'chest', 'Грудь', 'Тренировки грудных мышц', '💪', 1),
    ((SELECT id FROM modules WHERE slug = 'trainings'), 'back', 'Спина', 'Тренировки мышц спины', '🏋', 2),
    ((SELECT id FROM modules WHERE slug = 'trainings'), 'legs', 'Ноги', 'Тренировки ног', '🦵', 3),
    ((SELECT id FROM modules WHERE slug = 'trainings'), 'shoulders', 'Плечи', 'Тренировки дельтовидных мышц', '🤸', 4),
    ((SELECT id FROM modules WHERE slug = 'trainings'), 'arms', 'Руки', 'Тренировки бицепсов и трицепсов', '💪', 5),
    ((SELECT id FROM modules WHERE slug = 'trainings'), 'core', 'Пресс', 'Тренировки мышц кора', '🔥', 6);

-- Food categories
INSERT INTO module_categories (module_id, slug, name, description, icon, sort_order) VALUES
    ((SELECT id FROM modules WHERE slug = 'food'), 'breakfast', 'Завтраки', 'Рецепты полезных завтраков', '🌅', 1),
    ((SELECT id FROM modules WHERE slug = 'food'), 'lunch', 'Обеды', 'Рецепты обедов', '☀', 2),
    ((SELECT id FROM modules WHERE slug = 'food'), 'dinner', 'Ужины', 'Рецепты ужинов', '🌙', 3),
    ((SELECT id FROM modules WHERE slug = 'food'), 'snacks', 'Перекусы', 'Здоровые перекусы', '🍎', 4);

-- Health test questionnaire
INSERT INTO questionnaires (slug, title, description, is_active, sort_order) VALUES
    ('health_test', 'Тест здоровья', 'Пройдите тест, чтобы мы могли подобрать для вас оптимальную программу', TRUE, 1);

INSERT INTO questions (questionnaire_id, text, question_type, sort_order, is_required, metadata) VALUES
    ((SELECT id FROM questionnaires WHERE slug = 'health_test'),
     'Как часто вы занимаетесь физической активностью?', 'single_choice', 1, TRUE, '{}'),
    ((SELECT id FROM questionnaires WHERE slug = 'health_test'),
     'Есть ли у вас проблемы со здоровьем?', 'multiple_choice', 2, TRUE, '{}'),
    ((SELECT id FROM questionnaires WHERE slug = 'health_test'),
     'Оцените ваш текущий уровень энергии', 'scale', 3, TRUE, '{"min": 1, "max": 10, "min_label": "Очень низкий", "max_label": "Очень высокий"}'),
    ((SELECT id FROM questionnaires WHERE slug = 'health_test'),
     'Опишите ваши ожидания от программы', 'text', 4, FALSE, '{}');

-- Options for question 1
INSERT INTO question_options (question_id, text, value, sort_order)
SELECT q.id, o.text, o.value, o.sort_order
FROM questions q,
(VALUES
    ('Не занимаюсь', 'none', 1),
    ('1-2 раза в неделю', '1-2_per_week', 2),
    ('3-4 раза в неделю', '3-4_per_week', 3),
    ('5+ раз в неделю', '5+_per_week', 4)
) AS o(text, value, sort_order)
WHERE q.text = 'Как часто вы занимаетесь физической активностью?';

-- Options for question 2
INSERT INTO question_options (question_id, text, value, sort_order)
SELECT q.id, o.text, o.value, o.sort_order
FROM questions q,
(VALUES
    ('Нет проблем', 'none', 1),
    ('Проблемы с позвоночником', 'spine', 2),
    ('Проблемы с суставами', 'joints', 3),
    ('Сердечно-сосудистые заболевания', 'cardiovascular', 4),
    ('Другое', 'other', 5)
) AS o(text, value, sort_order)
WHERE q.text = 'Есть ли у вас проблемы со здоровьем?';

COMMIT;
