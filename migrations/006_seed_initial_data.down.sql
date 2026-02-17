BEGIN;

DELETE FROM question_options WHERE question_id IN (
    SELECT id FROM questions WHERE questionnaire_id IN (
        SELECT id FROM questionnaires WHERE slug = 'health_test'
    )
);
DELETE FROM questions WHERE questionnaire_id IN (
    SELECT id FROM questionnaires WHERE slug = 'health_test'
);
DELETE FROM questionnaires WHERE slug = 'health_test';
DELETE FROM module_categories WHERE module_id IN (
    SELECT id FROM modules WHERE slug IN ('lfk', 'trainings', 'food')
);
DELETE FROM modules WHERE slug IN ('lfk', 'trainings', 'food');
DELETE FROM subscription_plans WHERE slug IN ('monthly', 'quarterly', 'yearly');

COMMIT;
