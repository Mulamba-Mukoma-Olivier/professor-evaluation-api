import apiClient from './client';

export interface Course {
  id: number;
  code: string;
  name: string;
  description: string;
  department: string;
  academic_year: string;
}

export interface CreateCourseRequest {
  name: string;
  code: string;
  credits: number;
  is_active?: boolean;
}

export const coursesApi = {
  getAll: async (): Promise<Course[]> => {
    const response = await apiClient.get<{ courses: Course[] }>('/courses');
    return response.data.courses;
  },

  getById: async (id: number): Promise<Course> => {
    const response = await apiClient.get<Course>(`/courses/${id}`);
    return response.data;
  },

  create: async (data: CreateCourseRequest): Promise<Course> => {
    const response = await apiClient.post<Course>('/courses', data);
    return response.data;
  },
};
