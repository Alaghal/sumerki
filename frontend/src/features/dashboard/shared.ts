import type { ResourceValues, UnitType } from '../../api/client';

export const rulerStats = [
  'authority',
  'courage',
  'cunning',
  'honor',
  'cruelty',
  'ambition',
  'paranoia',
] as const;

export const resourceRows = ['gold', 'food', 'wood', 'stone', 'population'] as const;

export const costRows = ['gold', 'wood', 'stone'] as const;

export const armyCostRows = ['gold', 'food', 'wood', 'stone', 'population'] as const;

export const unitStatRows = ['attack', 'defense', 'speed', 'supply'] as const;

export const unitTypes: UnitType[] = ['militia', 'spearmen', 'archers', 'cavalry', 'scouts'];

export type ResourceKey = keyof ResourceValues;

export type DashboardFormatters = {
  formatDate: (value: string) => string;
  resourceLabel: (key: ResourceKey) => string;
  unitLabel: (unitType: UnitType) => string;
  resourceList: (values: Partial<ResourceValues>) => string;
  costList: (values: Partial<ResourceValues>, keys?: readonly ResourceKey[]) => string;
  unitList: <TUnit extends { unitType: UnitType; amountSent: number; amountLost: number; amountReturned: number }>(
    units: TUnit[],
  ) => string;
};
