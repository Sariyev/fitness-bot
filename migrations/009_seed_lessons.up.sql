BEGIN;

-- ==================== ЛФК: Грыжа ====================
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'hernia_basics', 'Основы ЛФК при грыже', 'Базовые упражнения и правила занятий при грыже позвоночника', 1
FROM module_categories c WHERE c.slug = 'hernia'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'hernia_stretching', 'Растяжка при грыже', 'Упражнения на растяжку для снятия напряжения и боли', 2
FROM module_categories c WHERE c.slug = 'hernia'
ON CONFLICT (category_id, slug) DO NOTHING;

-- ЛФК: Протрузии
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'protrusion_basics', 'Основы ЛФК при протрузиях', 'Безопасные упражнения для укрепления спины', 1
FROM module_categories c WHERE c.slug = 'protrusion'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'protrusion_strength', 'Укрепление при протрузиях', 'Упражнения для стабилизации позвоночника', 2
FROM module_categories c WHERE c.slug = 'protrusion'
ON CONFLICT (category_id, slug) DO NOTHING;

-- ЛФК: Сколиоз
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'scoliosis_basics', 'Основы ЛФК при сколиозе', 'Корректирующие упражнения для выравнивания позвоночника', 1
FROM module_categories c WHERE c.slug = 'scoliosis'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'scoliosis_correction', 'Коррекция осанки при сколиозе', 'Асимметричные упражнения для исправления искривления', 2
FROM module_categories c WHERE c.slug = 'scoliosis'
ON CONFLICT (category_id, slug) DO NOTHING;

-- ЛФК: Кифоз
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'kyphosis_basics', 'Основы ЛФК при кифозе', 'Упражнения для раскрытия грудного отдела', 1
FROM module_categories c WHERE c.slug = 'kyphosis'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'kyphosis_posture', 'Осанка при кифозе', 'Упражнения для укрепления мышц спины и коррекции сутулости', 2
FROM module_categories c WHERE c.slug = 'kyphosis'
ON CONFLICT (category_id, slug) DO NOTHING;

-- ЛФК: Лордоз
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'lordosis_basics', 'Основы ЛФК при лордозе', 'Упражнения для стабилизации поясничного отдела', 1
FROM module_categories c WHERE c.slug = 'lordosis'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'lordosis_core', 'Укрепление кора при лордозе', 'Упражнения для мышц пресса и поясницы', 2
FROM module_categories c WHERE c.slug = 'lordosis'
ON CONFLICT (category_id, slug) DO NOTHING;

-- ==================== Тренировки: Грудь ====================
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'chest_basics', 'Базовая тренировка груди', 'Жим лёжа, отжимания и разводки для начинающих', 1
FROM module_categories c WHERE c.slug = 'chest'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'chest_advanced', 'Продвинутая тренировка груди', 'Суперсеты и дроп-сеты для опытных', 2
FROM module_categories c WHERE c.slug = 'chest'
ON CONFLICT (category_id, slug) DO NOTHING;

-- Тренировки: Спина
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'back_basics', 'Базовая тренировка спины', 'Тяга, подтягивания и гиперэкстензия', 1
FROM module_categories c WHERE c.slug = 'back'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'back_advanced', 'Продвинутая тренировка спины', 'Тяга штанги в наклоне и тяга блока', 2
FROM module_categories c WHERE c.slug = 'back'
ON CONFLICT (category_id, slug) DO NOTHING;

-- Тренировки: Ноги
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'legs_basics', 'Базовая тренировка ног', 'Приседания, выпады и жим ногами', 1
FROM module_categories c WHERE c.slug = 'legs'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'legs_advanced', 'Продвинутая тренировка ног', 'Румынская тяга и болгарские выпады', 2
FROM module_categories c WHERE c.slug = 'legs'
ON CONFLICT (category_id, slug) DO NOTHING;

-- Тренировки: Плечи
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'shoulders_basics', 'Базовая тренировка плеч', 'Жим гантелей и махи в стороны', 1
FROM module_categories c WHERE c.slug = 'shoulders'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'shoulders_advanced', 'Продвинутая тренировка плеч', 'Протяжка и обратные разведения', 2
FROM module_categories c WHERE c.slug = 'shoulders'
ON CONFLICT (category_id, slug) DO NOTHING;

-- Тренировки: Руки
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'arms_basics', 'Базовая тренировка рук', 'Подъём на бицепс и разгибание на трицепс', 1
FROM module_categories c WHERE c.slug = 'arms'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'arms_advanced', 'Продвинутая тренировка рук', 'Суперсеты бицепс-трицепс', 2
FROM module_categories c WHERE c.slug = 'arms'
ON CONFLICT (category_id, slug) DO NOTHING;

-- Тренировки: Пресс
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'core_basics', 'Базовая тренировка пресса', 'Скручивания, планка и подъём ног', 1
FROM module_categories c WHERE c.slug = 'core'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'core_advanced', 'Продвинутая тренировка пресса', 'Вакуум, велосипед и боковая планка', 2
FROM module_categories c WHERE c.slug = 'core'
ON CONFLICT (category_id, slug) DO NOTHING;

-- ==================== Питание ====================
-- Завтраки
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'breakfast_oatmeal', 'Овсянка с фруктами', 'Полезный завтрак с медленными углеводами', 1
FROM module_categories c WHERE c.slug = 'breakfast'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'breakfast_eggs', 'Яичница с овощами', 'Белковый завтрак для активного дня', 2
FROM module_categories c WHERE c.slug = 'breakfast'
ON CONFLICT (category_id, slug) DO NOTHING;

-- Обеды
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'lunch_chicken', 'Куриная грудка с рисом', 'Классический обед для набора массы', 1
FROM module_categories c WHERE c.slug = 'lunch'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'lunch_fish', 'Рыба с овощами', 'Лёгкий обед с омега-3', 2
FROM module_categories c WHERE c.slug = 'lunch'
ON CONFLICT (category_id, slug) DO NOTHING;

-- Ужины
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'dinner_salad', 'Салат с тунцом', 'Лёгкий ужин с высоким содержанием белка', 1
FROM module_categories c WHERE c.slug = 'dinner'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'dinner_soup', 'Куриный суп', 'Питательный ужин для восстановления', 2
FROM module_categories c WHERE c.slug = 'dinner'
ON CONFLICT (category_id, slug) DO NOTHING;

-- Перекусы
INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'snack_nuts', 'Орехи и сухофрукты', 'Быстрый перекус с полезными жирами', 1
FROM module_categories c WHERE c.slug = 'snacks'
ON CONFLICT (category_id, slug) DO NOTHING;

INSERT INTO lessons (category_id, slug, title, description, sort_order)
SELECT c.id, 'snack_smoothie', 'Протеиновый смузи', 'Перекус после тренировки', 2
FROM module_categories c WHERE c.slug = 'snacks'
ON CONFLICT (category_id, slug) DO NOTHING;

-- ==================== Lesson Contents ====================

-- Грыжа: Основы
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Что такое межпозвоночная грыжа',
'Межпозвоночная грыжа — это выпячивание ядра межпозвоночного диска за пределы фиброзного кольца. ЛФК помогает укрепить мышечный корсет и снять нагрузку с позвоночника.

**Основные правила:**
• Не допускай резких движений
• Начинай с малой амплитуды
• При боли — немедленно прекрати
• Занимайся регулярно, но не перегружайся
• Делай упражнения на мягком коврике', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'hernia_basics' AND c.slug = 'hernia';

INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Упражнение: Кошка-корова',
'**Исходное положение:** встань на четвереньки, колени под тазом, руки под плечами.

1. На вдохе прогни спину вниз, подними голову (поза коровы)
2. На выдохе округли спину, опусти голову (поза кошки)
3. Повтори 10-15 раз

**Важно:** движения плавные, без рывков. Дыши ровно.', 2
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'hernia_basics' AND c.slug = 'hernia';

-- Грыжа: Растяжка
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Растяжка поясничного отдела',
'**Упражнение 1: Колени к груди**
Ляг на спину, подтяни оба колена к груди, обхвати руками. Удерживай 20-30 секунд.

**Упражнение 2: Поворот лёжа**
Ляг на спину, согни колени. Плавно опусти колени влево, задержись на 20 секунд, затем вправо.

**Упражнение 3: Поза ребёнка**
Сядь на пятки, наклонись вперёд, руки вытяни перед собой. Удерживай 30-60 секунд.

Делай каждое упражнение по 3 подхода.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'hernia_stretching' AND c.slug = 'hernia';

-- Протрузии: Основы
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'ЛФК при протрузиях',
'Протрузия — это начальная стадия грыжи, когда диск выпячивается, но фиброзное кольцо ещё цело. Правильные упражнения могут остановить прогрессирование.

**Что нельзя делать:**
• Поднимать тяжести с круглой спиной
• Делать резкие скручивания
• Прыгать и бегать по жёсткой поверхности

**Упражнение: Мост**
Ляг на спину, согни колени. Подними таз вверх, задержись на 5 секунд. 15 повторений × 3 подхода.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'protrusion_basics' AND c.slug = 'protrusion';

-- Протрузии: Укрепление
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Укрепление мышц спины',
'**Упражнение 1: Супермен**
Ляг на живот, руки вперёд. Одновременно подними руки и ноги, задержись на 3 секунды. 12 повторений × 3 подхода.

**Упражнение 2: Планка**
Встань на предплечья и носки. Держи тело ровно 30-60 секунд. 3 подхода.

**Упражнение 3: Птица-собака**
Встань на четвереньки. Вытяни правую руку и левую ногу, задержись на 5 секунд. Поменяй стороны. 10 повторений на каждую сторону.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'protrusion_strength' AND c.slug = 'protrusion';

-- Сколиоз: Основы
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Коррекция сколиоза упражнениями',
'Сколиоз — боковое искривление позвоночника. ЛФК направлена на укрепление слабых мышц и растяжку укороченных.

**Упражнение: Боковая планка**
Встань на предплечье, тело ровное. Удерживай 20-30 секунд на каждую сторону. Начни с той стороны, которая слабее.

**Упражнение: Растяжка в дверном проёме**
Встань в дверной проём, руки на уровне плеч. Сделай шаг вперёд, почувствуй растяжение в груди. Удерживай 30 секунд.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'scoliosis_basics' AND c.slug = 'scoliosis';

-- Сколиоз: Коррекция
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Асимметричные упражнения',
'При сколиозе важно делать упражнения с акцентом на слабую сторону.

**Упражнение 1: Тяга одной рукой**
Наклонись, одна рука на скамье. Тяни гантель к поясу рабочей рукой. 12 повторений × 3 подхода на слабую сторону, 2 на сильную.

**Упражнение 2: Боковые наклоны**
Стоя, наклоняйся в сторону выпуклости дуги. 15 повторений × 3 подхода.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'scoliosis_correction' AND c.slug = 'scoliosis';

-- Кифоз: Основы
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Упражнения при кифозе',
'Кифоз — чрезмерный прогиб грудного отдела (сутулость). Цель ЛФК: раскрыть грудную клетку и укрепить разгибатели спины.

**Упражнение: Раскрытие груди на валике**
Ляг на пенный ролик (или скрученное полотенце), расположи его под грудным отделом. Руки за голову, позволь спине прогнуться. 2-3 минуты.

**Упражнение: Тяга резинки к лицу**
Закрепи резинку на уровне лица. Тяни к лицу, разводя локти в стороны. 15 повторений × 3 подхода.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'kyphosis_basics' AND c.slug = 'kyphosis';

-- Кифоз: Осанка
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Коррекция осанки',
'**Упражнение 1: Y-подъёмы лёжа**
Ляг на живот, руки вытянуты в форме буквы Y. Подними руки и верх тела, сведи лопатки. 12 повторений × 3 подхода.

**Упражнение 2: Стена**
Встань спиной к стене: пятки, таз, лопатки и затылок касаются стены. Удерживай 1-2 минуты, повторяй 3 раза в день.

**Важно:** Следи за осанкой в течение дня. Каждый час делай перерыв и расправляй плечи.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'kyphosis_posture' AND c.slug = 'kyphosis';

-- Лордоз: Основы
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Упражнения при лордозе',
'Лордоз — чрезмерный прогиб в поясничном отделе. Цель: укрепить пресс и растянуть сгибатели бедра.

**Упражнение: Наклон таза лёжа**
Ляг на спину, согни колени. Прижми поясницу к полу, напрягая пресс. Удерживай 5 секунд. 15 повторений × 3 подхода.

**Упражнение: Растяжка квадрицепса**
Стоя, согни ногу назад и возьмись за стопу. Подтяни пятку к ягодице. 30 секунд на каждую ногу.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'lordosis_basics' AND c.slug = 'lordosis';

-- Лордоз: Кор
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Укрепление кора при лордозе',
'**Упражнение 1: Мёртвый жук**
Ляг на спину, руки вверх, колени согнуты на 90°. Выпрями правую руку за голову и левую ногу вперёд одновременно. Вернись. 10 повторений на каждую сторону × 3 подхода.

**Упражнение 2: Обратная планка**
Сядь, руки за спиной. Подними таз, тело ровное от плеч до пяток. 20-30 секунд × 3 подхода.

**Упражнение 3: Скручивания с фиксацией поясницы**
Ляг на спину, колени согнуты. Подними плечи, не отрывая поясницу. 15 повторений × 3 подхода.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'lordosis_core' AND c.slug = 'lordosis';

-- Грудь: Базовая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Базовые упражнения на грудные мышцы',
'**1. Жим лёжа со штангой**
3 подхода × 8-12 повторений. Опускай штангу к середине груди, жми вверх.

**2. Отжимания от пола**
3 подхода × 15-20 повторений. Руки чуть шире плеч, тело ровное.

**3. Разводка гантелей лёжа**
3 подхода × 12-15 повторений. Руки слегка согнуты, опускай до уровня груди.

**Отдых между подходами:** 60-90 секунд.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'chest_basics' AND c.slug = 'chest';

-- Грудь: Продвинутая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Продвинутые методы тренировки груди',
'**Суперсет 1:** Жим гантелей лёжа + отжимания (без отдыха между упражнениями)
4 подхода × 10 + до отказа

**Суперсет 2:** Жим на наклонной скамье + разводка на горизонтальной
4 подхода × 10 + 12

**Дроп-сет:** Жим в тренажёре — начни с тяжёлого веса, 8 повторений, сбрось вес на 20%, ещё 8, сбрось ещё, до отказа.

**Отдых:** 60 секунд между суперсетами.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'chest_advanced' AND c.slug = 'chest';

-- Спина: Базовая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Базовые упражнения на спину',
'**1. Подтягивания**
3 подхода × максимум повторений. Хват чуть шире плеч.

**2. Тяга верхнего блока**
3 подхода × 10-12 повторений. Тяни к верху груди, сводя лопатки.

**3. Гиперэкстензия**
3 подхода × 15 повторений. Поднимайся до прямой линии с ногами.

**Совет:** Всегда начинай тренировку спины с подтягиваний — самое эффективное упражнение.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'back_basics' AND c.slug = 'back';

-- Спина: Продвинутая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Продвинутая тренировка спины',
'**1. Тяга штанги в наклоне**
4 подхода × 8-10 повторений. Наклон 45°, тяни к поясу.

**2. Тяга гантели одной рукой**
3 подхода × 10-12 на каждую руку. Опирайся коленом на скамью.

**3. Пуловер с гантелей**
3 подхода × 12-15 повторений. Чувствуй растяжение широчайших.

**4. Шраги со штангой**
3 подхода × 15 повторений. Поднимай плечи вверх, сжимай трапеции.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'back_advanced' AND c.slug = 'back';

-- Ноги: Базовая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Базовые упражнения на ноги',
'**1. Приседания со штангой**
4 подхода × 8-12 повторений. Колени за линию носков, спина прямая.

**2. Выпады с гантелями**
3 подхода × 12 на каждую ногу. Шаг вперёд, колено 90°.

**3. Жим ногами в тренажёре**
3 подхода × 12-15 повторений. Ступни на ширине плеч.

**Важно:** День ног — самый энергозатратный. Хорошо разомнись перед тренировкой.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'legs_basics' AND c.slug = 'legs';

-- Ноги: Продвинутая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Продвинутая тренировка ног',
'**1. Румынская тяга**
4 подхода × 10-12. Штанга скользит по ногам, чувствуй заднюю поверхность бедра.

**2. Болгарские выпады**
3 подхода × 10 на каждую ногу. Задняя нога на скамье.

**3. Сгибание ног в тренажёре**
3 подхода × 12-15 повторений.

**4. Подъём на носки**
4 подхода × 15-20 повторений. Полная амплитуда — растяжение в нижней точке.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'legs_advanced' AND c.slug = 'legs';

-- Плечи: Базовая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Базовые упражнения на плечи',
'**1. Жим гантелей сидя**
3 подхода × 10-12 повторений. Опускай до уровня ушей.

**2. Махи гантелями в стороны**
3 подхода × 12-15 повторений. Слегка согни локти, поднимай до уровня плеч.

**3. Подъём гантелей перед собой**
3 подхода × 12 повторений. Поочерёдно каждой рукой.

**Совет:** Плечи — маленькая мышечная группа. Не бери слишком тяжёлые веса.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'shoulders_basics' AND c.slug = 'shoulders';

-- Плечи: Продвинутая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Продвинутая тренировка плеч',
'**1. Армейский жим стоя**
4 подхода × 8-10 повторений. Штанга от груди вверх.

**2. Протяжка штанги к подбородку**
3 подхода × 12 повторений. Узкий хват, тяни локти вверх.

**3. Обратные разведения в наклоне**
3 подхода × 15 повторений. Наклон 90°, разводи руки в стороны.

**4. Шраги с гантелями**
3 подхода × 15 повторений.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'shoulders_advanced' AND c.slug = 'shoulders';

-- Руки: Базовая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Базовые упражнения на руки',
'**Бицепс:**
1. Подъём штанги на бицепс стоя — 3×12
2. Сгибание рук с гантелями «молот» — 3×12

**Трицепс:**
1. Разгибание рук на блоке — 3×12
2. Французский жим лёжа — 3×10

**Отдых:** 60 секунд между подходами. Следи за полной амплитудой движения.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'arms_basics' AND c.slug = 'arms';

-- Руки: Продвинутая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Суперсеты бицепс-трицепс',
'**Суперсет 1:**
Подъём на бицепс со штангой + отжимания на брусьях
4 подхода × 10 + 10

**Суперсет 2:**
Сгибание рук на наклонной скамье + разгибание из-за головы
3 подхода × 12 + 12

**Суперсет 3:**
Концентрированный подъём + кикбэк с гантелей
3 подхода × 10 + 10

Без отдыха внутри суперсета, 90 секунд между суперсетами.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'arms_advanced' AND c.slug = 'arms';

-- Пресс: Базовая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Базовые упражнения на пресс',
'**1. Скручивания**
3 подхода × 20 повторений. Поднимай только лопатки, не тяни шею руками.

**2. Планка**
3 подхода × 30-60 секунд. Тело ровное, не проваливай поясницу.

**3. Подъём ног лёжа**
3 подхода × 15 повторений. Медленно опускай ноги, не касаясь пола.

**Делай эту тренировку 3-4 раза в неделю для лучших результатов.**', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'core_basics' AND c.slug = 'core';

-- Пресс: Продвинутая
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Продвинутые упражнения на пресс',
'**1. Вакуум**
Втяни живот максимально, удерживай 15-30 секунд. 5 подходов. Делай натощак.

**2. Велосипед**
3 подхода × 20 на каждую сторону. Касайся локтем противоположного колена.

**3. Боковая планка**
3 подхода × 30 секунд на каждую сторону.

**4. Скалолаз**
3 подхода × 30 секунд. Быстро подтягивай колени к груди в планке.', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'core_advanced' AND c.slug = 'core';

-- Завтраки: Овсянка
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Рецепт: Овсянка с фруктами и орехами',
'**Ингредиенты:**
• 80 г овсяных хлопьев
• 200 мл молока или воды
• 1 банан
• Горсть ягод (клубника, черника)
• 15 г орехов (грецкие или миндаль)
• 1 ч.л. мёда

**Приготовление:**
1. Свари овсянку на молоке (3-5 минут)
2. Нарежь банан, добавь сверху
3. Посыпь ягодами и орехами
4. Полей мёдом

**КБЖУ:** ~420 ккал | Б: 14 г | Ж: 12 г | У: 65 г', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'breakfast_oatmeal' AND c.slug = 'breakfast';

-- Завтраки: Яичница
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Рецепт: Яичница с овощами',
'**Ингредиенты:**
• 3 яйца
• 1 помидор
• 1/2 болгарского перца
• Шпинат (горсть)
• 1 ч.л. оливкового масла
• Соль, перец по вкусу

**Приготовление:**
1. Нарежь овощи
2. Разогрей масло на сковороде
3. Обжарь перец 2 минуты, добавь помидор и шпинат
4. Залей яйцами, накрой крышкой
5. Готовь 3-4 минуты на среднем огне

**КБЖУ:** ~320 ккал | Б: 22 г | Ж: 20 г | У: 8 г', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'breakfast_eggs' AND c.slug = 'breakfast';

-- Обеды: Курица
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Рецепт: Куриная грудка с рисом и овощами',
'**Ингредиенты:**
• 200 г куриной грудки
• 80 г риса (сухой вес)
• Брокколи — 100 г
• 1 ч.л. оливкового масла
• Специи: паприка, чеснок, соль

**Приготовление:**
1. Свари рис (15-20 минут)
2. Нарежь курицу, обваляй в специях
3. Обжарь на масле 5-7 минут
4. Отвари брокколи на пару 5 минут
5. Сервируй на тарелке

**КБЖУ:** ~480 ккал | Б: 45 г | Ж: 8 г | У: 60 г', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'lunch_chicken' AND c.slug = 'lunch';

-- Обеды: Рыба
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Рецепт: Запечённая рыба с овощами',
'**Ингредиенты:**
• 200 г филе лосося (или форели)
• 1 цукини
• 1 морковь
• Лимон — 2 дольки
• 1 ч.л. оливкового масла
• Укроп, соль, перец

**Приготовление:**
1. Нарежь овощи кружочками
2. Выложи на противень с фольгой
3. Положи рыбу сверху, полей маслом и лимоном
4. Запекай при 180°C 20-25 минут

**КБЖУ:** ~380 ккал | Б: 38 г | Ж: 18 г | У: 12 г', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'lunch_fish' AND c.slug = 'lunch';

-- Ужины: Салат
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Рецепт: Салат с тунцом',
'**Ингредиенты:**
• 1 банка тунца в собственном соку
• Листья салата — большая горсть
• 1 огурец
• 5 помидоров черри
• 1/2 авокадо
• 1 ст.л. оливкового масла
• Лимонный сок, соль

**Приготовление:**
1. Разложи листья на тарелке
2. Нарежь огурец, помидоры и авокадо
3. Выложи тунец в центр
4. Заправь маслом и лимоном

**КБЖУ:** ~350 ккал | Б: 35 г | Ж: 18 г | У: 10 г', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'dinner_salad' AND c.slug = 'dinner';

-- Ужины: Суп
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Рецепт: Куриный суп с овощами',
'**Ингредиенты:**
• 300 г куриного филе
• 2 картофелины
• 1 морковь
• 1 луковица
• 1.5 л воды
• Зелень, соль, перец

**Приготовление:**
1. Свари курицу 20 минут, вынь, нарежь
2. В бульон добавь нарезанный картофель и морковь
3. Обжарь лук, добавь в суп
4. Вари 15 минут
5. Верни курицу, добавь зелень

**КБЖУ:** ~280 ккал | Б: 30 г | Ж: 5 г | У: 30 г', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'dinner_soup' AND c.slug = 'dinner';

-- Перекусы: Орехи
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Орехи и сухофрукты — правила перекуса',
'**Оптимальная порция:** 30-40 г орехов + 20 г сухофруктов.

**Лучшие орехи для спортсменов:**
• Миндаль — витамин E и магний
• Грецкие — омега-3 жирные кислоты
• Кешью — железо и цинк
• Фундук — витамины группы B

**Сухофрукты:**
• Курага — калий для сердца
• Финики — быстрая энергия перед тренировкой
• Изюм — железо

**Важно:** Не превышай порцию — орехи калорийны (~550 ккал на 100 г).

**КБЖУ порции:** ~200 ккал | Б: 6 г | Ж: 14 г | У: 14 г', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'snack_nuts' AND c.slug = 'snacks';

-- Перекусы: Смузи
INSERT INTO lesson_contents (lesson_id, content_type, title, body, sort_order)
SELECT l.id, 'text', 'Рецепт: Протеиновый смузи',
'**Ингредиенты:**
• 1 скуп протеина (ваниль или шоколад)
• 1 банан
• 200 мл молока
• 1 ст.л. арахисовой пасты
• 3-4 кубика льда

**Приготовление:**
1. Загрузи всё в блендер
2. Взбивай 30-40 секунд
3. Перелей в стакан

**Когда пить:** В течение 30 минут после тренировки.

**КБЖУ:** ~380 ккал | Б: 32 г | Ж: 12 г | У: 40 г', 1
FROM lessons l JOIN module_categories c ON l.category_id = c.id WHERE l.slug = 'snack_smoothie' AND c.slug = 'snacks';

COMMIT;
