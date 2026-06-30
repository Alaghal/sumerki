import { useTranslation } from 'react-i18next';
import { NavLink } from 'react-router-dom';

const links = [
  'city',
  'ruler',
  'army',
  'outskirts',
  'reports',
  'patron',
] as const;

export function Sidebar() {
  const { t } = useTranslation('game');

  return (
    <aside className="rounded border border-stone-800 bg-dusk-900 p-3 md:sticky md:top-4 md:self-start">
      <nav className="grid gap-1">
        {links.map((key) => (
          <NavLink
            className="rounded px-3 py-2 text-sm text-stone-300 hover:bg-dusk-800 hover:text-stone-100"
            key={key}
            to="/app"
          >
            {t(`navigation.${key}`)}
          </NavLink>
        ))}
      </nav>
    </aside>
  );
}
