import apiClient from './client';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  matricule: string;
  name: string;
  email: string;
  password: string;
  role: string;
}

export interface User {
  id: number;
  matricule: string;
  name: string;
  email: string;
  role: string;
}

export interface LoginResponse {
  access_token: string;
  token_type: string;
  user: User;
}

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/auth/login', data);
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<User> => {
    const response = await apiClient.post<{ user: User }>('/auth/register', data);
    return response.data.user;
  },
};
