import { useTranslation } from 'react-i18next';

import type { AvailableMission, Kingdom, KingdomEvent, Mission, Neighbor, PatronPressure, PatronStatus, Raid } from '../../api/client';
import type { GameMode } from '../game/types';
import { localMapNodes } from './mapNodes';
import { MapLegend } from './MapLegend';
import type { LocalMapNode, LocalMapNodeID } from './types';

type LocalMapProps = {
  activeEvents: KingdomEvent[];
  activeMissions: Mission[];
  activeRaids: Raid[];
  availableMissions: AvailableMission[];
  currentMode: GameMode;
  kingdom: Kingdom;
  neighbors: Neighbor[];
  patronPressure: PatronPressure | null;
  patronStatus: PatronStatus | null;
  selectedNodeID: LocalMapNodeID;
  unreadReportsCount: number;
  onNodeSelect: (node: LocalMapNode, neighborID?: string) => void;
};

function nodeColor(type: LocalMapNode['type']) {
  switch (type) {
    case 'home':
      return 'border-dusk-gold text-dusk-gold';
    case 'pve':
      return 'border-emerald-300 text-emerald-200';
    case 'neighbor':
      return 'border-red-300 text-red-200';
    case 'patron':
      return 'border-violet-300 text-violet-200';
    case 'event':
      return 'border-sky-300 text-sky-200';
  }
}

function missionNodeNameKey(nodeID: LocalMapNodeID) {
  if (nodeID === 'black_forest') {
    return 'missions:black_forest_expedition.name';
  }
  if (nodeID === 'old_kurgan') {
    return 'missions:old_kurgan_expedition.name';
  }
  if (nodeID === 'dry_ford') {
    return 'missions:dry_ford_scouting.name';
  }
  return null;
}

export function LocalMap({
  activeEvents,
  activeMissions,
  activeRaids,
  availableMissions,
  currentMode,
  kingdom,
  neighbors,
  onNodeSelect,
  patronPressure,
  patronStatus,
  selectedNodeID,
  unreadReportsCount,
}: LocalMapProps) {
  const { t } = useTranslation(['map', 'missions', 'raids']);
  const selectedNode = localMapNodes.find((node) => node.id === selectedNodeID) ?? localMapNodes[0];
  const neighborSlots = [neighbors[0], neighbors[1]];

  function nodeTitle(node: LocalMapNode) {
    const missionNameKey = missionNodeNameKey(node.id);
    if (missionNameKey) {
      return t(missionNameKey);
    }
    if (node.id === 'home') {
      return kingdom.name;
    }
    if (node.id === 'neighbor_1') {
      return neighborSlots[0]?.name ?? t('map:nodes.neighbor_1.name');
    }
    if (node.id === 'neighbor_2') {
      return neighborSlots[1]?.name ?? t('map:nodes.neighbor_2.name');
    }
    return t(`map:nodes.${node.id}.name`);
  }

  function nodeDescription(node: LocalMapNode) {
    if (node.id === 'neighbor_1' && neighborSlots[0]) {
      return t('map:nodes.neighbor_1.knownDescription', { power: t(`raids:power.${neighborSlots[0].powerEstimate}`) });
    }
    if (node.id === 'neighbor_2' && neighborSlots[1]) {
      return t('map:nodes.neighbor_2.knownDescription', { power: t(`raids:power.${neighborSlots[1].powerEstimate}`) });
    }
    return t(`map:nodes.${node.id}.description`);
  }

  function nodeStatuses(node: LocalMapNode) {
    const statuses: string[] = [];
    if (node.missionKey && activeMissions.some((mission) => mission.missionKey === node.missionKey)) {
      statuses.push(t('map:status.activeMission'));
    }
    if (node.type === 'pve' && node.missionKey && availableMissions.some((mission) => mission.key === node.missionKey)) {
      statuses.push(t('map:status.availableMission'));
    }
    if (node.type === 'neighbor' && activeRaids.length > 0) {
      statuses.push(t('map:status.activeRaid'));
    }
    if (node.id === 'omens' && activeEvents.length > 0) {
      statuses.push(t('map:status.activeEvents', { count: activeEvents.length }));
    }
    if (node.id === 'omens' && unreadReportsCount > 0) {
      statuses.push(t('map:status.unreadReports', { count: unreadReportsCount }));
    }
    if (node.id === 'patron_road' && patronPressure) {
      statuses.push(t('map:status.patronPressure'));
    }
    if (node.id === 'patron_road' && patronStatus?.patron) {
      statuses.push(t('map:status.patronChosen'));
    }
    if (node.type === 'neighbor' && neighbors.length === 0) {
      statuses.push(t('map:status.noNeighbors'));
    }
    return statuses.slice(0, 2);
  }

  function selectNode(node: LocalMapNode) {
    const neighbor = node.id === 'neighbor_1' ? neighborSlots[0] : node.id === 'neighbor_2' ? neighborSlots[1] : undefined;
    onNodeSelect(node, neighbor?.kingdomId);
  }

  return (
    <section className="grid min-w-0 max-w-full gap-3 overflow-hidden rounded border border-stone-800 bg-dusk-900/70 p-4">
      <div className="flex min-w-0 flex-wrap items-start justify-between gap-3">
        <div className="min-w-0 max-w-full">
          <p className="text-xs uppercase tracking-normal text-dusk-gold">{t('map:subtitle')}</p>
          <h2 className="break-words text-2xl font-semibold text-stone-100">{t('map:title')}</h2>
        </div>
        <div className="min-w-0 break-words rounded border border-stone-800 bg-dusk-950 px-3 py-2 text-sm text-stone-300">
          {t(`map:actions.${selectedNode.mode}`)}
        </div>
      </div>

      <div className="relative min-h-[22rem] overflow-hidden rounded border border-stone-800 bg-dusk-950 sm:min-h-[29rem]">
        <svg aria-hidden="true" className="h-full min-h-[22rem] w-full sm:min-h-[29rem]" preserveAspectRatio="xMidYMid meet" viewBox="0 0 100 100">
          <defs>
            <radialGradient cx="50%" cy="48%" id="mapGlow" r="65%">
              <stop offset="0%" stopColor="#3f3426" stopOpacity="0.55" />
              <stop offset="55%" stopColor="#17130f" stopOpacity="0.85" />
              <stop offset="100%" stopColor="#0b0908" stopOpacity="1" />
            </radialGradient>
          </defs>
          <rect fill="url(#mapGlow)" height="100" width="100" />
          <path d="M50 52 C38 40 30 32 20 24" fill="none" stroke="#78716c" strokeDasharray="2 3" strokeOpacity="0.45" strokeWidth="0.7" />
          <path d="M50 52 C58 39 62 30 67 22" fill="none" stroke="#78716c" strokeDasharray="2 3" strokeOpacity="0.45" strokeWidth="0.7" />
          <path d="M50 52 C38 61 32 68 24 76" fill="none" stroke="#78716c" strokeDasharray="2 3" strokeOpacity="0.45" strokeWidth="0.7" />
          <path d="M50 52 C65 45 75 42 84 42" fill="none" stroke="#7f1d1d" strokeDasharray="2 3" strokeOpacity="0.45" strokeWidth="0.7" />
          <path d="M50 52 C66 58 73 66 80 72" fill="none" stroke="#7f1d1d" strokeDasharray="2 3" strokeOpacity="0.45" strokeWidth="0.7" />
          <path d="M50 52 C54 65 58 76 62 84" fill="none" stroke="#6d28d9" strokeDasharray="2 3" strokeOpacity="0.45" strokeWidth="0.7" />
          <path d="M50 52 C48 38 46 25 46 13" fill="none" stroke="#0369a1" strokeDasharray="2 3" strokeOpacity="0.45" strokeWidth="0.7" />
          <circle cx="50" cy="52" fill="#d8a94a" opacity="0.12" r="18" />
        </svg>

        {localMapNodes.map((node) => {
          const statuses = nodeStatuses(node);
          const isSelected = selectedNodeID === node.id;
          return (
            <button
              aria-pressed={isSelected}
              className={`absolute max-w-[5.5rem] -translate-x-1/2 -translate-y-1/2 rounded border bg-dusk-950/95 px-2 py-1 text-center text-[0.65rem] shadow transition hover:bg-dusk-800 focus:outline-none focus:ring-2 focus:ring-dusk-gold sm:max-w-[8.5rem] sm:px-3 sm:py-2 sm:text-xs ${nodeColor(node.type)} ${
                isSelected ? 'ring-2 ring-dusk-gold' : ''
              }`}
              key={node.id}
              onClick={() => selectNode(node)}
              style={{ left: `${node.x}%`, top: `${node.y}%` }}
              type="button"
            >
              <span className="block break-words font-semibold">{nodeTitle(node)}</span>
              {statuses.length > 0 ? (
                <span className="mt-1 hidden flex-wrap justify-center gap-1 sm:flex">
                  {statuses.map((status) => (
                    <span className="break-words rounded bg-dusk-800 px-1.5 py-0.5 text-[0.65rem] text-stone-200" key={status}>
                      {status}
                    </span>
                  ))}
                </span>
              ) : null}
            </button>
          );
        })}
      </div>

      <div className="grid min-w-0 gap-3 rounded border border-stone-800 bg-dusk-950/90 p-3">
        <div className="min-w-0">
          <h3 className="break-words font-semibold text-stone-100">{nodeTitle(selectedNode)}</h3>
          <p className="mt-1 break-words text-sm text-stone-400">{nodeDescription(selectedNode)}</p>
          {currentMode !== selectedNode.mode ? <p className="mt-2 break-words text-xs text-dusk-gold">{t(`map:actions.${selectedNode.mode}`)}</p> : null}
        </div>
        <MapLegend />
      </div>
    </section>
  );
}
