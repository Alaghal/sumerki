import { ApiError } from './client';
import i18n from '../i18n';

export function toUserMessage(error: unknown): string {
  if (error instanceof ApiError) {
    return i18n.t(`errors:${error.code}`, { defaultValue: i18n.t('errors:generic') });
  }

  return i18n.t('errors:generic');
}
