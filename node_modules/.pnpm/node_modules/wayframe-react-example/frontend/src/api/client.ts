import {
  Item,
  CreateItemRequest,
  UpdateItemRequest,
  HealthResponse,
  PaginatedResponse,
} from '../types/api';

const API_BASE_URL = '/api';

class APIError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'APIError';
  }
}

async function fetchJSON<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new APIError(response.status, errorData.error || errorData.message || 'Request failed');
  }

  return response.json();
}

export const api = {
  // Health check
  getHealth: async (): Promise<HealthResponse> => {
    return fetchJSON<HealthResponse>(`${API_BASE_URL}/health`);
  },

  // Items
  getItems: async (page = 1, perPage = 10): Promise<PaginatedResponse<Item>> => {
    return fetchJSON<PaginatedResponse<Item>>(
      `${API_BASE_URL}/items?page=${page}&per_page=${perPage}`
    );
  },

  getItem: async (id: number): Promise<Item> => {
    return fetchJSON<Item>(`${API_BASE_URL}/items/${id}`);
  },

  createItem: async (data: CreateItemRequest): Promise<Item> => {
    return fetchJSON<Item>(`${API_BASE_URL}/items`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  updateItem: async (id: number, data: UpdateItemRequest): Promise<Item> => {
    return fetchJSON<Item>(`${API_BASE_URL}/items/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  deleteItem: async (id: number): Promise<void> => {
    await fetch(`${API_BASE_URL}/items/${id}`, {
      method: 'DELETE',
    });
  },
};

export { APIError };

