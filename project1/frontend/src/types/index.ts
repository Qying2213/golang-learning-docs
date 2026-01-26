export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: string;
  user_id: string;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  due_date?: string;
  created_at: string;
  updated_at: string;
}

export type TaskStatus = 'pending' | 'in_progress' | 'completed';
export type TaskPriority = 'low' | 'medium' | 'high';

export interface AuthResponse {
  success: boolean;
  message?: string;
  data?: {
    token: string;
    user: User;
  };
  error?: string;
}

export interface TaskListResponse {
  success: boolean;
  message?: string;
  data?: {
    tasks: Task[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
  };
  error?: string;
}

export interface TaskResponse {
  success: boolean;
  message?: string;
  data?: Task;
  error?: string;
}

export interface CreateTaskData {
  title: string;
  description: string;
  priority?: TaskPriority;
  due_date?: string;
}

export interface UpdateTaskData {
  title?: string;
  description?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  due_date?: string;
}

export interface TaskFilters {
  status?: TaskStatus;
  priority?: TaskPriority;
  search?: string;
  page?: number;
  page_size?: number;
}
