import { Injectable } from '@angular/core';
import { Vehicle, VehicleDTO, VehicleFilters } from '../models/vehicle.model';
import { CrudService } from './crud.service';

@Injectable({ providedIn: 'root' })
export class VehicleService extends CrudService<Vehicle, VehicleDTO, VehicleFilters> {
  protected readonly singular = 'Vehicle';
  protected readonly plural = 'Vehicles';
  protected readonly label = 'Vehicle';
}
