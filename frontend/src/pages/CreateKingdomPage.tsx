import { FormEvent, useState } from 'react';

import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';

const cultures = [
  {
    value: 'northern_principality',
    label: 'Северные Княжества',
    description: 'крепости, дружины, честь и прямой военный путь.',
  },
  {
    value: 'lizard_grad',
    label: 'Ящерские Грады',
    description: 'болота, древние договоры, ловушки и холодная память.',
  },
  {
    value: 'free_posad',
    label: 'Вольные Посады',
    description: 'торговля, золото, рынки и наёмники.',
  },
];

export function CreateKingdomPage() {
  const [message, setMessage] = useState('');

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setMessage('Kingdom creation will connect to the backend in a later phase.');
  }

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-2xl items-center px-4 py-10">
      <Card title="Create Kingdom">
        <form className="grid gap-5" onSubmit={handleSubmit}>
          <label className="grid gap-2">
            <span>Kingdom name</span>
            <input className="rounded border border-stone-700 bg-dusk-950 px-3 py-2 text-stone-100" type="text" />
          </label>
          <label className="grid gap-2">
            <span>Culture</span>
            <select className="rounded border border-stone-700 bg-dusk-950 px-3 py-2 text-stone-100">
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
          <Button type="submit">Create kingdom</Button>
          {message ? <p className="text-dusk-gold">{message}</p> : null}
        </form>
      </Card>
    </main>
  );
}
