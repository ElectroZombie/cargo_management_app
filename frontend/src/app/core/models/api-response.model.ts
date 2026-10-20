/**
 * Standard response envelope returned by every Wails binding.
 * Field names match the Go JSON tags produced by the backend.
 */
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  /**
   * Stable machine-readable classification of `error`, so the UI can react to
   * a failure category without parsing the message text.
   */
  code?: ErrorCode;
}

export interface ListResponse<T> {
  success: boolean;
  data?: T[];
  error?: string;
  total?: number;
  code?: ErrorCode;
}

/** Failure categories reported by the backend binding layer. */
export type ErrorCode =
  | 'validation'
  | 'not_found'
  | 'conflict'
  | 'internal'
  | 'unknown';

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
