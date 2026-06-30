import type { GameMode } from '../game/types';

export type LocalMapNodeID =
  | 'home'
  | 'black_forest'
  | 'old_kurgan'
  | 'dry_ford'
  | 'neighbor_1'
  | 'neighbor_2'
  | 'patron_road'
  | 'omens';

export type LocalMapNodeType = 'home' | 'pve' | 'neighbor' | 'patron' | 'event';

export type LocalMapNode = {
  id: LocalMapNodeID;
  type: LocalMapNodeType;
  mode: GameMode;
  missionKey?: string;
  x: number;
  y: number;
};
