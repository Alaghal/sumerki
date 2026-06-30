import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

import { useSession } from '../../context/SessionContext';
import { LanguageSwitcher } from '../i18n/LanguageSwitcher';
import { Button } from '../ui/Button';

export function TopBar() {
  const navigate = useNavigate();
  const { t } = useTranslation(['common', 'game']);
  const { user, logout } = useSession();

  function handleLogout() {
    logout();
    navigate('/login', { replace: true });
  }

  return (
    <header className="border-b border-stone-800 bg-dusk-900/90">
      <div className="mx-auto flex min-h-16 w-full max-w-7xl flex-wrap items-center justify-between gap-3 px-4 py-3">
        <div className="min-w-0 max-w-full">
          <div className="flex flex-wrap items-center gap-2">
            <div className="break-words text-xl font-semibold tracking-normal text-stone-100">{t('game:app.name')}</div>
            <span className="rounded border border-dusk-gold/40 px-2 py-0.5 text-xs uppercase tracking-normal text-dusk-gold">
              {t('game:app.playtestLabel')}
            </span>
          </div>
          <div className="break-words text-sm text-dusk-gold">{t('game:app.subtitle')}</div>
        </div>
        <div className="flex min-w-0 max-w-full flex-wrap items-center justify-start gap-3 text-sm text-stone-300 sm:justify-end">
          <LanguageSwitcher />
          <span className="min-w-0 max-w-full break-words sm:max-w-[13rem]">{user?.email}</span>
          <Button onClick={handleLogout} type="button">
            {t('common:buttons.logout')}
          </Button>
        </div>
      </div>
    </header>
  );
}
