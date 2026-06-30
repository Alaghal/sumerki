# UI Copy Rules

## Language

* Default UI language: Russian.
* Secondary UI language: English.
* Do not mix languages in one UI state.
* Do not show raw enum values to players.

## Bad Examples

* `Logout`
* `Kingdom`
* `Culture`
* `Patron`
* `Settlement dashboard for`
* `northern_principality`
* `empire_of_dusk`
* `black_forest_expedition`
* `active`
* `completed`

## Good Russian Examples

* `Выйти`
* `Владение`
* `Культура`
* `Покровитель`
* `Воронья Сечь`
* `Северные Княжества`
* `Империя Заката`
* `Чёрный Лес`
* `В пути`
* `Завершено`

## Tone

* concise
* atmospheric but readable
* no huge lore dumps in UI controls
* no modern slang
* no technical vocabulary unless necessary
* reports and events may be more atmospheric than buttons

## Buttons

Buttons should be action verbs:

* `Обновить`
* `Отправить`
* `Улучшить`
* `Обучить`
* `Выбрать`
* `Разорвать связь`
* `Заплатить дань`

## Empty States

Empty states should explain the next useful action.

Bad:

* `No data`

Good:

* `Отчётов пока нет. Отправьте отряд в первую экспедицию.`
* `Событий пока нет. Обновите двор позже.`
* `Нет активных набегов.`

## Error Messages

Errors should be:

* short
* human-readable
* non-technical
* not expose stack traces or raw backend errors
