import type { ButtonHTMLAttributes, ReactNode } from 'react';

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  children: ReactNode;
};

export function Button({ children, className = '', ...props }: ButtonProps) {
  return (
    <button
      className={`max-w-full rounded bg-dusk-gold px-4 py-2 text-center text-sm font-semibold text-dusk-950 transition hover:bg-amber-300 disabled:cursor-not-allowed disabled:bg-stone-700 disabled:text-stone-400 ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}
