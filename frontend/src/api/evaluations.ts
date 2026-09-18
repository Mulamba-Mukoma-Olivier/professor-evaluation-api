import apiClient from './client';

export interface Evaluation {
  id: number;
  student_id: number;
  professor_id: number;
  course_id: number;
  criteria_scores: Record<number, number>;
  comments: string;
  created_at: string;
}

export interface CreateEvaluationRequest {
  professor_id: number;
  course_id: number;
  criteria_scores: Record<number, number>;
  comments: string;
}

export const evaluationsApi = {
  getAll: async (): Promise<Evaluation[]> => {
    const response = await apiClient.get<{ evaluations: Evaluation[] }>('/evaluations');
    return response.data.evaluations;
  },

  getById: async (id: number): Promise<Evaluation> => {
    const response = await apiClient.get<Evaluation>(`/evaluations/${id}`);
    return response.data;
  },

  create: async (studentId: number, data: CreateEvaluationRequest): Promise<Evaluation> => {
    const response = await apiClient.post<Evaluation>(`/evaluations/${studentId}`, data);
    return response.data;
  },
};
