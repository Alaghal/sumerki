import { FormEvent, useState } from 'react';
import { Link } from 'react-router-dom';

import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';

export function RegisterPage() {
  const [message, setMessage] = useState('');

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setMessage('Registration will connect to the backend in a later phase.');
  }

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-md items-center px-4 py-10">
      <Card title="Register">
        <form className="grid gap-4" onSubmit={handleSubmit}>
          <label className="grid gap-2">
            <span>Email</span>
            <input className="rounded border border-stone-700 bg-dusk-950 px-3 py-2 text-stone-100" type="email" />
          </label>
          <label className="grid gap-2">
            <span>Password</span>
            <input className="rounded border border-stone-700 bg-dusk-950 px-3 py-2 text-stone-100" type="password" />
          </label>
          <Button type="submit">Register</Button>
          {message ? <p className="text-dusk-gold">{message}</p> : null}
          <Link className="text-dusk-gold hover:text-amber-300" to="/login">
            Already have an account?
          </Link>
        </form>
      </Card>
    </main>
  );
}
