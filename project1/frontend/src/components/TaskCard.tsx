import React from 'react';
import type { Task, TaskStatus, TaskPriority } from '../types';

interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (id: string) => void;
  onStatusChange: (id: string, status: TaskStatus) => void;
}

const statusLabels: Record<TaskStatus, string> = {
  pending: 'Pending',
  in_progress: 'In Progress',
  completed: 'Completed',
};

const priorityLabels: Record<TaskPriority, string> = {
  low: 'Low',
  medium: 'Medium',
  high: 'High',
};

const TaskCard: React.FC<TaskCardProps> = ({ task, onEdit, onDelete, onStatusChange }) => {
  const formatDate = (dateString?: string) => {
    if (!dateString) return null;
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const isOverdue = task.due_date && new Date(task.due_date) < new Date() && task.status !== 'completed';

  return (
    <div className={`card hover:shadow-lg transition-shadow ${task.status === 'completed' ? 'opacity-75' : ''}`}>
      <div className="flex justify-between items-start mb-3">
        <h3 className={`text-lg font-semibold text-gray-900 ${task.status === 'completed' ? 'line-through' : ''}`}>
          {task.title}
        </h3>
        <div className="flex gap-2">
          <span className={`badge badge-${task.status.replace('_', '-')}`}>
            {statusLabels[task.status]}
          </span>
          <span className={`badge badge-${task.priority}`}>
            {priorityLabels[task.priority]}
          </span>
        </div>
      </div>

      {task.description && (
        <p className="text-gray-600 mb-4 line-clamp-2">{task.description}</p>
      )}

      <div className="flex justify-between items-center">
        <div className="text-sm text-gray-500">
          {task.due_date && (
            <span className={isOverdue ? 'text-red-600 font-medium' : ''}>
              Due: {formatDate(task.due_date)}
              {isOverdue && ' (Overdue)'}
            </span>
          )}
        </div>

        <div className="flex gap-2">
          <select
            value={task.status}
            onChange={(e) => onStatusChange(task.id, e.target.value as TaskStatus)}
            className="text-sm border border-gray-300 rounded-lg px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
          </select>
          <button
            onClick={() => onEdit(task)}
            className="btn btn-secondary text-sm px-3 py-1"
          >
            Edit
          </button>
          <button
            onClick={() => onDelete(task.id)}
            className="btn btn-danger text-sm px-3 py-1"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
};

export default TaskCard;
