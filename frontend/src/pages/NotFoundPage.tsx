import { Link } from 'react-router-dom';

import { Card } from '../components/ui/Card';

export function NotFoundPage() {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-md items-center px-4 py-10">
      <Card title="Not found">
        <p>This route does not exist yet.</p>
        <Link className="mt-4 inline-block text-dusk-gold hover:text-amber-300" to="/login">
          Return to login
        </Link>
      </Card>
    </main>
  );
}
