import apiClient from './client';

export interface Criterion {
  id: number;
  name: string;
  description: string;
  max_score: number;
  active: boolean;
}

export interface CreateCriterionRequest {
  name: string;
  description: string;
  max_score: number;
  is_active?: boolean;
}

export const criteriaApi = {
  getAll: async (): Promise<Criterion[]> => {
    const response = await apiClient.get<{ criteria: Criterion[] }>('/criteria');
    return response.data.criteria;
  },

  getActive: async (): Promise<Criterion[]> => {
    const response = await apiClient.get<{ criteria: Criterion[] }>('/criteria/active');
    return response.data.criteria;
  },

  getById: async (id: number): Promise<Criterion> => {
    const response = await apiClient.get<Criterion>(`/criteria/${id}`);
    return response.data;
  },

  create: async (data: CreateCriterionRequest): Promise<Criterion> => {
    const response = await apiClient.post<Criterion>('/criteria', data);
    return response.data;
  },
};
