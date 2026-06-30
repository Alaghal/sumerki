import { useTranslation } from 'react-i18next';

import { supportedLanguages, SupportedLanguage } from '../../i18n';

export function LanguageSwitcher() {
  const { i18n, t } = useTranslation('common');

  function handleLanguageChange(language: SupportedLanguage) {
    void i18n.changeLanguage(language);
  }

  return (
    <div aria-label={t('language.label')} className="flex items-center gap-1 rounded border border-stone-700 p-1">
      {supportedLanguages.map((language) => {
        const isActive = i18n.language === language.code;
        return (
          <button
            aria-pressed={isActive}
            className={`rounded px-2 py-1 text-xs font-semibold transition ${
              isActive ? 'bg-dusk-gold text-dusk-950' : 'text-stone-300 hover:bg-dusk-800 hover:text-stone-100'
            }`}
            key={language.code}
            onClick={() => handleLanguageChange(language.code)}
            type="button"
          >
            {t(`language.${language.code}`)}
          </button>
        );
      })}
    </div>
  );
}
