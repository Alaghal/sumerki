import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { Culture } from '../api/client';
import { toUserMessage } from '../api/errors';
import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';
import { useSession } from '../context/SessionContext';

const cultures = [
  {
    value: 'northern_principality' as const,
    label: 'Северные Княжества',
    description: 'крепости, дружины, честь и прямой военный путь.',
  },
  {
    value: 'lizard_grad' as const,
    label: 'Ящерские Грады',
    description: 'болота, древние договоры, ловушки и холодная память.',
  },
  {
    value: 'free_posad' as const,
    label: 'Вольные Посады',
    description: 'торговля, золото, рынки и наёмники.',
  },
];

export function CreateKingdomPage() {
  const navigate = useNavigate();
  const { createUserKingdom } = useSession();
  const [name, setName] = useState('');
  const [culture, setCulture] = useState<Culture>('northern_principality');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const trimmedName = name.trim();
    const nameLength = Array.from(trimmedName).length;

    if (!trimmedName) {
      setError('Введите название владения.');
      return;
    }

    if (nameLength < 3) {
      setError('Название слишком короткое.');
      return;
    }

    if (nameLength > 32) {
      setError('Название слишком длинное.');
      return;
    }

    setIsSubmitting(true);
    setError('');

    try {
      await createUserKingdom(trimmedName, culture);
      navigate('/app', { replace: true });
    } catch (caughtError) {
      setError(toUserMessage(caughtError));
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-2xl items-center px-4 py-10">
      <Card title="Create Kingdom">
        <form className="grid gap-5" onSubmit={handleSubmit}>
          <label className="grid gap-2">
            <span>Kingdom name</span>
            <input
              className="rounded border border-stone-700 bg-dusk-950 px-3 py-2 text-stone-100"
              maxLength={32}
              minLength={3}
              onChange={(event) => setName(event.target.value)}
              required
              type="text"
              value={name}
            />
          </label>
          <label className="grid gap-2">
            <span>Culture</span>
            <select
              className="rounded border border-stone-700 bg-dusk-950 px-3 py-2 text-stone-100"
              onChange={(event) => setCulture(event.target.value as Culture)}
              required
              value={culture}
            >
              {cultures.map((culture) => (
                <option key={culture.value} value={culture.value}>
                  {culture.label}
                </option>
              ))}
            </select>
          </label>
          <div className="grid gap-3 sm:grid-cols-3">
            {cultures.map((culture) => (
              <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={culture.value}>
                <h2 className="text-sm font-semibold text-stone-100">{culture.label}</h2>
                <p className="mt-2 text-sm text-stone-400">{culture.description}</p>
              </div>
            ))}
          </div>
          <Button disabled={isSubmitting} type="submit">
            {isSubmitting ? 'Creating kingdom...' : 'Create kingdom'}
          </Button>
          {error ? <p className="text-red-300">{error}</p> : null}
        </form>
      </Card>
    </main>
  );
}
