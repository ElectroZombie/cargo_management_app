export interface Load {
  id: number;
  date: string;
  load_number: string;
  loader_id: number;
  well_id: number;
  ticket_number: string;
  miles: number;
  driver_id: number;
  vec_id: string;
  trailer_number: number;
  net_weight: number;
  tons: number;
  ton_rate: number;
  total_value: number;
  status: boolean;
  ticket_id: string;
  created_at: string;
  updated_at: string;
}

export interface LoadDTO {
  date: string;
  load_number: string;
  loader_id: number;
  well_id: number;
  ticket_number: string;
  miles: number;
  driver_id: number;
  vec_id: string;
  trailer_number: number;
  net_weight: number;
  tons: number;
  ton_rate: number;
  total_value: number;
  status: boolean;
  ticket_id: string;
}

export interface LoadFilters {
  date?: string;
  load_number?: string;
  loader_id?: number;
  well_id?: number;
  ticket_number?: string;
  driver_id?: number;
  vec_id?: string;
  status?: boolean;
}
