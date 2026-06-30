import { useTranslation } from 'react-i18next';

const legendItems = [
  ['settlement', 'bg-dusk-gold'],
  ['mission', 'bg-emerald-300'],
  ['neighbor', 'bg-red-300'],
  ['patron', 'bg-violet-300'],
  ['event', 'bg-sky-300'],
  ['active', 'bg-dusk-gold'],
  ['warning', 'bg-red-400'],
] as const;

export function MapLegend() {
  const { t } = useTranslation('map');

  return (
    <div className="flex flex-wrap gap-2 text-xs text-stone-400">
      {legendItems.map(([key, colorClass]) => (
        <span className="inline-flex items-center gap-2 rounded border border-stone-800 bg-dusk-950/90 px-2 py-1" key={key}>
          <span className={`h-2 w-2 rounded-full ${colorClass}`} />
          {t(`legend.${key}`)}
        </span>
      ))}
    </div>
  );
}
