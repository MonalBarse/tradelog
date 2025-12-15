import { atomWithStorage } from 'jotai/utils'

export interface User {
  id: number;
  email: string;
  role: 'user' | 'admin';
}

interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
}

export const authState = atomWithStorage<AuthState>('auth_state', {
  isAuthenticated: false,
  user: null,
  token: null,
})