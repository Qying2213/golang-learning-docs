import React, { useState, useEffect } from 'react';
import type { Task, TaskPriority, CreateTaskData, UpdateTaskData } from '../types';

interface TaskFormProps {
  task?: Task;
  onSubmit: (data: CreateTaskData | UpdateTaskData) => Promise<void>;
  onCancel: () => void;
  isLoading?: boolean;
}

const TaskForm: React.FC<TaskFormProps> = ({ task, onSubmit, onCancel, isLoading }) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState<TaskPriority>('medium');
  const [dueDate, setDueDate] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      setDescription(task.description || '');
      setPriority(task.priority);
      setDueDate(task.due_date ? task.due_date.split('T')[0] : '');
    }
  }, [task]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!title.trim()) {
      setError('Title is required');
      return;
    }

    try {
      const data: CreateTaskData | UpdateTaskData = {
        title: title.trim(),
        description: description.trim(),
        priority,
        due_date: dueDate ? new Date(dueDate).toISOString() : undefined,
      };
      await onSubmit(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save task');
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-lg">
        <div className="p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">
            {task ? 'Edit Task' : 'Create New Task'}
          </h2>

          {error && (
            <div className="mb-4 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                Title *
              </label>
              <input
                id="title"
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="input"
                placeholder="Enter task title"
                required
              />
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="input min-h-[100px]"
                placeholder="Enter task description (optional)"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="priority" className="block text-sm font-medium text-gray-700 mb-1">
                  Priority
                </label>
                <select
                  id="priority"
                  value={priority}
                  onChange={(e) => setPriority(e.target.value as TaskPriority)}
                  className="input"
                >
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                </select>
              </div>

              <div>
                <label htmlFor="dueDate" className="block text-sm font-medium text-gray-700 mb-1">
                  Due Date
                </label>
                <input
                  id="dueDate"
                  type="date"
                  value={dueDate}
                  onChange={(e) => setDueDate(e.target.value)}
                  className="input"
                />
              </div>
            </div>

            <div className="flex gap-3 pt-4">
              <button
                type="submit"
                disabled={isLoading}
                className="flex-1 btn btn-primary disabled:opacity-50"
              >
                {isLoading ? 'Saving...' : task ? 'Update Task' : 'Create Task'}
              </button>
              <button
                type="button"
                onClick={onCancel}
                className="flex-1 btn btn-secondary"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default TaskForm;
