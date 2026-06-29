import type { ReactNode } from 'react';

import { Navigate, Route, Routes } from 'react-router-dom';

import { useSession } from '../context/SessionContext';
import { CreateKingdomPage } from '../pages/CreateKingdomPage';
import { DashboardPage } from '../pages/DashboardPage';
import { LoginPage } from '../pages/LoginPage';
import { NotFoundPage } from '../pages/NotFoundPage';
import { RegisterPage } from '../pages/RegisterPage';

function SessionLoading() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-dusk-950 px-4 text-stone-100">
      <div className="rounded border border-stone-800 bg-dusk-900 px-5 py-4 text-sm text-stone-300">
        Загрузка...
      </div>
    </main>
  );
}

type RouteGuardProps = {
  children: ReactNode;
};

function PublicRoute({ children }: RouteGuardProps) {
  const { token, user, kingdom, loading } = useSession();

  if (loading || (token && !user)) {
    return <SessionLoading />;
  }

  if (token && user && kingdom) {
    return <Navigate to="/app" replace />;
  }

  if (token && user && !kingdom) {
    return <Navigate to="/create-kingdom" replace />;
  }

  return children;
}

function RequireAuth({ children }: RouteGuardProps) {
  const { token, user, loading } = useSession();

  if (loading || (token && !user)) {
    return <SessionLoading />;
  }

  if (!token || !user) {
    return <Navigate to="/login" replace />;
  }

  return children;
}

function RequireKingdom({ children }: RouteGuardProps) {
  const { kingdom } = useSession();

  if (!kingdom) {
    return <Navigate to="/create-kingdom" replace />;
  }

  return children;
}

function CreateKingdomRoute() {
  const { kingdom } = useSession();

  return (
    <RequireAuth>
      {kingdom ? <Navigate to="/app" replace /> : <CreateKingdomPage />}
    </RequireAuth>
  );
}

export function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/login" replace />} />
      <Route
        path="/login"
        element={
          <PublicRoute>
            <LoginPage />
          </PublicRoute>
        }
      />
      <Route
        path="/register"
        element={
          <PublicRoute>
            <RegisterPage />
          </PublicRoute>
        }
      />
      <Route path="/create-kingdom" element={<CreateKingdomRoute />} />
      <Route
        path="/app"
        element={
          <RequireAuth>
            <RequireKingdom>
              <DashboardPage />
            </RequireKingdom>
          </RequireAuth>
        }
      />
      <Route path="/404" element={<NotFoundPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
