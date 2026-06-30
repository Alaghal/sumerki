import type { LocalMapNode } from './types';

export const localMapNodes: LocalMapNode[] = [
  { id: 'home', mode: 'city', type: 'home', x: 50, y: 52 },
  { id: 'black_forest', missionKey: 'black_forest_expedition', mode: 'missions', type: 'pve', x: 20, y: 24 },
  { id: 'old_kurgan', missionKey: 'old_kurgan_expedition', mode: 'missions', type: 'pve', x: 67, y: 22 },
  { id: 'dry_ford', missionKey: 'dry_ford_scouting', mode: 'missions', type: 'pve', x: 24, y: 76 },
  { id: 'neighbor_1', mode: 'raids', type: 'neighbor', x: 84, y: 42 },
  { id: 'neighbor_2', mode: 'raids', type: 'neighbor', x: 80, y: 72 },
  { id: 'patron_road', mode: 'patron', type: 'patron', x: 62, y: 84 },
  { id: 'omens', mode: 'events', type: 'event', x: 46, y: 13 },
];
