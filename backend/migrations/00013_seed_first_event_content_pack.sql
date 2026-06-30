-- +goose Up
WITH upserted AS (
    INSERT INTO game_events (event_key, category, title, body, conditions_json)
    VALUES
        ('cracked_granary_roof', 'economy', 'Трещина в амбаре', 'После ночного дождя амбарная крыша дала течь. Если ждать, зерно начнёт портиться.', '{}'::jsonb),
        ('lucky_market_day', 'economy', 'Удачный торг', 'На рынок пришли чужие купцы. Их кошели тяжелы, а вопросы осторожны.', '{}'::jsonb),
        ('stone_in_old_well', 'economy', 'Камень в старом колодце', 'Рабочие чистили старый колодец и нашли каменную кладку. Под ней может быть древний ход или просто хороший камень.', '{}'::jsonb),
        ('forest_tax_dispute', 'economy', 'Спор о лесной дани', 'Лесорубы спорят с казначеем. Они говорят, что новая дань оставит их семьи без хлеба.', '{}'::jsonb),
        ('ruler_public_judgment', 'ruler', 'Суд у ворот', 'Два рода спорят за землю у ручья. Оба требуют, чтобы правитель решил дело лично.', '{}'::jsonb),
        ('ruler_fever', 'ruler', 'Жар правителя', 'Правитель проснулся в жару. Лекари спорят: дать покой или показать людям, что двор не слаб.', '{}'::jsonb),
        ('court_whispers', 'ruler', 'Шёпот при дворе', 'Советники говорят, что кто-то слишком часто встречается с чужими купцами.', '{}'::jsonb),
        ('ruler_generous_feast', 'ruler', 'Пир для двора', 'После удачных дней дружина просит пира. Казначей шипит, что зима не за горами.', '{}'::jsonb),
        ('rusty_spears_found', 'military', 'Ржавые копья', 'В старом сарае нашли связку копий. Железо плохое, но в руках ополчения и такое сгодится.', '{}'::jsonb),
        ('young_hunters_offer', 'military', 'Молодые охотники', 'Несколько охотников предлагают служить разведчиками. Они знают лес, но плохо слушают приказы.', '{}'::jsonb),
        ('militia_drill', 'military', 'Учение ополчения', 'Воевода просит день на учение. Работы встанут, зато люди научатся держать строй.', '{}'::jsonb),
        ('deserter_at_night', 'military', 'Дезертир ночью', 'Ночью стража поймала человека из ополчения. Он собирался уйти к родне за рекой.', '{}'::jsonb),
        ('imperial_tax_scribe', 'patron', 'Имперский писарь', 'Писарь Заката просит сверить списки дворов. Он улыбается так, будто уже знает ответ.', '{"requiresPatron":"empire_of_dusk"}'::jsonb),
        ('imperial_road_offer', 'patron', 'Чертёж дороги', 'Имперский мастер предлагает улучшить дорогу к рынку. За работу он просит камень и покорность срокам.', '{"requiresPatron":"empire_of_dusk"}'::jsonb),
        ('old_pact_oath_stone', 'patron', 'Камень клятвы', 'Волхвы Старого Договора просят обновить клятву у старого камня. Люди ждут знака.', '{"requiresPatron":"old_pact"}'::jsonb),
        ('old_pact_refugees', 'patron', 'Люди с Рубежа', 'С Рубежа пришли беженцы. Старый Договор просит принять хотя бы часть семей.', '{"requiresPatron":"old_pact"}'::jsonb),
        ('black_milk_morning', 'dark_omen', 'Чёрное молоко', 'Утром одна семья принесла ведро молока цвета сажи. Корова здорова, но никто не хочет пить это.', '{}'::jsonb),
        ('whispering_grain', 'dark_omen', 'Шёпот в зерне', 'На складе слышат шёпот, когда пересыпают зерно. Слова никто не может повторить.', '{}'::jsonb),
        ('red_moon_over_walls', 'dark_omen', 'Красная луна', 'Луна поднялась красной, как уголь под золой. Стража просит двойной караул.', '{}'::jsonb),
        ('child_names_shadow', 'dark_omen', 'Тень зовёт имя', 'Ребёнок у колодца сказал, что его имя позвала тень. Взрослые спорят, был ли это ветер.', '{}'::jsonb)
    ON CONFLICT (event_key) DO UPDATE
    SET category = EXCLUDED.category,
        title = EXCLUDED.title,
        body = EXCLUDED.body,
        conditions_json = EXCLUDED.conditions_json,
        is_active = true,
        updated_at = now()
    RETURNING id, event_key
)
INSERT INTO event_choices (game_event_id, choice_key, label, description, effects_json, result_title, result_body)
SELECT id, 'repair_now', 'Починить сразу', 'Потратить дерево и сохранить припасы.', '{"resourceDelta":{"wood":-40,"food":40},"kingdomDelta":{"honor":1}}'::jsonb, 'Амбар спасён', 'Крыша держит. Люди ворчат на работу, но зерно осталось сухим.'
FROM upserted WHERE event_key = 'cracked_granary_roof'
UNION ALL
SELECT id, 'ignore_leak', 'Подождать', 'Не тратить дерево, но рискнуть припасами.', '{"resourceDelta":{"food":-40,"wood":20},"kingdomDelta":{"honor":-1}}'::jsonb, 'Сырая кладовая', 'Доски остались на складе, но часть зерна почернела от сырости.'
FROM upserted WHERE event_key = 'cracked_granary_roof'
UNION ALL
SELECT id, 'sell_surplus', 'Продать излишки', 'Обменять часть еды и дерева на золото.', '{"resourceDelta":{"gold":120,"food":-40,"wood":-30}}'::jsonb, 'Звонкая сделка', 'Купцы ушли довольными. Казна стала тяжелее, склады чуть пустее.'
FROM upserted WHERE event_key = 'lucky_market_day'
UNION ALL
SELECT id, 'refuse_trade', 'Беречь запасы', 'Отказаться от сделки и сохранить припасы.', '{"resourceDelta":{"food":20,"wood":20},"kingdomDelta":{"honor":1}}'::jsonb, 'Закрытые ворота', 'Купцы ушли к соседям. Зато в городе спокойнее смотрят на зиму.'
FROM upserted WHERE event_key = 'lucky_market_day'
UNION ALL
SELECT id, 'dismantle_stone', 'Разобрать кладку', 'Получить камень для стройки.', '{"resourceDelta":{"stone":90,"food":-10}}'::jsonb, 'Камень поднят', 'Кладку разобрали без находок. Зато телеги привезли крепкий камень.'
FROM upserted WHERE event_key = 'stone_in_old_well'
UNION ALL
SELECT id, 'seal_well', 'Запечатать колодец', 'Не тревожить старое место.', '{"resourceDelta":{"stone":-10},"kingdomDelta":{"honor":1}}'::jsonb, 'Колодец закрыт', 'Старики одобрили решение. Никто не знает, что было под кладкой, и это всех устраивает.'
FROM upserted WHERE event_key = 'stone_in_old_well'
UNION ALL
SELECT id, 'lower_tax', 'Снизить дань', 'Потерять золото, но сохранить спокойствие.', '{"resourceDelta":{"gold":-50,"wood":40},"kingdomDelta":{"honor":1}}'::jsonb, 'Топоры снова стучат', 'Лесорубы вернулись к работе. Казначей недоволен, но народ запомнил мягкость двора.'
FROM upserted WHERE event_key = 'forest_tax_dispute'
UNION ALL
SELECT id, 'enforce_tax', 'Взять своё', 'Получить золото ценой недовольства.', '{"resourceDelta":{"gold":80,"wood":-30},"kingdomDelta":{"dread":1,"honor":-1}}'::jsonb, 'Дань собрана', 'Казна пополнилась. В лесу теперь работают молча и смотрят в землю.'
FROM upserted WHERE event_key = 'forest_tax_dispute'
UNION ALL
SELECT id, 'judge_by_oath', 'Судить по клятве', 'Укрепить честь, но обидеть одну сторону.', '{"kingdomDelta":{"honor":2},"resourceDelta":{"gold":-20}}'::jsonb, 'Клятва решает', 'Решение приняли не все, но никто не посмел спорить с произнесённой клятвой.'
FROM upserted WHERE event_key = 'ruler_public_judgment'
UNION ALL
SELECT id, 'judge_by_payment', 'Взять выкуп', 'Пополнить казну и закрыть спор силой власти.', '{"resourceDelta":{"gold":70},"kingdomDelta":{"dread":1,"honor":-1}}'::jsonb, 'Спор куплен', 'Землю поделили по грамоте, а казна стала тяжелее. У ручья теперь говорят тише.'
FROM upserted WHERE event_key = 'ruler_public_judgment'
UNION ALL
SELECT id, 'let_rest', 'Дать покой', 'Потратить еду и травы, сохранив лицо двора.', '{"resourceDelta":{"food":-40},"kingdomDelta":{"honor":1}}'::jsonb, 'Тихая комната', 'Двор притих на день. К вечеру жар спал, но слухи уже прошли по рынку.'
FROM upserted WHERE event_key = 'ruler_fever'
UNION ALL
SELECT id, 'appear_publicly', 'Выйти к людям', 'Показать силу, но рискнуть здоровьем.', '{"kingdomDelta":{"dread":1,"honor":1},"resourceDelta":{"food":-20}}'::jsonb, 'Бледный выход', 'Люди увидели правителя у ворот. Он стоял прямо, хотя руки дрожали под плащом.'
FROM upserted WHERE event_key = 'ruler_fever'
UNION ALL
SELECT id, 'investigate_quietly', 'Проверить тихо', 'Потратить золото на тайную проверку.', '{"resourceDelta":{"gold":-40},"kingdomDelta":{"honor":1}}'::jsonb, 'Тихая проверка', 'Следы ведут к мелкому писарю. Двору не пришлось устраивать показную бурю.'
FROM upserted WHERE event_key = 'court_whispers'
UNION ALL
SELECT id, 'make_example', 'Наказать показательно', 'Укрепить страх, но потерять доверие.', '{"kingdomDelta":{"dread":2,"honor":-1}}'::jsonb, 'Громкое наказание', 'Люди поняли намёк. Теперь при дворе говорят меньше, но слушают внимательнее.'
FROM upserted WHERE event_key = 'court_whispers'
UNION ALL
SELECT id, 'hold_feast', 'Устроить пир', 'Потратить еду и золото ради верности.', '{"resourceDelta":{"food":-80,"gold":-40},"kingdomDelta":{"honor":2}}'::jsonb, 'Тёплый зал', 'Вечером в зале пели громче ветра. Наутро казна похудела, но лица стали мягче.'
FROM upserted WHERE event_key = 'ruler_generous_feast'
UNION ALL
SELECT id, 'refuse_feast', 'Отказать', 'Сохранить ресурсы и показать строгость.', '{"resourceDelta":{"food":40,"gold":20},"kingdomDelta":{"dread":1,"honor":-1}}'::jsonb, 'Пустые столы', 'Дружина разошлась молча. Никто не спорил, но песни этой ночью не было.'
FROM upserted WHERE event_key = 'ruler_generous_feast'
UNION ALL
SELECT id, 'repair_spears', 'Починить копья', 'Потратить дерево и получить копейщиков.', '{"resourceDelta":{"wood":-40,"gold":-20},"unitDelta":{"spearmen":3}}'::jsonb, 'Копья выпрямлены', 'Кузнец ругался на ржавчину, но к вечеру три новых копья стояли у стены.'
FROM upserted WHERE event_key = 'rusty_spears_found'
UNION ALL
SELECT id, 'melt_parts', 'Разобрать на лом', 'Получить немного золота и камня через продажу.', '{"resourceDelta":{"gold":50,"stone":20}}'::jsonb, 'Старое продано', 'Копья ушли купцам как лом. Войско не выросло, зато казна звякнула.'
FROM upserted WHERE event_key = 'rusty_spears_found'
UNION ALL
SELECT id, 'accept_hunters', 'Принять охотников', 'Получить разведчиков, потратив еду.', '{"resourceDelta":{"food":-30},"unitDelta":{"scouts":2}}'::jsonb, 'Следы приняты', 'Охотники поклялись служить. Их глаза всё ещё чаще смотрят на лес, чем на ворота.'
FROM upserted WHERE event_key = 'young_hunters_offer'
UNION ALL
SELECT id, 'refuse_hunters', 'Отказать', 'Сохранить припасы.', '{"resourceDelta":{"food":20},"kingdomDelta":{"honor":-1}}'::jsonb, 'Лес забрал своих', 'Охотники ушли без поклона. Наутро в лесу стало больше чужих следов.'
FROM upserted WHERE event_key = 'young_hunters_offer'
UNION ALL
SELECT id, 'hold_drill', 'Провести учение', 'Потратить еду, получить ополчение и честь.', '{"resourceDelta":{"food":-50},"unitDelta":{"militia":4},"kingdomDelta":{"honor":1}}'::jsonb, 'Строй у ворот', 'К вечеру руки устали, но копья уже не смотрели в разные стороны.'
FROM upserted WHERE event_key = 'militia_drill'
UNION ALL
SELECT id, 'skip_drill', 'Не отвлекать', 'Сохранить еду и работу.', '{"resourceDelta":{"food":40,"wood":20}}'::jsonb, 'Работа важнее', 'Воевода недоволен. Поля и лесопилки зато не пустовали.'
FROM upserted WHERE event_key = 'militia_drill'
UNION ALL
SELECT id, 'forgive_deserter', 'Простить', 'Сохранить человека, но ослабить страх.', '{"unitDelta":{"militia":1},"kingdomDelta":{"honor":1,"dread":-1}}'::jsonb, 'Прощённый страх', 'Его вернули в строй. Некоторые считают это милостью, другие слабостью.'
FROM upserted WHERE event_key = 'deserter_at_night'
UNION ALL
SELECT id, 'punish_deserter', 'Наказать', 'Укрепить страх, потеряв бойца.', '{"unitDelta":{"militia":-1},"kingdomDelta":{"dread":2,"honor":-1}}'::jsonb, 'Ночная кара', 'После наказания никто не говорил у костров. Утром строй был ровнее обычного.'
FROM upserted WHERE event_key = 'deserter_at_night'
UNION ALL
SELECT id, 'cooperate', 'Дать списки', 'Улучшить отношения с Империей ценой доверия людей.', '{"patronFavorDelta":3,"kingdomDelta":{"honor":-1},"resourceDelta":{"gold":30}}'::jsonb, 'Списки приняты', 'Писарь ушёл довольным. Люди заметили, как долго его перо скребло по бумаге.'
FROM upserted WHERE event_key = 'imperial_tax_scribe'
UNION ALL
SELECT id, 'delay_scribe', 'Тянуть время', 'Сохранить лицо перед людьми, ухудшив отношение Империи.', '{"patronFavorDelta":-2,"kingdomDelta":{"honor":1}}'::jsonb, 'Чернила не высохли', 'Писарь уехал с пустой папкой. Он не угрожал, но его молчание было хуже угроз.'
FROM upserted WHERE event_key = 'imperial_tax_scribe'
UNION ALL
SELECT id, 'accept_road_plan', 'Принять план', 'Потратить камень, получить золото и favor.', '{"resourceDelta":{"stone":-50,"gold":90},"patronFavorDelta":2}'::jsonb, 'Ровная дорога', 'Дорога стала лучше. По ней легче идут телеги и чужие сапоги.'
FROM upserted WHERE event_key = 'imperial_road_offer'
UNION ALL
SELECT id, 'refuse_road_plan', 'Отказаться', 'Сохранить камень, но охладить Империю.', '{"resourceDelta":{"stone":20},"patronFavorDelta":-2,"kingdomDelta":{"honor":1}}'::jsonb, 'Кривая тропа', 'Мастер свернул чертёж. Старые дороги остались старыми, зато своими.'
FROM upserted WHERE event_key = 'imperial_road_offer'
UNION ALL
SELECT id, 'renew_oath', 'Обновить клятву', 'Потратить еду и укрепить favor.', '{"resourceDelta":{"food":-40},"patronFavorDelta":3,"kingdomDelta":{"honor":1}}'::jsonb, 'Клятва звучит', 'Слова легли на камень тяжело и ровно. Даже ветер на миг стих.'
FROM upserted WHERE event_key = 'old_pact_oath_stone'
UNION ALL
SELECT id, 'avoid_rite', 'Уклониться', 'Сохранить припасы, но ослабить доверие Договора.', '{"resourceDelta":{"food":20},"patronFavorDelta":-2}'::jsonb, 'Камень молчит', 'Обряд перенесли. Волхвы ничего не сказали, но их молчание заметили все.'
FROM upserted WHERE event_key = 'old_pact_oath_stone'
UNION ALL
SELECT id, 'accept_families', 'Принять семьи', 'Потратить еду, получить население и favor.', '{"resourceDelta":{"food":-70,"population":4},"patronFavorDelta":2,"kingdomDelta":{"honor":1}}'::jsonb, 'Новые очаги', 'В пустых избах снова дым. Дети боятся громких звуков, но уже помогают носить воду.'
FROM upserted WHERE event_key = 'old_pact_refugees'
UNION ALL
SELECT id, 'refuse_families', 'Закрыть ворота', 'Сохранить еду, потеряв доверие.', '{"resourceDelta":{"food":40},"patronFavorDelta":-3,"kingdomDelta":{"dread":1,"honor":-1}}'::jsonb, 'Ворота закрыты', 'Люди ушли дальше по дороге. Их следы ещё долго были видны в грязи.'
FROM upserted WHERE event_key = 'old_pact_refugees'
UNION ALL
SELECT id, 'burn_the_milk', 'Сжечь молоко', 'Потратить еду, успокоить людей.', '{"resourceDelta":{"food":-30},"kingdomDelta":{"honor":1}}'::jsonb, 'Чёрный дым', 'Молоко сожгли за стеной. Запах был странным, но люди разошлись спокойнее.'
FROM upserted WHERE event_key = 'black_milk_morning'
UNION ALL
SELECT id, 'sell_to_stranger', 'Продать страннику', 'Получить золото и дурной шёпот.', '{"resourceDelta":{"gold":60},"kingdomDelta":{"dread":1,"honor":-1}}'::jsonb, 'Странная сделка', 'Странник заплатил без торга. После него на снегу не осталось следов.'
FROM upserted WHERE event_key = 'black_milk_morning'
UNION ALL
SELECT id, 'bless_granary', 'Очистить амбар', 'Потратить еду и золото, снизить страх.', '{"resourceDelta":{"food":-30,"gold":-20},"kingdomDelta":{"honor":1}}'::jsonb, 'Тихое зерно', 'После обряда амбар замолчал. Крысы тоже пропали, что тревожит кладовщика.'
FROM upserted WHERE event_key = 'whispering_grain'
UNION ALL
SELECT id, 'ignore_whispers', 'Не слушать', 'Сохранить ресурсы и принять тревогу.', '{"resourceDelta":{"food":40},"kingdomDelta":{"dread":1}}'::jsonb, 'Шёпот остался', 'Зерно всё ещё шуршит слишком похоже на слова. Но хлеб из него получается сытный.'
FROM upserted WHERE event_key = 'whispering_grain'
UNION ALL
SELECT id, 'double_watch', 'Удвоить караул', 'Потратить еду, получить ополчение и спокойствие.', '{"resourceDelta":{"food":-40},"unitDelta":{"militia":2},"kingdomDelta":{"honor":1}}'::jsonb, 'Ночь под стражей', 'Ночь прошла без нападения. На рассвете все сделали вид, что не ждали худшего.'
FROM upserted WHERE event_key = 'red_moon_over_walls'
UNION ALL
SELECT id, 'let_them_sleep', 'Дать спать', 'Сохранить еду, но усилить дурную славу.', '{"resourceDelta":{"food":30},"kingdomDelta":{"dread":1}}'::jsonb, 'Тихая ночь', 'Никто не умер. Но утром на воротах нашли красную пыль.'
FROM upserted WHERE event_key = 'red_moon_over_walls'
UNION ALL
SELECT id, 'close_well_day', 'Закрыть колодец', 'Потерять немного еды, успокоить людей.', '{"resourceDelta":{"food":-20},"kingdomDelta":{"honor":1}}'::jsonb, 'Вода под замком', 'Колодец закрыли до утра. Ночью цепь на крышке звякнула сама.'
FROM upserted WHERE event_key = 'child_names_shadow'
UNION ALL
SELECT id, 'mock_fear', 'Высмеять страх', 'Укрепить жёсткость, но потерять честь.', '{"kingdomDelta":{"dread":1,"honor":-1}}'::jsonb, 'Смех у колодца', 'Люди засмеялись слишком громко. Ребёнок больше не подходит к воде один.'
FROM upserted WHERE event_key = 'child_names_shadow'
ON CONFLICT (game_event_id, choice_key) DO UPDATE
SET label = EXCLUDED.label,
    description = EXCLUDED.description,
    effects_json = EXCLUDED.effects_json,
    result_title = EXCLUDED.result_title,
    result_body = EXCLUDED.result_body,
    updated_at = now();

-- +goose Down
DELETE FROM game_events
WHERE event_key IN (
    'cracked_granary_roof',
    'lucky_market_day',
    'stone_in_old_well',
    'forest_tax_dispute',
    'ruler_public_judgment',
    'ruler_fever',
    'court_whispers',
    'ruler_generous_feast',
    'rusty_spears_found',
    'young_hunters_offer',
    'militia_drill',
    'deserter_at_night',
    'imperial_tax_scribe',
    'imperial_road_offer',
    'old_pact_oath_stone',
    'old_pact_refugees',
    'black_milk_morning',
    'whispering_grain',
    'red_moon_over_walls',
    'child_names_shadow'
);
