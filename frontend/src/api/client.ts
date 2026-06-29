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

export { API_BASE_URL, AUTH_TOKEN_KEY, clearStoredToken, readStoredToken, storeToken };
