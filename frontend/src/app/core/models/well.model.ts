export interface Well {
  id: number;
  name: string;
  client: string;
  location: string;
  created_at: string;
  updated_at: string;
}

export interface WellDTO {
  name: string;
  client: string;
  location: string;
}

export interface WellFilters {
  name?: string;
  client?: string;
  location?: string;
}
