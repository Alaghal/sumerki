import { FormEvent, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';

import { toUserMessage } from '../api/errors';
import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';
import { useSession } from '../context/SessionContext';

export function RegisterPage() {
  const navigate = useNavigate();
  const { t } = useTranslation('auth');
  const { registerUser } = useSession();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);
    setError('');

    try {
      await registerUser(email, password);
      navigate('/create-kingdom', { replace: true });
    } catch (caughtError) {
      setError(toUserMessage(caughtError));
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-md items-center px-4 py-10">
      <Card title={t('register.title')}>
        <form className="grid gap-4" onSubmit={handleSubmit}>
          <label className="grid gap-2">
            <span>{t('register.email')}</span>
            <input
              autoComplete="email"
              className="rounded border border-stone-700 bg-dusk-950 px-3 py-2 text-stone-100"
              onChange={(event) => setEmail(event.target.value)}
              required
              type="email"
              value={email}
            />
          </label>
          <label className="grid gap-2">
            <span>{t('register.password')}</span>
            <input
              autoComplete="new-password"
              className="rounded border border-stone-700 bg-dusk-950 px-3 py-2 text-stone-100"
              minLength={8}
              onChange={(event) => setPassword(event.target.value)}
              required
              type="password"
              value={password}
            />
          </label>
          <Button disabled={isSubmitting} type="submit">
            {isSubmitting ? t('register.loading') : t('register.submit')}
          </Button>
          {error ? <p className="text-red-300">{error}</p> : null}
          <Link className="text-dusk-gold hover:text-amber-300" to="/login">
            {t('register.toLogin')}
          </Link>
        </form>
      </Card>
    </main>
  );
}
