-- +goose Up
ALTER TABLE mission_reports
DROP CONSTRAINT mission_reports_type_valid;

ALTER TABLE mission_reports
ADD CONSTRAINT mission_reports_type_valid CHECK (
    type IN (
        'pve_mission',
        'pvp_raid_attacker',
        'pvp_raid_defender',
        'event'
    )
);

CREATE TABLE game_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_key TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    trigger_type TEXT NOT NULL DEFAULT 'pool',
    weight INTEGER NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT true,
    cooldown_seconds INTEGER NOT NULL DEFAULT 21600,
    expires_after_seconds INTEGER NOT NULL DEFAULT 86400,
    conditions_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT game_events_category_valid CHECK (
        category IN ('economy', 'ruler', 'military', 'patron', 'dark_omen')
    ),
    CONSTRAINT game_events_trigger_type_valid CHECK (trigger_type IN ('pool')),
    CONSTRAINT game_events_weight_positive CHECK (weight > 0),
    CONSTRAINT game_events_cooldown_nonnegative CHECK (cooldown_seconds >= 0),
    CONSTRAINT game_events_expires_after_positive CHECK (expires_after_seconds > 0)
);

CREATE INDEX game_events_event_key_idx ON game_events (event_key);
CREATE INDEX game_events_category_idx ON game_events (category);
CREATE INDEX game_events_is_active_idx ON game_events (is_active);

CREATE TABLE event_choices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_event_id UUID NOT NULL REFERENCES game_events(id) ON DELETE CASCADE,
    choice_key TEXT NOT NULL,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    effects_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    result_title TEXT NOT NULL,
    result_body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT event_choices_game_event_choice_unique UNIQUE (game_event_id, choice_key)
);

CREATE INDEX event_choices_game_event_id_idx ON event_choices (game_event_id);

CREATE TABLE kingdom_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL REFERENCES kingdoms(id) ON DELETE CASCADE,
    game_event_id UUID NOT NULL REFERENCES game_events(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'active',
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    resolved_at TIMESTAMPTZ NULL,
    selected_choice_key TEXT NULL,
    result_json JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT kingdom_events_status_valid CHECK (status IN ('active', 'resolved', 'expired'))
);

CREATE INDEX kingdom_events_kingdom_id_idx ON kingdom_events (kingdom_id);
CREATE INDEX kingdom_events_game_event_id_idx ON kingdom_events (game_event_id);
CREATE INDEX kingdom_events_status_idx ON kingdom_events (status);
CREATE INDEX kingdom_events_expires_at_idx ON kingdom_events (expires_at);
CREATE INDEX kingdom_events_generated_at_idx ON kingdom_events (generated_at);

WITH inserted AS (
    INSERT INTO game_events (event_key, category, title, body, conditions_json)
    VALUES
        ('found_old_idol', 'economy', 'Старый идол в лесу', 'Лесорубы нашли в корнях чёрный идол. Люди спорят, продать его купцам или отдать волхвам.', '{}'::jsonb),
        ('ruler_bad_dream', 'ruler', 'Сон правителя', 'Правитель проснулся до рассвета. Ему снился волк без головы у городских ворот.', '{}'::jsonb),
        ('volunteers_at_gate', 'military', 'Добровольцы у ворот', 'Несколько молодых людей пришли к воротам и просят взять их в ополчение.', '{}'::jsonb),
        ('patron_envoy', 'patron', 'Посол покровителя', 'К воротам прибыл посол. Он говорит мягко, но печать на грамоте тяжёлая.', '{"requiresPatron":"any"}'::jsonb),
        ('black_birds_over_walls', 'dark_omen', 'Чёрные птицы над стенами', 'На стенах всю ночь сидели чёрные птицы. Утром люди нашли следы когтей на дверях амбара.', '{}'::jsonb)
    RETURNING id, event_key
)
INSERT INTO event_choices (game_event_id, choice_key, label, description, effects_json, result_title, result_body)
SELECT id, 'sell_to_merchants', 'Продать купцам', 'Купцы не задают вопросов, если цена хорошая.', '{"resourceDelta":{"gold":80},"kingdomDelta":{"honor":-1}}'::jsonb, 'Идол ушёл с купцами', 'Купцы забрали находку без лишних слов. В казне стало тяжелее, но старики смотрят на лес тревожнее.'
FROM inserted WHERE event_key = 'found_old_idol'
UNION ALL
SELECT id, 'give_to_volkhvs', 'Отдать волхвам', 'Волхвы примут знак и накормят людей после обряда.', '{"resourceDelta":{"food":30},"kingdomDelta":{"honor":1}}'::jsonb, 'Волхвы приняли идол', 'Идол исчез в дыму трав. Люди говорят тише, зато двор выглядит благочестивее.'
FROM inserted WHERE event_key = 'found_old_idol'
UNION ALL
SELECT id, 'calm_the_court', 'Успокоить двор', 'Сказать, что сон не властен над живыми.', '{"kingdomDelta":{"honor":1}}'::jsonb, 'Двор успокоен', 'Правитель вышел к людям спокойным. Страх отступил, а слово двора стало твёрже.'
FROM inserted WHERE event_key = 'ruler_bad_dream'
UNION ALL
SELECT id, 'use_the_omen', 'Использовать знамение', 'Пусть враги услышат, что даже сны служат твоей воле.', '{"kingdomDelta":{"dread":1}}'::jsonb, 'Знамение стало оружием', 'Сон пересказали так, что соседи вспомнили о воротах и замках. Во дворе стало тише.'
FROM inserted WHERE event_key = 'ruler_bad_dream'
UNION ALL
SELECT id, 'accept_volunteers', 'Принять добровольцев', 'Накормить их и поставить в ополчение.', '{"unitDelta":{"militia":3},"resourceDelta":{"food":-20}}'::jsonb, 'Новые люди в строю', 'Добровольцы получили похлёбку и копья. Строй стал плотнее.'
FROM inserted WHERE event_key = 'volunteers_at_gate'
UNION ALL
SELECT id, 'send_home', 'Отправить домой', 'Пусть лучше принесут пользу в хозяйстве.', '{"resourceDelta":{"food":10}}'::jsonb, 'Добровольцы вернулись к дворам', 'Юноши ушли без обиды. К вечеру их семьи прислали немного припасов.'
FROM inserted WHERE event_key = 'volunteers_at_gate'
UNION ALL
SELECT id, 'receive_politely', 'Принять вежливо', 'Встретить посла хлебом, солью и тяжёлой улыбкой.', '{"patronFavorDelta":2,"resourceDelta":{"gold":-20}}'::jsonb, 'Посол доволен', 'Посол уехал с подарками и правильными словами. Печать на грамоте будто стала легче.'
FROM inserted WHERE event_key = 'patron_envoy'
UNION ALL
SELECT id, 'keep_distance', 'Держать дистанцию', 'Выслушать посла без лишних поклонов.', '{"patronFavorDelta":-1,"kingdomDelta":{"honor":1}}'::jsonb, 'Двор сохранил лицо', 'Посол уехал холодным, но люди заметили, что двор не спешит сгибать шею.'
FROM inserted WHERE event_key = 'patron_envoy'
UNION ALL
SELECT id, 'burn_incense', 'Жечь благовония', 'Потратить припасы на обряд очищения.', '{"resourceDelta":{"food":-30},"kingdomDelta":{"honor":1}}'::jsonb, 'Дым поднялся над стенами', 'Птицы исчезли до полудня. Люди всё ещё шепчутся, но благодарят за обряд.'
FROM inserted WHERE event_key = 'black_birds_over_walls'
UNION ALL
SELECT id, 'ignore_omen', 'Не обращать внимания', 'Пусть птицы сидят, пока не проголодаются.', '{"resourceDelta":{"food":20},"kingdomDelta":{"dread":1}}'::jsonb, 'Знамение оставили без ответа', 'Птицы улетели сами. Люди решили, что двор не боится даже дурных крыльев.'
FROM inserted WHERE event_key = 'black_birds_over_walls';

-- +goose Down
DROP TABLE IF EXISTS kingdom_events;
DROP TABLE IF EXISTS event_choices;
DROP TABLE IF EXISTS game_events;

DELETE FROM mission_reports
WHERE type = 'event';

ALTER TABLE mission_reports
DROP CONSTRAINT mission_reports_type_valid;

ALTER TABLE mission_reports
ADD CONSTRAINT mission_reports_type_valid CHECK (
    type IN (
        'pve_mission',
        'pvp_raid_attacker',
        'pvp_raid_defender'
    )
);
