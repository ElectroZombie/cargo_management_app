export interface Loader {
  id: number;
  name: string;
  phone: string;
  email: string;
  address: string;
  company: string;
  created_at: string;
  updated_at: string;
}

export interface LoaderDTO {
  name: string;
  phone: string;
  email: string;
  address: string;
  company: string;
}

export interface LoaderFilters {
  name?: string;
  phone?: string;
  email?: string;
  address?: string;
  company?: string;
}
