import { Injectable } from '@angular/core';
import { Well, WellDTO, WellFilters } from '../models/well.model';
import { CrudService } from './crud.service';

@Injectable({ providedIn: 'root' })
export class WellService extends CrudService<Well, WellDTO, WellFilters> {
  protected readonly singular = 'Well';
  protected readonly plural = 'Wells';
  protected readonly label = 'Well';
}
