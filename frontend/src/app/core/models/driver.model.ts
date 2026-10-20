export interface Driver {
  id: number;
  name: string;
  license_number: string;
  phone: string;
  email: string;
  address: string;
  created_at: string;
  updated_at: string;
}

export interface DriverDTO {
  name: string;
  license_number: string;
  phone: string;
  email: string;
  address: string;
}

export interface DriverFilters {
  name?: string;
  license_number?: string;
  phone?: string;
  email?: string;
  address?: string;
}
