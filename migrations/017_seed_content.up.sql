BEGIN;

-- ============================================================================
-- EXERCISES
-- ============================================================================

INSERT INTO exercises (name, technique, common_mistakes, easier_modification, harder_modification, rest_seconds) VALUES
('Отжимания', 'Примите упор лежа, руки на ширине плеч. Опускайтесь до касания грудью пола, затем отжимайтесь вверх, сохраняя тело прямым. Локти держите под углом 45 градусов к телу.', 'Прогиб в пояснице, разведение локтей в стороны, неполная амплитуда движения, задержка дыхания.', 'Отжимания с колен или от стены. Уменьшите амплитуду движения.', 'Отжимания с ногами на возвышении, с хлопком, на одной руке или с дополнительным весом на спине.', 90),
('Приседания', 'Встаньте прямо, ноги на ширине плеч. Отводя таз назад, опускайтесь до параллели бедер с полом. Колени не выходят за носки. Спина прямая, взгляд вперед.', 'Колени выходят за носки, отрыв пяток от пола, округление спины, сведение коленей внутрь.', 'Приседания до комфортной глубины, придерживаясь за опору. Можно использовать стул.', 'Приседания с прыжком, на одной ноге (пистолетиком), с дополнительным весом (гантели, штанга).', 60),
('Планка', 'Примите упор на предплечьях и носках. Тело образует прямую линию от головы до пяток. Напрягите пресс и ягодицы. Удерживайте позицию.', 'Прогиб или округление спины, поднятие таза вверх, задержка дыхания, опускание головы.', 'Планка на коленях, уменьшение времени удержания, планка на вытянутых руках.', 'Планка с поочередным подъемом рук или ног, боковая планка, динамическая планка.', 45),
('Выпады', 'Из положения стоя сделайте широкий шаг вперед. Опуститесь, сгибая оба колена под 90 градусов. Переднее колено над пяткой. Вернитесь в исходное положение.', 'Переднее колено выходит за носок, узкий шаг, наклон корпуса вперед, потеря равновесия.', 'Выпады назад (они легче), придерживаясь за опору, укороченная амплитуда.', 'Выпады в прыжке, болгарские выпады (задняя нога на возвышении), выпады с гантелями.', 60),
('Берпи', 'Из положения стоя присядьте, упритесь руками в пол. Прыжком примите упор лежа, выполните отжимание. Прыжком подтяните ноги к рукам и выпрыгните вверх с хлопком над головой.', 'Пропуск отжимания, недостаточная амплитуда прыжка, прогиб в пояснице при отжимании.', 'Без прыжка и отжимания, просто переход в планку и обратно. Шаговые берпи.', 'Берпи с прыжком на возвышение, с дополнительным весом, с двойным отжиманием.', 90),
('Скручивания', 'Лежа на спине, ноги согнуты, стопы на полу. Руки за головой. Поднимайте верхнюю часть корпуса, напрягая пресс. Поясница прижата к полу.', 'Отрыв поясницы от пола, тяга головы руками, резкие рывковые движения, задержка дыхания.', 'Скручивания с меньшей амплитудой, руки вытянуты вперед для облегчения.', 'Скручивания на наклонной скамье, с поворотом корпуса, с поднятыми ногами.', 45),
('Подтягивания', 'Возьмитесь за перекладину хватом чуть шире плеч. Подтянитесь вверх до касания подбородком перекладины. Опуститесь в исходное положение контролируя движение.', 'Раскачивание тела, неполная амплитуда, рывковые движения, отсутствие контроля при опускании.', 'Австралийские подтягивания (низкая перекладина), негативные подтягивания, с резинкой для помощи.', 'Подтягивания с дополнительным весом, на одной руке, с уголком (ноги перед собой).', 120),
('Жим гантелей лежа', 'Лежа на скамье, гантели в руках над грудью. Опустите гантели в стороны, сгибая локти. Выжмите вверх, сводя гантели в верхней точке.', 'Слишком большой вес, отрыв таза от скамьи, неконтролируемое опускание гантелей, разная траектория рук.', 'Уменьшение веса, жим на наклонной скамье с меньшим углом.', 'Увеличение веса, жим на наклонной скамье вверх, медленные негативы.', 90),
('Тяга гантели в наклоне', 'Упритесь рукой и коленом в скамью. Второй ногой стойте на полу. Гантель в свободной руке. Подтягивайте гантель к поясу, отводя локоть назад. Спина прямая.', 'Скручивание корпуса, тяга гантели не к поясу а к груди, рывковые движения, округление спины.', 'Уменьшение веса, тяга двумя руками стоя в наклоне для стабильности.', 'Увеличение веса, тяга с паузой в верхней точке, медленный негатив.', 60),
('Жим гантелей стоя', 'Стоя, гантели на уровне плеч. Выжмите гантели вверх над головой. Опустите контролируя движение. Держите пресс напряженным, не прогибайтесь в пояснице.', 'Прогиб в пояснице, выведение гантелей вперед, неполная амплитуда, разная траектория рук.', 'Жим сидя на скамье с опорой для спины, уменьшение веса.', 'Жим стоя с большим весом, жим одной рукой, жим с паузой в верхней точке.', 75),
('Приседания с гантелями', 'Гантели в руках вдоль тела или на плечах. Выполните приседание как описано выше. Гантели добавляют нагрузку.', 'Те же ошибки что и в обычных приседаниях плюс потеря контроля над гантелями.', 'Приседания с одной легкой гантелью, приседания с кубковым хватом (гантель перед грудью).', 'Приседания с тяжелыми гантелями, приседания с паузой внизу, болгарские сплит-приседания с гантелями.', 90),
('Румынская тяга', 'Стоя, ноги на ширине плеч, гантели или штанга в руках. Наклонитесь вперед с прямой спиной, отводя таз назад. Гриф скользит вдоль ног. Почувствуйте растяжение задней поверхности бедра.', 'Округление спины, сгибание коленей как в приседе, отведение грифа от ног.', 'Румынская тяга с легким весом, с гантелями (легче контролировать), частичная амплитуда.', 'Румынская тяга на одной ноге, с большим весом, с паузой в нижней точке.', 75),
('Подъемы на носки', 'Стоя, поднимитесь на носки максимально высоко. Задержитесь на секунду. Медленно опуститесь. Можно держать гантели для утяжеления.', 'Неполная амплитуда, отсутствие паузы в верхней точке, раскачивание корпуса.', 'Подъемы без веса, подъемы сидя (меньше нагрузки).', 'Подъемы с большим весом, на одной ноге, с паузой в верхней точке.', 45),
('Разведение гантелей лежа', 'Лежа на скамье, гантели над грудью, руки почти прямые. Разведите гантели в стороны до растяжения грудных мышц. Сведите обратно.', 'Слишком большой вес, сильное сгибание локтей (превращается в жим), неконтролируемое опускание.', 'Уменьшение веса, меньшая амплитуда разведения.', 'Увеличение веса, разведение на наклонной скамье, медленный негатив.', 60),
('Сгибания рук с гантелями', 'Стоя, гантели в руках вдоль тела. Согните руки в локтях, поднимая гантели к плечам. Опустите контролируя движение. Локти неподвижны.', 'Раскачивание корпуса, отведение локтей назад или вперед, неполная амплитуда, рывковые движения.', 'Сгибания сидя с опорой для спины, молотковый хват, уменьшение веса.', 'Сгибания с паузой в верхней точке, медленный негатив (4-5 секунд), концентрированные сгибания.', 60),
('Французский жим лежа', 'Лежа на скамье, гантели или гриф над грудью, руки вертикально. Согните руки в локтях, опуская вес за голову. Разогните руки. Локти неподвижны.', 'Разведение локтей в стороны, движение в плечевом суставе, слишком большой вес.', 'Уменьшение веса, французский жим с гантелями (легче контролировать), частичная амплитуда.', 'Французский жим с большим весом, с паузой в нижней точке, медленный негатив.', 75),
('Гиперэкстензия', 'На специальном тренажере или лежа на животе на полу. Поднимите корпус вверх до прямой линии с ногами. Руки скрещены на груди или за головой.', 'Переразгибание вверх (выше линии тела), рывковые движения, слишком быстрое выполнение.', 'Гиперэкстензия на полу (лежа на животе), руки вдоль тела, частичная амплитуда.', 'Гиперэкстензия с дополнительным весом (блин на груди), медленный темп, с паузой вверху.', 60),
('Боковая планка', 'Лежа на боку, упор на предплечье, ноги прямые. Поднимите таз, тело образует прямую линию. Свободная рука на поясе или вверх. Удерживайте позицию.', 'Опускание таза, скручивание корпуса, потеря равновесия.', 'Боковая планка с упором на колено, на вытянутой руке, уменьшение времени.', 'Боковая планка с подъемом верхней ноги, с опусканием и подъемом таза, с весом.', 45),
('Горный альпинист', 'Примите упор лежа. Поочередно подтягивайте колени к груди в быстром темпе. Тело остается в положении планки, таз не поднимается.', 'Подъем таза вверх, медленный темп, отсутствие контроля над корпусом, прогиб в пояснице.', 'Медленный темп выполнения, шаговые подтягивания коленей без прыжка.', 'Ускорение темпа, подтягивание колена к противоположному локтю, с ногами на возвышении.', 45),
('Махи гантелями в стороны', 'Стоя, гантели в руках вдоль тела. Поднимите гантели в стороны до уровня плеч. Опустите контролируя движение. Локти слегка согнуты.', 'Поднятие плеч вверх, раскачивание корпуса, слишком высокий подъем гантелей, прямые руки.', 'Махи сидя (убирает инерцию), с легкими гантелями, частичная амплитуда.', 'Махи с тяжелыми гантелями, с паузой в верхней точке, медленный негатив.', 60);

-- ============================================================================
-- PROGRAMS
-- ============================================================================

INSERT INTO programs (slug, name, description, goal, format, level, duration_weeks, is_active, sort_order) VALUES
('home-beginner-weight-loss', 'Похудение дома для начинающих', 'Эффективная программа для снижения веса в домашних условиях. Не требует специального оборудования. Сочетание кардио и силовых упражнений для максимального жиросжигания.', 'weight_loss', 'home', 'beginner', 4, true, 1),
('home-intermediate-maintenance', 'Поддержание формы дома', 'Программа для поддержания достигнутых результатов и общего тонуса мышц. Разнообразные тренировки для всего тела. Средний уровень интенсивности.', 'maintenance', 'home', 'intermediate', 6, true, 2),
('gym-beginner-muscle-gain', 'Набор массы в зале для начинающих', 'Базовая программа для набора мышечной массы в тренажерном зале. Фокус на базовых упражнениях и правильной технике. Постепенное увеличение нагрузки.', 'muscle_gain', 'gym', 'beginner', 8, true, 3),
('gym-intermediate-strength', 'Развитие силы в зале', 'Программа для развития силовых показателей. Работа с большими весами в базовых упражнениях. Включает периодизацию нагрузки.', 'strength', 'gym', 'intermediate', 6, true, 4),
('home-beginner-fitness', 'Общая физическая подготовка дома', 'Универсальная программа для улучшения общей физической формы. Развитие силы, выносливости и гибкости. Тренировки 3-4 раза в неделю по 30-40 минут.', 'endurance', 'home', 'beginner', 4, true, 5),
('gym-advanced-muscle-gain', 'Продвинутый набор массы в зале', 'Интенсивная программа для опытных атлетов. Сплит-тренировки с фокусом на отдельные мышечные группы. Дропсеты, суперсеты, частичные повторения.', 'muscle_gain', 'gym', 'advanced', 8, true, 6);

-- ============================================================================
-- WORKOUTS (using currval to reference program IDs)
-- ============================================================================

INSERT INTO workouts (program_id, slug, name, description, goal, format, level, duration_minutes, equipment, sort_order, week_number, day_number, is_active) VALUES
-- Program 1: Home beginner weight loss
((SELECT id FROM programs WHERE slug='home-beginner-weight-loss'), 'hw-w1d1', 'Кардио-разминка', 'Вводная кардио-тренировка. Умеренная интенсивность, фокус на технике.', 'weight_loss', 'home', 'beginner', 30, '{}', 1, 1, 1, true),
((SELECT id FROM programs WHERE slug='home-beginner-weight-loss'), 'hw-w1d3', 'Силовая всего тела', 'Базовые силовые упражнения с собственным весом для всех мышечных групп.', 'weight_loss', 'home', 'beginner', 35, '{}', 2, 1, 3, true),
((SELECT id FROM programs WHERE slug='home-beginner-weight-loss'), 'hw-w2d1', 'Интервальная тренировка', 'Чередование интенсивных и восстановительных интервалов для жиросжигания.', 'weight_loss', 'home', 'beginner', 35, '{}', 3, 2, 1, true),
((SELECT id FROM programs WHERE slug='home-beginner-weight-loss'), 'hw-w3d2', 'Комплексная жиросжигающая', 'Интенсивная тренировка с комбинированными упражнениями.', 'weight_loss', 'home', 'beginner', 40, '{}', 4, 3, 2, true),

-- Program 2: Home intermediate maintenance
((SELECT id FROM programs WHERE slug='home-intermediate-maintenance'), 'hm-w1d1', 'Верх тела и кор', 'Проработка мышц верхней части тела и пресса.', 'maintenance', 'home', 'intermediate', 40, ARRAY['коврик', 'стул'], 1, 1, 1, true),
((SELECT id FROM programs WHERE slug='home-intermediate-maintenance'), 'hm-w2d3', 'Ноги и ягодицы', 'Целевая тренировка нижней части тела.', 'maintenance', 'home', 'intermediate', 45, ARRAY['коврик'], 2, 2, 3, true),
((SELECT id FROM programs WHERE slug='home-intermediate-maintenance'), 'hm-w3d2', 'Функциональная тренировка', 'Комплексные многосуставные упражнения для координации и баланса.', 'maintenance', 'home', 'intermediate', 45, ARRAY['коврик', 'стул'], 3, 3, 2, true),
((SELECT id FROM programs WHERE slug='home-intermediate-maintenance'), 'hm-w5d1', 'HIIT тренировка', 'Высокоинтенсивная интервальная тренировка для ускорения метаболизма.', 'maintenance', 'home', 'intermediate', 35, ARRAY['коврик'], 4, 5, 1, true),

-- Program 3: Gym beginner muscle gain
((SELECT id FROM programs WHERE slug='gym-beginner-muscle-gain'), 'gb-w1d1', 'Грудь и трицепс', 'Базовая тренировка грудных мышц и трицепсов с гантелями.', 'muscle_gain', 'gym', 'beginner', 50, ARRAY['скамья', 'гантели'], 1, 1, 1, true),
((SELECT id FROM programs WHERE slug='gym-beginner-muscle-gain'), 'gb-w1d3', 'Спина и бицепс', 'Тяговые упражнения для спины и проработка бицепсов.', 'muscle_gain', 'gym', 'beginner', 50, ARRAY['турник', 'гантели'], 2, 1, 3, true),
((SELECT id FROM programs WHERE slug='gym-beginner-muscle-gain'), 'gb-w2d2', 'Ноги', 'Тренировка всех мышц ног с гантелями и штангой.', 'muscle_gain', 'gym', 'beginner', 55, ARRAY['штанга', 'гантели'], 3, 2, 2, true),
((SELECT id FROM programs WHERE slug='gym-beginner-muscle-gain'), 'gb-w3d1', 'Плечи и пресс', 'Тренировка дельтовидных мышц и кора.', 'muscle_gain', 'gym', 'beginner', 45, ARRAY['гантели', 'коврик'], 4, 3, 1, true),

-- Program 4: Gym intermediate strength
((SELECT id FROM programs WHERE slug='gym-intermediate-strength'), 'gs-w1d1', 'Тяжелая тяга', 'Силовая тренировка с акцентом на тяговые движения.', 'strength', 'gym', 'intermediate', 60, ARRAY['штанга', 'гантели'], 1, 1, 1, true),
((SELECT id FROM programs WHERE slug='gym-intermediate-strength'), 'gs-w1d3', 'Тяжелый жим', 'Силовая тренировка с акцентом на жимовые движения.', 'strength', 'gym', 'intermediate', 60, ARRAY['штанга', 'гантели', 'скамья'], 2, 1, 3, true),
((SELECT id FROM programs WHERE slug='gym-intermediate-strength'), 'gs-w2d2', 'Тяжелый присед', 'Силовая тренировка ног с приседаниями.', 'strength', 'gym', 'intermediate', 65, ARRAY['штанга', 'гантели'], 3, 2, 2, true),

-- Program 5: Home beginner fitness
((SELECT id FROM programs WHERE slug='home-beginner-fitness'), 'hf-w1d1', 'Базовая тренировка', 'Комплексная тренировка для всего тела с простыми упражнениями.', 'endurance', 'home', 'beginner', 35, ARRAY['коврик'], 1, 1, 1, true),
((SELECT id FROM programs WHERE slug='home-beginner-fitness'), 'hf-w2d2', 'Кардио и выносливость', 'Аэробная тренировка для развития выносливости.', 'endurance', 'home', 'beginner', 30, ARRAY['коврик'], 2, 2, 2, true),
((SELECT id FROM programs WHERE slug='home-beginner-fitness'), 'hf-w3d1', 'Сила и гибкость', 'Сочетание силовых упражнений и растяжки.', 'endurance', 'home', 'beginner', 40, ARRAY['коврик'], 3, 3, 1, true),

-- Program 6: Gym advanced muscle gain
((SELECT id FROM programs WHERE slug='gym-advanced-muscle-gain'), 'ga-w1d1', 'Грудь — объёмная', 'Высокообъемная тренировка груди с дропсетами и суперсетами.', 'muscle_gain', 'gym', 'advanced', 70, ARRAY['штанга', 'гантели', 'скамья', 'кроссовер'], 1, 1, 1, true),
((SELECT id FROM programs WHERE slug='gym-advanced-muscle-gain'), 'ga-w2d2', 'Спина — ширина и толщина', 'Разделенная тренировка спины: вертикальные и горизонтальные тяги.', 'muscle_gain', 'gym', 'advanced', 75, ARRAY['турник', 'штанга', 'гантели', 'блочный тренажер'], 2, 2, 2, true),
((SELECT id FROM programs WHERE slug='gym-advanced-muscle-gain'), 'ga-w3d1', 'Ноги — полная проработка', 'Тяжелая тренировка ног с высоким объемом.', 'muscle_gain', 'gym', 'advanced', 80, ARRAY['штанга', 'гантели', 'жим ногами'], 3, 3, 1, true),
((SELECT id FROM programs WHERE slug='gym-advanced-muscle-gain'), 'ga-w4d3', 'Руки и плечи — пампинг', 'Объемная пампинг-тренировка рук и плеч с суперсетами.', 'muscle_gain', 'gym', 'advanced', 65, ARRAY['гантели', 'штанга', 'блочный тренажер'], 4, 4, 3, true);

-- ============================================================================
-- WORKOUT_EXERCISES (reference by slug subqueries)
-- ============================================================================

-- hw-w1d1: Кардио-разминка
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, duration_seconds, sort_order) VALUES
((SELECT id FROM workouts WHERE slug='hw-w1d1'), (SELECT id FROM exercises WHERE name='Приседания' LIMIT 1), 3, '15', 0, 1),
((SELECT id FROM workouts WHERE slug='hw-w1d1'), (SELECT id FROM exercises WHERE name='Берпи' LIMIT 1), 3, '8', 0, 2),
((SELECT id FROM workouts WHERE slug='hw-w1d1'), (SELECT id FROM exercises WHERE name='Горный альпинист' LIMIT 1), 3, '20', 0, 3),
((SELECT id FROM workouts WHERE slug='hw-w1d1'), (SELECT id FROM exercises WHERE name='Планка' LIMIT 1), 2, '30 сек', 30, 4);

-- hw-w1d3: Силовая всего тела
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, duration_seconds, sort_order) VALUES
((SELECT id FROM workouts WHERE slug='hw-w1d3'), (SELECT id FROM exercises WHERE name='Отжимания' LIMIT 1), 3, '10', 0, 1),
((SELECT id FROM workouts WHERE slug='hw-w1d3'), (SELECT id FROM exercises WHERE name='Выпады' LIMIT 1), 3, '12 на ногу', 0, 2),
((SELECT id FROM workouts WHERE slug='hw-w1d3'), (SELECT id FROM exercises WHERE name='Скручивания' LIMIT 1), 3, '15', 0, 3),
((SELECT id FROM workouts WHERE slug='hw-w1d3'), (SELECT id FROM exercises WHERE name='Боковая планка' LIMIT 1), 2, '30 сек', 30, 4);

-- hw-w2d1: Интервальная
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, duration_seconds, sort_order) VALUES
((SELECT id FROM workouts WHERE slug='hw-w2d1'), (SELECT id FROM exercises WHERE name='Берпи' LIMIT 1), 4, '10', 0, 1),
((SELECT id FROM workouts WHERE slug='hw-w2d1'), (SELECT id FROM exercises WHERE name='Приседания' LIMIT 1), 4, '20', 0, 2),
((SELECT id FROM workouts WHERE slug='hw-w2d1'), (SELECT id FROM exercises WHERE name='Горный альпинист' LIMIT 1), 4, '30 сек', 30, 3),
((SELECT id FROM workouts WHERE slug='hw-w2d1'), (SELECT id FROM exercises WHERE name='Отжимания' LIMIT 1), 3, '12', 0, 4);

-- gb-w1d1: Грудь и трицепс
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, duration_seconds, sort_order) VALUES
((SELECT id FROM workouts WHERE slug='gb-w1d1'), (SELECT id FROM exercises WHERE name='Жим гантелей лежа' LIMIT 1), 4, '10', 0, 1),
((SELECT id FROM workouts WHERE slug='gb-w1d1'), (SELECT id FROM exercises WHERE name='Отжимания' LIMIT 1), 3, '12', 0, 2),
((SELECT id FROM workouts WHERE slug='gb-w1d1'), (SELECT id FROM exercises WHERE name='Разведение гантелей лежа' LIMIT 1), 3, '12', 0, 3),
((SELECT id FROM workouts WHERE slug='gb-w1d1'), (SELECT id FROM exercises WHERE name='Французский жим лежа' LIMIT 1), 3, '12', 0, 4);

-- gb-w1d3: Спина и бицепс
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, duration_seconds, sort_order) VALUES
((SELECT id FROM workouts WHERE slug='gb-w1d3'), (SELECT id FROM exercises WHERE name='Подтягивания' LIMIT 1), 4, '8', 0, 1),
((SELECT id FROM workouts WHERE slug='gb-w1d3'), (SELECT id FROM exercises WHERE name='Тяга гантели в наклоне' LIMIT 1), 4, '12', 0, 2),
((SELECT id FROM workouts WHERE slug='gb-w1d3'), (SELECT id FROM exercises WHERE name='Сгибания рук с гантелями' LIMIT 1), 3, '12', 0, 3),
((SELECT id FROM workouts WHERE slug='gb-w1d3'), (SELECT id FROM exercises WHERE name='Гиперэкстензия' LIMIT 1), 3, '15', 0, 4);

-- gb-w2d2: Ноги
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, duration_seconds, sort_order) VALUES
((SELECT id FROM workouts WHERE slug='gb-w2d2'), (SELECT id FROM exercises WHERE name='Приседания с гантелями' LIMIT 1), 4, '12', 0, 1),
((SELECT id FROM workouts WHERE slug='gb-w2d2'), (SELECT id FROM exercises WHERE name='Выпады' LIMIT 1), 3, '12 на ногу', 0, 2),
((SELECT id FROM workouts WHERE slug='gb-w2d2'), (SELECT id FROM exercises WHERE name='Румынская тяга' LIMIT 1), 3, '12', 0, 3),
((SELECT id FROM workouts WHERE slug='gb-w2d2'), (SELECT id FROM exercises WHERE name='Подъемы на носки' LIMIT 1), 4, '15', 0, 4);

-- gb-w3d1: Плечи и пресс
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, duration_seconds, sort_order) VALUES
((SELECT id FROM workouts WHERE slug='gb-w3d1'), (SELECT id FROM exercises WHERE name='Жим гантелей стоя' LIMIT 1), 4, '10', 0, 1),
((SELECT id FROM workouts WHERE slug='gb-w3d1'), (SELECT id FROM exercises WHERE name='Махи гантелями в стороны' LIMIT 1), 3, '12', 0, 2),
((SELECT id FROM workouts WHERE slug='gb-w3d1'), (SELECT id FROM exercises WHERE name='Скручивания' LIMIT 1), 4, '20', 0, 3),
((SELECT id FROM workouts WHERE slug='gb-w3d1'), (SELECT id FROM exercises WHERE name='Планка' LIMIT 1), 3, '45 сек', 45, 4);

-- ============================================================================
-- REHAB COURSES
-- ============================================================================

INSERT INTO rehab_courses (slug, category, name, description, warnings, is_active, sort_order) VALUES
('hernia-recovery', 'hernia', 'Восстановление при грыже позвоночника', 'Программа реабилитации для людей с грыжами межпозвоночных дисков. Укрепление мышечного корсета, улучшение подвижности позвоночника и снижение боли. Постепенное увеличение нагрузки за 14 дней.', 'При острой боли или онемении конечностей обратитесь к врачу. Не выполняйте упражнения через сильную боль. Программа не заменяет медицинское лечение.', true, 1),
('protrusion-care', 'protrusion', 'Программа при протрузии дисков', 'Курс для людей с протрузиями межпозвоночных дисков. Укрепление мышц спины, улучшение осанки, снятие напряжения. Помогает предотвратить переход протрузии в грыжу.', 'Избегайте резких движений и осевой нагрузки на позвоночник. При усилении боли снизьте интенсивность. Обязательна консультация врача.', true, 2),
('scoliosis-correction', 'scoliosis', 'Коррекция сколиоза', 'Программа для коррекции сколиотической деформации. Асимметричные упражнения для балансировки мышечного тонуса. Работа над осанкой и симметрией тела.', 'При сколиозе 3-4 степени обязательна консультация ортопеда. Избегайте упражнений на турнике и висов. Контролируйте состояние у специалиста.', true, 3),
('kyphosis-fix', 'kyphosis', 'Исправление кифоза (сутулости)', 'Программа для коррекции грудного кифоза. Раскрытие грудной клетки, укрепление мышц верхней части спины, растяжка грудных мышц. Техники постурального контроля.', 'Не форсируйте растяжку — действуйте постепенно. При болях в грудном отделе проконсультируйтесь с врачом.', true, 4),
('lordosis-balance', 'lordosis', 'Работа с гиперлордозом', 'Курс для коррекции избыточного поясничного лордоза. Укрепление мышц живота и ягодиц, растяжка сгибателей бедра. Возвращение нейтрального положения таза.', 'Избегайте глубоких прогибов в пояснице. Держите пресс в напряжении. Регулярная практика важнее интенсивности.', true, 5);

-- ============================================================================
-- REHAB SESSIONS (14 days per course, stages 1→3)
-- ============================================================================

-- Hernia recovery
INSERT INTO rehab_sessions (course_id, day_number, stage, duration_minutes, description, sort_order) VALUES
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 1, 1, 15, 'Базовая растяжка и активация глубоких мышц спины. Дыхательные упражнения.', 1),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 2, 1, 15, 'Укрепление мышц кора лёжа на спине.', 2),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 3, 1, 20, 'Мобилизация позвоночника через мягкие движения.', 3),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 4, 1, 20, 'Изометрические упражнения для стабилизации позвоночника.', 4),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 5, 1, 20, 'Комплексная тренировка 1-й стадии.', 5),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 6, 2, 25, 'Увеличение амплитуды. Упражнения на четвереньках.', 6),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 7, 2, 25, 'Упражнения на баланс и координацию.', 7),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 8, 2, 25, 'Динамические упражнения для мышц спины.', 8),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 9, 2, 30, 'Функциональные упражнения для повседневных движений.', 9),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 10, 2, 30, 'Комплексная тренировка 2-й стадии.', 10),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 11, 3, 30, 'Упражнения стоя с увеличенной нагрузкой.', 11),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 12, 3, 35, 'Силовые упражнения для спины и кора.', 12),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 13, 3, 35, 'Комплексные многосуставные упражнения.', 13),
((SELECT id FROM rehab_courses WHERE slug='hernia-recovery'), 14, 3, 40, 'Финальная тренировка. Полный комплекс для самостоятельных занятий.', 14);

-- Protrusion care
INSERT INTO rehab_sessions (course_id, day_number, stage, duration_minutes, description, sort_order) VALUES
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 1, 1, 15, 'Мягкая разминка и декомпрессия позвоночника.', 1),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 2, 1, 15, 'Статические упражнения для укрепления мышц.', 2),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 3, 1, 18, 'Лёгкая растяжка мышц спины и ног.', 3),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 4, 1, 18, 'Улучшение питания дисков через мягкие движения.', 4),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 5, 1, 20, 'Базовый комплекс первого этапа.', 5),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 6, 2, 22, 'Укрепление боковых мышц спины.', 6),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 7, 2, 25, 'Работа над правильной осанкой.', 7),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 8, 2, 25, 'Динамическая стабилизация и баланс.', 8),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 9, 2, 28, 'Функциональный тренинг для спины.', 9),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 10, 2, 28, 'Комплекс второй стадии.', 10),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 11, 3, 30, 'Активные упражнения с контролируемой нагрузкой.', 11),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 12, 3, 32, 'Силовая выносливость мышц спины.', 12),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 13, 3, 35, 'Комплексные упражнения для спины и кора.', 13),
((SELECT id FROM rehab_courses WHERE slug='protrusion-care'), 14, 3, 35, 'Итоговая тренировка для самостоятельных занятий.', 14);

-- Scoliosis correction
INSERT INTO rehab_sessions (course_id, day_number, stage, duration_minutes, description, sort_order) VALUES
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 1, 1, 20, 'Оценка осанки. Базовые симметричные упражнения.', 1),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 2, 1, 20, 'Односторонние упражнения для балансировки тонуса.', 2),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 3, 1, 22, 'Укрепление ослабленной стороны.', 3),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 4, 1, 22, 'Осознание положения тела. Упражнения перед зеркалом.', 4),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 5, 1, 25, 'Дыхательные упражнения для раскрытия грудной клетки.', 5),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 6, 2, 25, 'Усиление асимметричной работы.', 6),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 7, 2, 28, 'Упражнения для деротации позвоночника.', 7),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 8, 2, 28, 'Силовые асимметричные упражнения для спины.', 8),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 9, 2, 30, 'Комплексная коррекция осанки.', 9),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 10, 2, 30, 'Функциональные упражнения с коррекцией.', 10),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 11, 3, 32, 'Продвинутые асимметричные упражнения.', 11),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 12, 3, 35, 'Длительное удержание корректирующих позиций.', 12),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 13, 3, 35, 'Динамические упражнения для закрепления.', 13),
((SELECT id FROM rehab_courses WHERE slug='scoliosis-correction'), 14, 3, 40, 'Полный корректирующий комплекс.', 14);

-- Kyphosis fix
INSERT INTO rehab_sessions (course_id, day_number, stage, duration_minutes, description, sort_order) VALUES
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 1, 1, 18, 'Мягкое раскрытие грудной клетки. Растяжка грудных мышц.', 1),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 2, 1, 18, 'Укрепление мышц между лопатками.', 2),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 3, 1, 20, 'Растяжка передней поверхности тела.', 3),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 4, 1, 20, 'Мягкая экстензия грудного отдела.', 4),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 5, 1, 22, 'Комплекс для улучшения осанки.', 5),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 6, 2, 25, 'Интенсивное укрепление верхней части спины.', 6),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 7, 2, 25, 'Работа с лопатками в различных плоскостях.', 7),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 8, 2, 28, 'Силовые упражнения для ромбовидных мышц.', 8),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 9, 2, 28, 'Функциональные упражнения для осанки.', 9),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 10, 2, 30, 'Комплексная тренировка 2-го этапа.', 10),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 11, 3, 30, 'Продвинутые упражнения для раскрытия груди.', 11),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 12, 3, 33, 'Силовая выносливость разгибателей грудного отдела.', 12),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 13, 3, 35, 'Интеграция осанки в повседневные движения.', 13),
((SELECT id FROM rehab_courses WHERE slug='kyphosis-fix'), 14, 3, 35, 'Финальный комплекс для поддержания осанки.', 14);

-- Lordosis balance
INSERT INTO rehab_sessions (course_id, day_number, stage, duration_minutes, description, sort_order) VALUES
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 1, 1, 18, 'Обучение нейтральному положению таза.', 1),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 2, 1, 18, 'Укрепление нижней части живота.', 2),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 3, 1, 20, 'Растяжка сгибателей бедра.', 3),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 4, 1, 20, 'Упражнения для ягодичных мышц.', 4),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 5, 1, 22, 'Базовый комплекс коррекции лордоза.', 5),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 6, 2, 25, 'Усиление работы мышц кора.', 6),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 7, 2, 25, 'Динамический контроль положения таза.', 7),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 8, 2, 28, 'Силовые упражнения для ягодиц.', 8),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 9, 2, 28, 'Интеграция нейтрального положения таза.', 9),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 10, 2, 30, 'Комплексная тренировка с контролем таза.', 10),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 11, 3, 30, 'Продвинутые упражнения для кора.', 11),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 12, 3, 33, 'Силовая выносливость мышц живота и ягодиц.', 12),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 13, 3, 35, 'Функциональный тренинг с коррекцией.', 13),
((SELECT id FROM rehab_courses WHERE slug='lordosis-balance'), 14, 3, 35, 'Итоговый комплекс для поддержания правильного положения.', 14);

-- ============================================================================
-- MEAL PLANS
-- ============================================================================

INSERT INTO meal_plans (slug, name, goal, day_number, calories, protein, fat, carbs, is_active, sort_order) VALUES
('weight-loss-plan', 'План питания для похудения', 'weight_loss', 1, 1500, 110.0, 50.0, 130.0, true, 1),
('muscle-gain-plan', 'План питания для набора массы', 'muscle_gain', 1, 2500, 180.0, 80.0, 280.0, true, 2),
('maintenance-plan', 'План питания для поддержания', 'maintenance', 1, 2000, 140.0, 65.0, 200.0, true, 3);

-- ============================================================================
-- MEALS
-- ============================================================================

INSERT INTO meals (meal_plan_id, meal_type, name, recipe, calories, protein, fat, carbs, alternatives, sort_order) VALUES
-- Weight loss plan
((SELECT id FROM meal_plans WHERE slug='weight-loss-plan'), 'breakfast', 'Овсяная каша с ягодами', 'Овсяные хлопья (50г) залить кипятком, настоять 5 мин. Добавить ягоды (100г), грецкие орехи (10г), корицу.', 320, 12.0, 10.0, 48.0, 'Гречневая каша с кефиром, творог с фруктами, омлет из 2 яиц с овощами', 1),
((SELECT id FROM meal_plans WHERE slug='weight-loss-plan'), 'lunch', 'Куриная грудка с гречкой', 'Куриная грудка (150г) запечённая с травами. Гречка (80г сухой). Салат из огурцов и помидоров (150г) с оливковым маслом.', 450, 48.0, 12.0, 42.0, 'Индейка или рыба, бурый рис или киноа вместо гречки', 2),
((SELECT id FROM meal_plans WHERE slug='weight-loss-plan'), 'dinner', 'Запечённая рыба с овощами', 'Филе белой рыбы (180г) с лимоном и травами. Овощи запечённые (200г): кабачки, баклажаны, перец. Зелень.', 380, 40.0, 14.0, 18.0, 'Куриные котлеты на пару, морепродукты, мясо кролика', 3),
((SELECT id FROM meal_plans WHERE slug='weight-loss-plan'), 'snack', 'Творог с зеленью', 'Обезжиренный творог (150г) с укропом и петрушкой. Щепотка соли. Можно с цельнозерновым хлебцем.', 150, 25.0, 2.0, 8.0, 'Греческий йогурт с орехами, яблоко с арахисовой пастой', 4),

-- Muscle gain plan
((SELECT id FROM meal_plans WHERE slug='muscle-gain-plan'), 'breakfast', 'Омлет с сыром и тостами', 'Омлет из 4 яиц с молоком (50мл) и сыром (30г). Цельнозерновые тосты (2шт) с авокадо (50г). Помидоры черри.', 580, 38.0, 30.0, 40.0, 'Блины из овсянки с творогом, сырники с мёдом', 1),
((SELECT id FROM meal_plans WHERE slug='muscle-gain-plan'), 'lunch', 'Говядина с макаронами', 'Говядина постная (200г) тушёная. Макароны из твёрдых сортов (100г сухих). Овощной салат (150г) с оливковым маслом.', 720, 52.0, 22.0, 78.0, 'Курица с рисом, индейка с картофелем', 2),
((SELECT id FROM meal_plans WHERE slug='muscle-gain-plan'), 'dinner', 'Лосось с киноа', 'Стейк лосося (200г) на гриле. Киноа (100г сухой). Тушёные овощи (200г): брокколи, цветная капуста, морковь.', 680, 48.0, 26.0, 62.0, 'Тунец с бурым рисом, курица с бататом', 3),
((SELECT id FROM meal_plans WHERE slug='muscle-gain-plan'), 'snack', 'Творожная запеканка с бананом', 'Творог 5% (200г), яйца (2шт), манка (30г), мёд (20г). Запечь в духовке. Подавать с бананом и ягодами.', 520, 42.0, 2.0, 100.0, 'Протеиновый коктейль с бананом и овсянкой, орехи с сухофруктами', 4),

-- Maintenance plan
((SELECT id FROM meal_plans WHERE slug='maintenance-plan'), 'breakfast', 'Мюсли с йогуртом и фруктами', 'Мюсли (60г) с греческим йогуртом (150г). Банан, яблоко, ягоды. Семена чиа (5г).', 420, 18.0, 12.0, 58.0, 'Смузи-боул, овсянка с протеином, тосты с авокадо и яйцом', 1),
((SELECT id FROM meal_plans WHERE slug='maintenance-plan'), 'lunch', 'Индейка с бурым рисом', 'Грудка индейки (160г) на пару или гриле. Бурый рис (80г сухого). Овощной салат (200г) с лимонным соусом.', 520, 44.0, 14.0, 58.0, 'Куриные котлеты с гречкой, рыба с картофелем', 2),
((SELECT id FROM meal_plans WHERE slug='maintenance-plan'), 'dinner', 'Куриное филе с овощным рагу', 'Куриное филе (150г) с розмарином. Рагу из кабачков, баклажанов, помидоров, перца и лука (250г). Зелень.', 460, 42.0, 16.0, 32.0, 'Рыба с овощами на пару, морепродукты с овощами', 3),
((SELECT id FROM meal_plans WHERE slug='maintenance-plan'), 'snack', 'Сэндвич с курицей', 'Цельнозерновой хлеб (2 ломтика), куриная грудка (80г), салат, помидор, огурец. Йогуртовый соус.', 400, 30.0, 12.0, 42.0, 'Хлебцы с творожным сыром и рыбой, йогурт с гранолой', 4);

COMMIT;
