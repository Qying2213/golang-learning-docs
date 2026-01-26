import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { useTasks } from '../hooks/useTasks';
import type { Task, TaskFilters, TaskStatus, CreateTaskData, UpdateTaskData } from '../types';
import TaskCard from '../components/TaskCard';
import TaskForm from '../components/TaskForm';
import TaskFiltersComponent from '../components/TaskFilters';

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { 
    tasks, 
    total, 
    page, 
    totalPages, 
    isLoading, 
    error, 
    fetchTasks, 
    createTask, 
    updateTask, 
    deleteTask,
    clearError
  } = useTasks();

  const [filters, setFilters] = useState<TaskFilters>({ page: 1, page_size: 10 });
  const [showForm, setShowForm] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  useEffect(() => {
    fetchTasks(filters);
  }, [filters, fetchTasks]);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleCreateTask = async (data: CreateTaskData) => {
    await createTask(data);
    setShowForm(false);
    fetchTasks(filters);
  };

  const handleUpdateTask = async (data: UpdateTaskData) => {
    if (editingTask) {
      await updateTask(editingTask.id, data);
      setEditingTask(null);
      fetchTasks(filters);
    }
  };

  const handleDeleteTask = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this task?')) {
      await deleteTask(id);
      fetchTasks(filters);
    }
  };

  const handleStatusChange = async (id: string, status: TaskStatus) => {
    await updateTask(id, { status });
    fetchTasks(filters);
  };

  const handlePageChange = (newPage: number) => {
    setFilters({ ...filters, page: newPage });
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Task Manager</h1>
              <p className="text-sm text-gray-600">Welcome back, {user?.name}!</p>
            </div>
            <button onClick={handleLogout} className="btn btn-secondary">
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Actions Bar */}
        <div className="mb-8 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
          <div>
            <h2 className="text-xl font-semibold text-gray-900">My Tasks</h2>
            <p className="text-sm text-gray-600 mt-1">{total} total tasks</p>
          </div>
          <button 
            onClick={() => setShowForm(true)} 
            className="btn btn-primary"
          >
            + Create Task
          </button>
        </div>

        {/* Filters */}
        <div className="mb-6">
          <TaskFiltersComponent 
            filters={filters} 
            onFiltersChange={setFilters} 
          />
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg flex justify-between items-center">
            <span>{error}</span>
            <button onClick={clearError} className="text-red-900 hover:text-red-700">
              ✕
            </button>
          </div>
        )}

        {/* Loading State */}
        {isLoading && tasks.length === 0 && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p className="mt-2 text-gray-600">Loading tasks...</p>
          </div>
        )}

        {/* Empty State */}
        {!isLoading && tasks.length === 0 && (
          <div className="text-center py-12 bg-white rounded-xl shadow-sm">
            <p className="text-gray-600 text-lg">No tasks found.</p>
            <button 
              onClick={() => setShowForm(true)} 
              className="mt-4 btn btn-primary"
            >
              Create Your First Task
            </button>
          </div>
        )}

        {/* Tasks Grid */}
        {tasks.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {tasks.map((task) => (
              <TaskCard
                key={task.id}
                task={task}
                onEdit={setEditingTask}
                onDelete={handleDeleteTask}
                onStatusChange={handleStatusChange}
              />
            ))}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="mt-8 flex justify-center gap-2">
            <button
              onClick={() => handlePageChange(page - 1)}
              disabled={page === 1}
              className="btn btn-secondary disabled:opacity-50"
            >
              Previous
            </button>
            <span className="flex items-center px-4 py-2 text-gray-700">
              Page {page} of {totalPages}
            </span>
            <button
              onClick={() => handlePageChange(page + 1)}
              disabled={page === totalPages}
              className="btn btn-secondary disabled:opacity-50"
            >
              Next
            </button>
          </div>
        )}
      </main>

      {/* Task Form Modal */}
      {showForm && (
        <TaskForm
          onSubmit={handleCreateTask}
          onCancel={() => setShowForm(false)}
          isLoading={isLoading}
        />
      )}

      {editingTask && (
        <TaskForm
          task={editingTask}
          onSubmit={handleUpdateTask}
          onCancel={() => setEditingTask(null)}
          isLoading={isLoading}
        />
      )}
    </div>
  );
};

export default Dashboard;
