export interface Vehicle {
  id: number;
  vec_id: string;
  license_id: string;
  trailer_number: number;
  vin: string;
  license_expiration: string;
  driver_id: number;
  overweight_permit_id: string;
  created_at: string;
  updated_at: string;
}

export interface VehicleDTO {
  vec_id: string;
  license_id: string;
  trailer_number: number;
  vin: string;
  license_expiration: string;
  driver_id: number;
  overweight_permit_id: string;
}

export interface VehicleFilters {
  vec_id?: string;
  license_id?: string;
  vin?: string;
  license_expiration?: string;
  driver_id?: number;
  overweight_permit_id?: string;
}
