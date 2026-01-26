import React from 'react';
import type { TaskStatus, TaskPriority, TaskFilters as Filters } from '../types';

interface TaskFiltersProps {
  filters: Filters;
  onFiltersChange: (filters: Filters) => void;
}

const TaskFiltersComponent: React.FC<TaskFiltersProps> = ({ filters, onFiltersChange }) => {
  return (
    <div className="flex flex-wrap gap-4 items-center">
      <div className="flex-1 min-w-[200px]">
        <input
          type="text"
          placeholder="Search tasks..."
          value={filters.search || ''}
          onChange={(e) => onFiltersChange({ ...filters, search: e.target.value, page: 1 })}
          className="input"
        />
      </div>

      <div>
        <select
          value={filters.status || ''}
          onChange={(e) => onFiltersChange({ 
            ...filters, 
            status: e.target.value as TaskStatus || undefined,
            page: 1 
          })}
          className="input"
        >
          <option value="">All Status</option>
          <option value="pending">Pending</option>
          <option value="in_progress">In Progress</option>
          <option value="completed">Completed</option>
        </select>
      </div>

      <div>
        <select
          value={filters.priority || ''}
          onChange={(e) => onFiltersChange({ 
            ...filters, 
            priority: e.target.value as TaskPriority || undefined,
            page: 1 
          })}
          className="input"
        >
          <option value="">All Priority</option>
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
        </select>
      </div>

      {(filters.status || filters.priority || filters.search) && (
        <button
          onClick={() => onFiltersChange({ page: 1 })}
          className="btn btn-secondary"
        >
          Clear Filters
        </button>
      )}
    </div>
  );
};

export default TaskFiltersComponent;
