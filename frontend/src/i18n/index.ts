import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import enAuth from './resources/en/auth.json';
import enCommon from './resources/en/common.json';
import enErrors from './resources/en/errors.json';
import enGame from './resources/en/game.json';
import enKingdom from './resources/en/kingdom.json';
import ruAuth from './resources/ru/auth.json';
import ruCommon from './resources/ru/common.json';
import ruErrors from './resources/ru/errors.json';
import ruGame from './resources/ru/game.json';
import ruKingdom from './resources/ru/kingdom.json';

export const languageStorageKey = 'sumerki.ui.language';

export const supportedLanguages = [
  { code: 'ru', label: 'Русский' },
  { code: 'en', label: 'English' },
] as const;

export type SupportedLanguage = (typeof supportedLanguages)[number]['code'];

function isSupportedLanguage(language: string | null): language is SupportedLanguage {
  return supportedLanguages.some((supportedLanguage) => supportedLanguage.code === language);
}

function getInitialLanguage(): SupportedLanguage {
  const savedLanguage = window.localStorage.getItem(languageStorageKey);
  return isSupportedLanguage(savedLanguage) ? savedLanguage : 'ru';
}

void i18n.use(initReactI18next).init({
  fallbackLng: 'en',
  interpolation: {
    escapeValue: false,
  },
  lng: getInitialLanguage(),
  ns: ['common', 'game', 'auth', 'kingdom', 'errors'],
  resources: {
    en: {
      auth: enAuth,
      common: enCommon,
      errors: enErrors,
      game: enGame,
      kingdom: enKingdom,
    },
    ru: {
      auth: ruAuth,
      common: ruCommon,
      errors: ruErrors,
      game: ruGame,
      kingdom: ruKingdom,
    },
  },
  supportedLngs: supportedLanguages.map((language) => language.code),
});

i18n.on('languageChanged', (language) => {
  if (isSupportedLanguage(language)) {
    window.localStorage.setItem(languageStorageKey, language);
  }
});

export default i18n;
