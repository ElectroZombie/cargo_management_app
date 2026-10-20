import { Injectable } from '@angular/core';
import { Load, LoadDTO, LoadFilters } from '../models/load.model';
import { ListResponse } from '../models/api-response.model';
import { CrudService } from './crud.service';

@Injectable({ providedIn: 'root' })
export class LoadService extends CrudService<Load, LoadDTO, LoadFilters> {
  protected readonly singular = 'Load';
  protected readonly plural = 'Loads';
  protected readonly label = 'Load';

  /** Loads assigned to a specific driver. */
  async byDriver(driverId: number): Promise<Load[]> {
    return this.queryList('GetLoadsByDriver', driverId);
  }

  /** Loads originating from a specific well. */
  async byWell(wellId: number): Promise<Load[]> {
    return this.queryList('GetLoadsByWell', wellId);
  }

  /** Loads filtered by active/inactive status. */
  async byStatus(status: boolean): Promise<Load[]> {
    return this.queryList('GetLoadsByStatus', status);
  }

  private async queryList(method: string, arg: unknown): Promise<Load[]> {
    this.loadingSubject.next(true);
    try {
      const response = await this.wails.invoke<ListResponse<Load>>(method, arg);
      if (response?.success) {
        this.itemsSubject.next(response.data ?? []);
        return response.data ?? [];
      }
      this.notification.error(response?.error ?? 'Failed to load loads.');
      return [];
    } catch (error) {
      this.notification.error('Failed to load loads.');
      console.error(error);
      return [];
    } finally {
      this.loadingSubject.next(false);
    }
  }
}
