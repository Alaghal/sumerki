import { useNavigate } from 'react-router-dom';

import { useSession } from '../../context/SessionContext';
import { Button } from '../ui/Button';

export function TopBar() {
  const navigate = useNavigate();
  const { user, logout } = useSession();

  function handleLogout() {
    logout();
    navigate('/login', { replace: true });
  }

  return (
    <header className="border-b border-stone-800 bg-dusk-900/90">
      <div className="mx-auto flex min-h-16 w-full max-w-7xl flex-wrap items-center justify-between gap-3 px-4 py-3">
        <div>
          <div className="text-xl font-semibold tracking-normal text-stone-100">Sumerki</div>
          <div className="text-sm text-dusk-gold">Сумеречье</div>
        </div>
        <div className="flex items-center gap-3 text-sm text-stone-300">
          <span className="max-w-[13rem] truncate">{user?.email}</span>
          <Button onClick={handleLogout} type="button">
            Logout
          </Button>
        </div>
      </div>
    </header>
  );
}
