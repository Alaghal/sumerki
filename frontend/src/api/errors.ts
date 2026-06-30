import { ApiError } from './client';

const errorMessages: Record<string, string> = {
  invalid_credentials: 'Неверная почта или пароль.',
  email_already_exists: 'Такой email уже зарегистрирован.',
  password_too_short: 'Пароль должен быть не короче 8 символов.',
  invalid_email: 'Введите корректный email.',
  kingdom_already_exists: 'У вас уже есть владение.',
  invalid_culture: 'Неверная культура.',
  kingdom_name_too_short: 'Название слишком короткое.',
  kingdom_name_too_long: 'Название слишком длинное.',
  invalid_building_type: 'Неверный тип здания.',
  building_not_found: 'Здание не найдено.',
  building_already_upgrading: 'Здание уже улучшается.',
  building_max_level: 'Здание уже максимального уровня.',
  insufficient_resources: 'Недостаточно ресурсов.',
  invalid_unit_type: 'Неверный тип войск.',
  invalid_training_amount: 'Количество должно быть от 1 до 50.',
  barracks_level_too_low: 'Недостаточный уровень казармы.',
  invalid_mission_key: 'Неверный поход.',
  invalid_unit_amount: 'Количество войск должно быть целым и неотрицательным.',
  insufficient_units: 'Недостаточно доступных войск.',
  mission_requirements_not_met: 'Требования похода не выполнены.',
  event_expired: 'Событие истекло.',
  event_already_resolved: 'Событие уже решено.',
  invalid_event_choice: 'Неверный выбор.',
  event_choice_not_available: 'Этот выбор недоступен.',
  event_not_found: 'Событие не найдено.',
};

export function toUserMessage(error: unknown): string {
  if (error instanceof ApiError && errorMessages[error.code]) {
    return errorMessages[error.code];
  }

  return 'Что-то пошло не так. Попробуйте ещё раз.';
}
