import { useEffect, useState } from 'react';

import { getMyRuler, Ruler } from '../api/client';
import { AppShell } from '../components/layout/AppShell';
import { Card } from '../components/ui/Card';
import { useSession } from '../context/SessionContext';

const cultureLabels = {
  northern_principality: 'Северные Княжества',
  lizard_grad: 'Ящерские Грады',
  free_posad: 'Вольные Посады',
};

const patronLabels = {
  independent: 'Independent',
  empire_of_dusk: 'Empire of Dusk',
  old_pact: 'Old Pact',
};

const healthLabels = {
  healthy: 'Здоров',
  wounded: 'Ранен',
  sick: 'Болен',
};

const rulerStats = [
  ['Власть', 'authority'],
  ['Храбрость', 'courage'],
  ['Хитрость', 'cunning'],
  ['Честь', 'honor'],
  ['Жестокость', 'cruelty'],
  ['Амбиция', 'ambition'],
  ['Паранойя', 'paranoia'],
] as const;

export function DashboardPage() {
  const { token, user, kingdom } = useSession();
  const [ruler, setRuler] = useState<Ruler | null>(null);
  const [rulerLoading, setRulerLoading] = useState(true);
  const [rulerError, setRulerError] = useState('');

  useEffect(() => {
    let isActive = true;

    async function loadRuler() {
      if (!token || !kingdom) {
        return;
      }

      setRulerLoading(true);
      setRulerError('');

      try {
        const response = await getMyRuler(token);
        if (isActive) {
          setRuler(response.ruler);
        }
      } catch {
        if (isActive) {
          setRulerError('Не удалось загрузить правителя.');
        }
      } finally {
        if (isActive) {
          setRulerLoading(false);
        }
      }
    }

    loadRuler();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  if (!user || !kingdom) {
    return null;
  }

  return (
    <AppShell>
      <div className="grid gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-stone-100">{kingdom.name}</h1>
          <p className="mt-1 text-sm text-stone-400">Settlement dashboard for {user.email}</p>
        </div>
        <div className="grid gap-4 lg:grid-cols-2">
          <Card title="Kingdom">
            <dl className="grid gap-2">
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">Culture</dt>
                <dd className="text-right text-stone-100">{cultureLabels[kingdom.culture]}</dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">Patron</dt>
                <dd className="text-right text-stone-100">
                  {kingdom.patron ? patronLabels[kingdom.patron] : 'Без покровителя'}
                </dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">Player</dt>
                <dd className="max-w-[14rem] truncate text-right text-stone-100">{user.email}</dd>
              </div>
            </dl>
          </Card>
          <Card title="Resources">
            Gold, food, wood, stone, and population will appear here when the resources system is implemented.
          </Card>
          <Card title="Ruler">
            {rulerLoading ? <p>Загрузка правителя...</p> : null}
            {rulerError ? <p className="text-red-300">{rulerError}</p> : null}
            {ruler && !rulerLoading && !rulerError ? (
              <div className="grid gap-4">
                <dl className="grid gap-2">
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">Name</dt>
                    <dd className="text-right text-stone-100">{ruler.name}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">Age</dt>
                    <dd className="text-right text-stone-100">{ruler.age}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">Health</dt>
                    <dd className="text-right text-stone-100">{healthLabels[ruler.healthStatus]}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">Culture</dt>
                    <dd className="text-right text-stone-100">{cultureLabels[ruler.culture]}</dd>
                  </div>
                </dl>
                <dl className="grid gap-2">
                  {rulerStats.map(([label, key]) => (
                    <div className="flex justify-between gap-4" key={key}>
                      <dt className="text-stone-400">{label}</dt>
                      <dd className="text-right text-stone-100">{ruler[key]}</dd>
                    </div>
                  ))}
                </dl>
              </div>
            ) : null}
          </Card>
          <Card title="Buildings">Town hall, farms, walls, and other structures are placeholders.</Card>
          <Card title="Army">Militia and scouts are not trained yet.</Card>
          <Card title="Reports">Mission and raid reports will be listed here.</Card>
        </div>
      </div>
    </AppShell>
  );
}
