export interface Item {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface CreateItemRequest {
  name: string;
  description: string;
}

export interface UpdateItemRequest {
  name?: string;
  description?: string;
}

export interface HealthResponse {
  status: string;
  database: string;
  timestamp: string;
}

export interface ErrorResponse {
  error: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  per_page: number;
  total: number;
}

