import { Load } from './load.model';

export interface DashboardStats {
  total_loads: number;
  active_loads: number;
  total_net_weight: number;
  total_tons: number;
  total_value: number;
  driver_count: number;
  loader_count: number;
  vehicle_count: number;
  well_count: number;
  recent_loads: Load[];
}
