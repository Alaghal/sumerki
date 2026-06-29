const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';
const AUTH_TOKEN_KEY = 'sumerki.auth.token';

export type User = {
  id: string;
  email: string;
};

export type Culture = 'northern_principality' | 'lizard_grad' | 'free_posad';
export type Patron = 'independent' | 'empire_of_dusk' | 'old_pact';

export type Kingdom = {
  id: string;
  userId: string;
  name: string;
  culture: Culture;
  patron: Patron | null;
  createdAt: string;
  updatedAt: string;
};

export type HealthStatus = 'healthy' | 'wounded' | 'sick';

export type Ruler = {
  id: string;
  kingdomId: string;
  name: string;
  age: number;
  culture: Culture;
  authority: number;
  courage: number;
  cunning: number;
  honor: number;
  cruelty: number;
  ambition: number;
  paranoia: number;
  healthStatus: HealthStatus;
  createdAt: string;
  updatedAt: string;
};

export type ResourceValues = {
  gold: number;
  food: number;
  wood: number;
  stone: number;
  population: number;
};

export type Resources = ResourceValues & {
  kingdomId: string;
  productionPerHour: ResourceValues;
  lastCalculatedAt: string;
  updatedAt: string;
};

export type BuildingType =
  | 'town_hall'
  | 'farm'
  | 'lumberyard'
  | 'quarry'
  | 'barracks'
  | 'market'
  | 'walls'
  | 'shrine';

export type BuildingNextUpgrade = {
  targetLevel: number;
  cost: ResourceValues;
  durationSeconds: number;
  canUpgrade: boolean;
  blockedReason: string | null;
};

export type Building = {
  id: string;
  kingdomId: string;
  type: BuildingType;
  label: string;
  level: number;
  maxLevel: number;
  isUpgrading: boolean;
  upgradeStartedAt: string | null;
  upgradeFinishesAt: string | null;
  nextUpgrade: BuildingNextUpgrade | null;
  effects: string[];
  createdAt: string;
  updatedAt: string;
};

export type UnitType = 'militia' | 'spearmen' | 'archers' | 'cavalry' | 'scouts';

export type UnitStats = {
  attack: number;
  defense: number;
  speed: number;
  supply: number;
};

export type UnitRequirements = {
  barracksLevel: number;
  isMet: boolean;
};

export type Unit = {
  type: UnitType;
  label: string;
  amount: number;
  stats: UnitStats;
  cost: ResourceValues;
  secondsPerUnit: number;
  requirements: UnitRequirements;
};

export type TrainingOrder = {
  id: string;
  unitType: UnitType;
  unitLabel: string;
  amount: number;
  status: 'training' | 'completed';
  startedAt: string;
  finishesAt: string;
  completedAt: string | null;
};

export type ArmySummary = {
  totalUnits: number;
  totalAttack: number;
  totalDefense: number;
  totalSupply: number;
};

export type Army = {
  kingdomId: string;
  units: Unit[];
  trainingOrders: TrainingOrder[];
  summary: ArmySummary;
};

type AuthResponse = {
  user: User;
  token: string;
};

type MeResponse = {
  user: User;
};

type KingdomResponse = {
  kingdom: Kingdom | null;
};

type RulerResponse = {
  ruler: Ruler;
};

type ResourcesResponse = {
  resources: Resources;
};

type BuildingsResponse = {
  buildings: Building[];
};

type BuildingUpgradeResponse = {
  building: Building;
  resources: Resources;
};

type ArmyResponse = {
  army: Army;
};

type TrainUnitsRequest = {
  unitType: UnitType;
  amount: number;
};

type TrainUnitsResponse = {
  trainingOrder: TrainingOrder;
  resources: Resources;
};

type ApiErrorResponse = {
  error?: {
    code?: string;
    message?: string;
  };
};

type RequestOptions = {
  method?: string;
  body?: unknown;
  token?: string | null;
};

export class ApiError extends Error {
  code: string;
  status: number;

  constructor(code: string, message: string, status: number) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
  }
}

function readStoredToken(): string | null {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

function storeToken(token: string) {
  localStorage.setItem(AUTH_TOKEN_KEY, token);
}

function clearStoredToken() {
  localStorage.removeItem(AUTH_TOKEN_KEY);
}

async function request<TResponse>(path: string, options: RequestOptions = {}): Promise<TResponse> {
  const headers = new Headers();
  headers.set('Accept', 'application/json');

  if (options.body !== undefined) {
    headers.set('Content-Type', 'application/json');
  }

  const token = options.token === undefined ? readStoredToken() : options.token;
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method ?? 'GET',
    headers,
    body: options.body === undefined ? undefined : JSON.stringify(options.body),
  });

  const data = (await response.json().catch(() => null)) as ApiErrorResponse | TResponse | null;

  if (!response.ok) {
    const errorData = data as ApiErrorResponse | null;
    throw new ApiError(
      errorData?.error?.code ?? 'request_failed',
      errorData?.error?.message ?? 'Request failed',
      response.status,
    );
  }

  return data as TResponse;
}

export function register(email: string, password: string) {
  return request<AuthResponse>('/api/auth/register', {
    method: 'POST',
    body: { email, password },
    token: null,
  });
}

export function login(email: string, password: string) {
  return request<AuthResponse>('/api/auth/login', {
    method: 'POST',
    body: { email, password },
    token: null,
  });
}

export function getMe(token?: string) {
  return request<MeResponse>('/api/me', { token });
}

export function getMyKingdom(token?: string) {
  return request<KingdomResponse>('/api/kingdoms/me', { token });
}

export function createKingdom(name: string, culture: Culture, token?: string) {
  return request<KingdomResponse>('/api/kingdoms', {
    method: 'POST',
    body: { name, culture },
    token,
  });
}

export function getMyRuler(token?: string) {
  return request<RulerResponse>('/api/ruler/me', { token });
}

export function getMyResources(token?: string) {
  return request<ResourcesResponse>('/api/resources/me', { token });
}

export function getMyBuildings(token?: string) {
  return request<BuildingsResponse>('/api/buildings/me', { token });
}

export function upgradeBuilding(type: BuildingType, token?: string) {
  return request<BuildingUpgradeResponse>(`/api/buildings/${type}/upgrade`, {
    method: 'POST',
    token,
  });
}

export function getMyArmy(token?: string) {
  return request<ArmyResponse>('/api/army/me', { token });
}

export function trainUnits(unitType: UnitType, amount: number, token?: string) {
  return request<TrainUnitsResponse>('/api/army/train', {
    method: 'POST',
    body: { unitType, amount } satisfies TrainUnitsRequest,
    token,
  });
}

export { API_BASE_URL, AUTH_TOKEN_KEY, clearStoredToken, readStoredToken, storeToken };
