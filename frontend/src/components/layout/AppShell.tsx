import type { ReactNode } from 'react';

import { Sidebar } from './Sidebar';
import { TopBar } from './TopBar';

type AppShellProps = {
  children: ReactNode;
};

export function AppShell({ children }: AppShellProps) {
  return (
    <div className="min-h-screen bg-dusk-950 text-stone-100">
      <TopBar />
      <div className="mx-auto grid w-full max-w-7xl grid-cols-1 gap-0 px-4 pb-8 pt-4 md:grid-cols-[220px_1fr] md:gap-6">
        <Sidebar />
        <main className="min-w-0 py-4 md:py-0">{children}</main>
      </div>
    </div>
  );
}
