import api from './api';
import type {
  Task,
  TaskListResponse,
  TaskResponse,
  CreateTaskData,
  UpdateTaskData,
  TaskFilters,
} from '../types';

export const taskService = {
  getTasks: async (filters?: TaskFilters) => {
    const params = new URLSearchParams();
    if (filters?.status) params.append('status', filters.status);
    if (filters?.priority) params.append('priority', filters.priority);
    if (filters?.search) params.append('search', filters.search);
    if (filters?.page) params.append('page', filters.page.toString());
    if (filters?.page_size) params.append('page_size', filters.page_size.toString());

    const response = await api.get<TaskListResponse>(`/tasks?${params.toString()}`);
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || 'Failed to fetch tasks');
  },

  getTask: async (id: string): Promise<Task> => {
    const response = await api.get<TaskResponse>(`/tasks/${id}`);
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || 'Failed to fetch task');
  },

  createTask: async (data: CreateTaskData): Promise<Task> => {
    const response = await api.post<TaskResponse>('/tasks', data);
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || 'Failed to create task');
  },

  updateTask: async (id: string, data: UpdateTaskData): Promise<Task> => {
    const response = await api.put<TaskResponse>(`/tasks/${id}`, data);
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || 'Failed to update task');
  },

  deleteTask: async (id: string): Promise<void> => {
    const response = await api.delete(`/tasks/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error || 'Failed to delete task');
    }
  },
};
