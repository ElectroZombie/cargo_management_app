import { Injectable } from '@angular/core';
import { Driver, DriverDTO, DriverFilters } from '../models/driver.model';
import { CrudService } from './crud.service';

@Injectable({ providedIn: 'root' })
export class DriverService extends CrudService<Driver, DriverDTO, DriverFilters> {
  protected readonly singular = 'Driver';
  protected readonly plural = 'Drivers';
  protected readonly label = 'Driver';
}
