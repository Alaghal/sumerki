import { createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';

import {
  ApiError,
  clearStoredToken,
  createKingdom,
  Culture,
  getMe,
  getMyKingdom,
  Kingdom,
  login,
  readStoredToken,
  register,
  storeToken,
  User,
} from '../api/client';

type AuthResult = {
  kingdom: Kingdom | null;
};

type SessionContextValue = {
  token: string | null;
  user: User | null;
  kingdom: Kingdom | null;
  loading: boolean;
  error: string | null;
  registerUser: (email: string, password: string) => Promise<AuthResult>;
  loginUser: (email: string, password: string) => Promise<AuthResult>;
  createUserKingdom: (name: string, culture: Culture) => Promise<Kingdom>;
  logout: () => void;
  clearError: () => void;
};

const SessionContext = createContext<SessionContextValue | null>(null);

const authFailureCodes = new Set([
  'missing_authorization_header',
  'invalid_authorization_header',
  'invalid_token',
  'expired_token',
  'user_not_found',
]);

type SessionProviderProps = {
  children: ReactNode;
};

export function SessionProvider({ children }: SessionProviderProps) {
  const [token, setToken] = useState<string | null>(() => readStoredToken());
  const [user, setUser] = useState<User | null>(null);
  const [kingdom, setKingdom] = useState<Kingdom | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const clearSession = useCallback(() => {
    clearStoredToken();
    setToken(null);
    setUser(null);
    setKingdom(null);
  }, []);

  const loadSession = useCallback(
    async (sessionToken: string): Promise<AuthResult> => {
      try {
        setError(null);
        const meResponse = await getMe(sessionToken);
        const kingdomResponse = await getMyKingdom(sessionToken);
        setUser(meResponse.user);
        setKingdom(kingdomResponse.kingdom);
        return { kingdom: kingdomResponse.kingdom };
      } catch (caughtError) {
        if (caughtError instanceof ApiError && authFailureCodes.has(caughtError.code)) {
          clearSession();
        }

        throw caughtError;
      }
    },
    [clearSession],
  );

  useEffect(() => {
    let isActive = true;

    async function restoreSession() {
      if (!token) {
        setLoading(false);
        return;
      }

      try {
        await loadSession(token);
      } catch (caughtError) {
        if (isActive) {
          setError(caughtError instanceof Error ? caughtError.message : 'Session restore failed');
        }
      } finally {
        if (isActive) {
          setLoading(false);
        }
      }
    }

    restoreSession();

    return () => {
      isActive = false;
    };
  }, [loadSession, token]);

  const registerUser = useCallback(
    async (email: string, password: string): Promise<AuthResult> => {
      setError(null);
      const response = await register(email, password);
      storeToken(response.token);
      setToken(response.token);
      setUser(response.user);
      setKingdom(null);
      await getMe(response.token);
      return { kingdom: null };
    },
    [],
  );

  const loginUser = useCallback(
    async (email: string, password: string): Promise<AuthResult> => {
      setError(null);
      const response = await login(email, password);
      storeToken(response.token);
      setToken(response.token);
      setUser(response.user);
      return loadSession(response.token);
    },
    [loadSession],
  );

  const createUserKingdom = useCallback(
    async (name: string, culture: Culture): Promise<Kingdom> => {
      if (!token) {
        throw new ApiError('missing_token', 'Missing auth token', 401);
      }

      setError(null);
      const response = await createKingdom(name, culture, token);

      if (!response.kingdom) {
        throw new ApiError('kingdom_missing', 'Kingdom was not returned', 500);
      }

      setKingdom(response.kingdom);
      return response.kingdom;
    },
    [token],
  );

  const logout = useCallback(() => {
    clearSession();
    setError(null);
  }, [clearSession]);

  const value = useMemo(
    () => ({
      token,
      user,
      kingdom,
      loading,
      error,
      registerUser,
      loginUser,
      createUserKingdom,
      logout,
      clearError: () => setError(null),
    }),
    [token, user, kingdom, loading, error, registerUser, loginUser, createUserKingdom, logout],
  );

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
}

export function useSession() {
  const value = useContext(SessionContext);

  if (!value) {
    throw new Error('useSession must be used inside SessionProvider');
  }

  return value;
}
