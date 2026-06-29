import { AppShell } from '../components/layout/AppShell';
import { Card } from '../components/ui/Card';

const sections = [
  { title: 'Kingdom', body: 'Воронья Сечь waits for its first real data.' },
  { title: 'Resources', body: 'Gold, food, wood, stone, and population will appear here later.' },
  { title: 'Ruler', body: 'A ruler card arrives in a later phase.' },
  { title: 'Buildings', body: 'Town hall, farms, walls, and other structures are placeholders.' },
  { title: 'Army', body: 'Militia and scouts are not trained yet.' },
  { title: 'Reports', body: 'Mission and raid reports will be listed here.' },
];

export function DashboardPage() {
  return (
    <AppShell>
      <div className="grid gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-stone-100">Город</h1>
          <p className="mt-1 text-sm text-stone-400">Placeholder dashboard shell for future game systems.</p>
        </div>
        <div className="grid gap-4 lg:grid-cols-2">
          {sections.map((section) => (
            <Card key={section.title} title={section.title}>
              {section.body}
            </Card>
          ))}
        </div>
      </div>
    </AppShell>
  );
}
