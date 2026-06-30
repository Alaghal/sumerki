import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import enAuth from './resources/en/auth.json';
import enBuildings from './resources/en/buildings.json';
import enCommon from './resources/en/common.json';
import enErrors from './resources/en/errors.json';
import enEvents from './resources/en/events.json';
import enGame from './resources/en/game.json';
import enKingdom from './resources/en/kingdom.json';
import enMap from './resources/en/map.json';
import enMissions from './resources/en/missions.json';
import enPatrons from './resources/en/patrons.json';
import enRaids from './resources/en/raids.json';
import enReports from './resources/en/reports.json';
import enResources from './resources/en/resources.json';
import enUnits from './resources/en/units.json';
import ruAuth from './resources/ru/auth.json';
import ruBuildings from './resources/ru/buildings.json';
import ruCommon from './resources/ru/common.json';
import ruErrors from './resources/ru/errors.json';
import ruEvents from './resources/ru/events.json';
import ruGame from './resources/ru/game.json';
import ruKingdom from './resources/ru/kingdom.json';
import ruMap from './resources/ru/map.json';
import ruMissions from './resources/ru/missions.json';
import ruPatrons from './resources/ru/patrons.json';
import ruRaids from './resources/ru/raids.json';
import ruReports from './resources/ru/reports.json';
import ruResources from './resources/ru/resources.json';
import ruUnits from './resources/ru/units.json';

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
  ns: [
    'common',
    'game',
    'auth',
    'kingdom',
    'errors',
    'resources',
    'buildings',
    'units',
    'missions',
    'reports',
    'patrons',
    'events',
    'raids',
    'map',
  ],
  resources: {
    en: {
      auth: enAuth,
      buildings: enBuildings,
      common: enCommon,
      errors: enErrors,
      events: enEvents,
      game: enGame,
      kingdom: enKingdom,
      map: enMap,
      missions: enMissions,
      patrons: enPatrons,
      raids: enRaids,
      reports: enReports,
      resources: enResources,
      units: enUnits,
    },
    ru: {
      auth: ruAuth,
      buildings: ruBuildings,
      common: ruCommon,
      errors: ruErrors,
      events: ruEvents,
      game: ruGame,
      kingdom: ruKingdom,
      map: ruMap,
      missions: ruMissions,
      patrons: ruPatrons,
      raids: ruRaids,
      reports: ruReports,
      resources: ruResources,
      units: ruUnits,
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
