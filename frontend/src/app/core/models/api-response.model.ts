/**
 * Standard response envelope returned by every Wails binding.
 * Field names match the Go JSON tags produced by the backend.
 */
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface ListResponse<T> {
  success: boolean;
  data?: T[];
  error?: string;
  total?: number;
}

/** Pagination, sorting, and filtering options shared by list queries. */
export interface QueryOptions {
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_dir?: 'ASC' | 'DESC';
}

export interface Paginated<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}
