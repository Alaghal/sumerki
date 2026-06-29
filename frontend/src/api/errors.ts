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
};

export function toUserMessage(error: unknown): string {
  if (error instanceof ApiError && errorMessages[error.code]) {
    return errorMessages[error.code];
  }

  return 'Что-то пошло не так. Попробуйте ещё раз.';
}
