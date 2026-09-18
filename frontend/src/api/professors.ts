import apiClient from './client';

export interface Professor {
  id: number;
  matricule: string;
  first_name: string;
  last_name: string;
  email: string;
  department: string;
  grade: string;
  active: boolean;
}

export interface CreateProfessorRequest {
  name: string;
  email: string;
  department: string;
  is_active?: boolean;
}

export const professorsApi = {
  getAll: async (): Promise<Professor[]> => {
    const response = await apiClient.get<{ professors: Professor[] }>('/professors');
    return response.data.professors;
  },

  getActive: async (): Promise<Professor[]> => {
    const response = await apiClient.get<{ professors: Professor[] }>('/professors/active');
    return response.data.professors;
  },

  getById: async (id: number): Promise<Professor> => {
    const response = await apiClient.get<Professor>(`/professors/${id}`);
    return response.data;
  },

  create: async (data: CreateProfessorRequest): Promise<Professor> => {
    const response = await apiClient.post<Professor>('/professors', data);
    return response.data;
  },
};
