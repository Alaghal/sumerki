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

export function DashboardPage() {
  const { user, kingdom } = useSession();

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
          <Card title="Ruler">A ruler card arrives in a later phase.</Card>
          <Card title="Buildings">Town hall, farms, walls, and other structures are placeholders.</Card>
          <Card title="Army">Militia and scouts are not trained yet.</Card>
          <Card title="Reports">Mission and raid reports will be listed here.</Card>
        </div>
      </div>
    </AppShell>
  );
}
