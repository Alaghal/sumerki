import type { ReactNode } from 'react';

type CardProps = {
  title?: string;
  children: ReactNode;
};

export function Card({ title, children }: CardProps) {
  return (
    <section className="rounded border border-stone-800 bg-dusk-900 p-4 shadow-sm shadow-black/20">
      {title ? <h2 className="mb-3 text-base font-semibold text-stone-100">{title}</h2> : null}
      <div className="text-sm leading-6 text-stone-300">{children}</div>
    </section>
  );
}
