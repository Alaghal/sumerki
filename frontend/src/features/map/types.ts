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

export type SelectedContext =
  | { kind: 'home' }
  | { kind: 'city' }
  | { kind: 'army' }
  | { kind: 'missions' }
  | { kind: 'mission'; missionKey: string }
  | { kind: 'raids' }
  | { kind: 'neighbor'; kingdomId?: string; slot: 'neighbor_1' | 'neighbor_2' }
  | { kind: 'patron' }
  | { kind: 'events' }
  | { kind: 'reports' };

export function contextFromMapNode(node: LocalMapNode, neighborID?: string): SelectedContext {
  if (node.missionKey) {
    return { kind: 'mission', missionKey: node.missionKey };
  }
  if (node.id === 'neighbor_1' || node.id === 'neighbor_2') {
    return { kind: 'neighbor', kingdomId: neighborID, slot: node.id };
  }
  if (node.id === 'patron_road') {
    return { kind: 'patron' };
  }
  if (node.id === 'omens') {
    return { kind: 'events' };
  }
  return { kind: 'home' };
}
