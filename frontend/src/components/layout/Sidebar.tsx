import { NavLink } from 'react-router-dom';

const links = [
  'Город',
  'Правитель',
  'Армия',
  'Окрестности',
  'Отчёты',
  'Покровитель',
];

export function Sidebar() {
  return (
    <aside className="rounded border border-stone-800 bg-dusk-900 p-3 md:sticky md:top-4 md:self-start">
      <nav className="grid gap-1">
        {links.map((label) => (
          <NavLink
            className="rounded px-3 py-2 text-sm text-stone-300 hover:bg-dusk-800 hover:text-stone-100"
            key={label}
            to="/app"
          >
            {label}
          </NavLink>
        ))}
      </nav>
    </aside>
  );
}
