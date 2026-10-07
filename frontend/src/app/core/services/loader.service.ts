import { Injectable } from '@angular/core';
import { Loader, LoaderDTO, LoaderFilters } from '../models/loader.model';
import { CrudService } from './crud.service';

@Injectable({ providedIn: 'root' })
export class LoaderService extends CrudService<Loader, LoaderDTO, LoaderFilters> {
  protected readonly singular = 'Loader';
  protected readonly plural = 'Loaders';
  protected readonly label = 'Loader';
}
