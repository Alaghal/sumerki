# UI Copy Rules

## Language Rules

* Default language is Russian.
* English is the second supported language.
* Do not mix languages inside one UI state.
* No hardcoded player-facing strings after UI copy migration.
* Do not show raw enum/API values to players.
* Technical IDs, UUIDs, event keys, mission keys, and unit keys should not appear in normal UI.

## Bad Examples

* `Logout`
* `Kingdom`
* `Culture`
* `Patron`
* `Settlement dashboard for`
* `Favor`
* `northern_principality`
* `empire_of_dusk`
* `black_forest_expedition`
* `active`
* `completed`
* `target_too_weak`

## Good Russian Examples

* `Выйти`
* `Владение`
* `Культура`
* `Покровитель`
* `Двор`
* `Влияние`
* `Северные Княжества`
* `Империя Заката`
* `Чёрный Лес`
* `В пути`
* `Завершено`
* `Цель слишком слаба`

## Tone

* concise
* atmospheric but readable
* serious dark Slavic frontier mood
* no huge lore dumps in buttons
* no modern slang
* no technical jargon
* reports and events may be more atmospheric than controls
* action buttons should stay clear and short

## Buttons

Use action verbs.

Russian examples:

* `Обновить`
* `Отправить`
* `Улучшить`
* `Обучить`
* `Выбрать`
* `Разорвать связь`
* `Заплатить дань`
* `Попросить отсрочку`
* `Открыть отчёт`

English examples:

* `Refresh`
* `Send`
* `Upgrade`
* `Train`
* `Choose`
* `Break tie`
* `Pay tribute`
* `Ask for delay`
* `Open report`

## Empty States

Empty states should explain the next useful action.

Bad:

* `No data`

Good:

* `Отчётов пока нет. Отправьте отряд в первую экспедицию.`
* `Событий пока нет. Обновите двор позже.`
* `Нет активных набегов.`
* `No reports yet. Send a party on its first expedition.`
* `No events yet. Check the court again later.`

## Error Messages

Errors should be:

* short
* human-readable
* non-technical
* not expose stack traces
* not expose raw backend errors

## Date And Time

* Use localized date/time formatting.
* Avoid raw ISO strings in player-facing UI.
* Show timers as relative when possible:
  * `Осталось 2 мин.`
  * `Завершится сегодня в 18:40`

## Numbers

* Use readable grouped numbers if values grow.
* Keep resource values integer in MVP.
* Avoid showing internal calculations unless needed.
