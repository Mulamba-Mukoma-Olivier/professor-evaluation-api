import apiClient from './client';

export interface ProfessorResult {
  professor_id: number;
  professor_name: string;
  average_score: number;
  total_evaluations: number;
  criteria_averages: Record<string, number>;
}

export const resultsApi = {
  getProfessorResults: async (professorId: number): Promise<ProfessorResult> => {
    const response = await apiClient.get<ProfessorResult>(`/results/professors/${professorId}`);
    return response.data;
  },
};
