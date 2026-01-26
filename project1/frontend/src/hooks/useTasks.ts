import { useState, useCallback } from 'react';
import type { Task, TaskFilters, CreateTaskData, UpdateTaskData } from '../types';
import { taskService } from '../services/tasks';

interface TasksState {
  tasks: Task[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  isLoading: boolean;
  error: string | null;
}

export const useTasks = () => {
  const [state, setState] = useState<TasksState>({
    tasks: [],
    total: 0,
    page: 1,
    pageSize: 10,
    totalPages: 0,
    isLoading: false,
    error: null,
  });

  const fetchTasks = useCallback(async (filters?: TaskFilters) => {
    setState((prev) => ({ ...prev, isLoading: true, error: null }));
    try {
      const data = await taskService.getTasks(filters);
      setState({
        tasks: data.tasks,
        total: data.total,
        page: data.page,
        pageSize: data.page_size,
        totalPages: data.total_pages,
        isLoading: false,
        error: null,
      });
    } catch (err) {
      setState((prev) => ({
        ...prev,
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to fetch tasks',
      }));
    }
  }, []);

  const createTask = useCallback(async (data: CreateTaskData) => {
    setState((prev) => ({ ...prev, isLoading: true, error: null }));
    try {
      const newTask = await taskService.createTask(data);
      setState((prev) => ({
        ...prev,
        tasks: [newTask, ...prev.tasks],
        total: prev.total + 1,
        isLoading: false,
      }));
      return newTask;
    } catch (err) {
      setState((prev) => ({
        ...prev,
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to create task',
      }));
      throw err;
    }
  }, []);

  const updateTask = useCallback(async (id: string, data: UpdateTaskData) => {
    setState((prev) => ({ ...prev, isLoading: true, error: null }));
    try {
      const updatedTask = await taskService.updateTask(id, data);
      setState((prev) => ({
        ...prev,
        tasks: prev.tasks.map((task) => (task.id === id ? updatedTask : task)),
        isLoading: false,
      }));
      return updatedTask;
    } catch (err) {
      setState((prev) => ({
        ...prev,
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to update task',
      }));
      throw err;
    }
  }, []);

  const deleteTask = useCallback(async (id: string) => {
    setState((prev) => ({ ...prev, isLoading: true, error: null }));
    try {
      await taskService.deleteTask(id);
      setState((prev) => ({
        ...prev,
        tasks: prev.tasks.filter((task) => task.id !== id),
        total: prev.total - 1,
        isLoading: false,
      }));
    } catch (err) {
      setState((prev) => ({
        ...prev,
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to delete task',
      }));
      throw err;
    }
  }, []);

  const clearError = useCallback(() => {
    setState((prev) => ({ ...prev, error: null }));
  }, []);

  return {
    ...state,
    fetchTasks,
    createTask,
    updateTask,
    deleteTask,
    clearError,
  };
};
