import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ApiResponse, ListResponse, QueryOptions } from '../models/api-response.model';
import { NotificationService } from './notification.service';
import { WailsService } from './wails.service';

export interface EntityIdentity {
  id: number;
}

/**
 * Shared state + CRUD behaviour for operational entities.
 *
 * Concrete services provide the Wails binding names and the user-facing
 * labels; all request/response and state handling lives here so behaviour
 * stays consistent across drivers, loaders, vehicles, wells, and loads.
 */
@Injectable()
export abstract class CrudService<T extends EntityIdentity, DTO, Filters> {
  protected abstract readonly singular: string;
  protected abstract readonly plural: string;
  protected abstract readonly label: string;

  protected readonly itemsSubject = new BehaviorSubject<T[]>([]);
  protected readonly selectedSubject = new BehaviorSubject<T | null>(null);
  protected readonly loadingSubject = new BehaviorSubject<boolean>(false);
  protected readonly totalSubject = new BehaviorSubject<number>(0);

  constructor(
    protected readonly wails: WailsService,
    protected readonly notification: NotificationService
  ) {}

  items(): Observable<T[]> {
    return this.itemsSubject.asObservable();
  }

  isLoading(): Observable<boolean> {
    return this.loadingSubject.asObservable();
  }

  total(): Observable<number> {
    return this.totalSubject.asObservable();
  }

  selected(): Observable<T | null> {
    return this.selectedSubject.asObservable();
  }

  /** Load a filtered, sorted, paginated page of entities. */
  async load(filters: Filters = {} as Filters, options: QueryOptions = {}): Promise<void> {
    this.loadingSubject.next(true);
    try {
      const response = await this.wails.invoke<ListResponse<T>>(
        `List${this.plural}`,
        filters,
        options
      );
      if (response?.success) {
        this.itemsSubject.next(response.data ?? []);
        this.totalSubject.next(response.total ?? response.data?.length ?? 0);
      } else {
        this.notification.error(response?.error ?? `Failed to load ${this.label}s.`);
      }
    } catch (error) {
      this.notification.error(`Failed to load ${this.label}s.`);
      console.error(error);
    } finally {
      this.loadingSubject.next(false);
    }
  }

  /** Load a single entity by id. */
  async get(id: number): Promise<T | null> {
    this.loadingSubject.next(true);
    try {
      const response = await this.wails.invoke<ApiResponse<T>>(`Get${this.singular}`, id);
      if (response?.success && response.data) {
        this.selectedSubject.next(response.data);
        return response.data;
      }
      this.notification.error(response?.error ?? `${this.label} not found.`);
      return null;
    } catch (error) {
      this.notification.error(`Failed to load ${this.label}.`);
      console.error(error);
      return null;
    } finally {
      this.loadingSubject.next(false);
    }
  }

  /** Create a new entity and refresh the collection. */
  async create(dto: DTO): Promise<T | null> {
    this.loadingSubject.next(true);
    try {
      const response = await this.wails.invoke<ApiResponse<T>>(`Create${this.singular}`, dto);
      if (response?.success && response.data) {
        this.notification.success(`${this.label} created successfully.`);
        await this.refresh();
        return response.data;
      }
      this.notification.error(response?.error ?? `Failed to create ${this.label}.`);
      return null;
    } catch (error) {
      this.notification.error(`Failed to create ${this.label}.`);
      console.error(error);
      return null;
    } finally {
      this.loadingSubject.next(false);
    }
  }

  /** Update an existing entity and refresh the collection. */
  async update(id: number, dto: DTO): Promise<T | null> {
    this.loadingSubject.next(true);
    try {
      const response = await this.wails.invoke<ApiResponse<T>>(
        `Update${this.singular}`,
        id,
        dto
      );
      if (response?.success) {
        this.notification.success(`${this.label} updated successfully.`);
        await this.refresh();
        return response.data ?? null;
      }
      this.notification.error(response?.error ?? `Failed to update ${this.label}.`);
      return null;
    } catch (error) {
      this.notification.error(`Failed to update ${this.label}.`);
      console.error(error);
      return null;
    } finally {
      this.loadingSubject.next(false);
    }
  }

  /** Delete an entity and refresh the collection. */
  async remove(id: number): Promise<boolean> {
    this.loadingSubject.next(true);
    try {
      const response = await this.wails.invoke<ApiResponse<unknown>>(
        `Delete${this.singular}`,
        id
      );
      if (response?.success) {
        this.notification.success(`${this.label} deleted successfully.`);
        await this.refresh();
        return true;
      }
      this.notification.error(response?.error ?? `Failed to delete ${this.label}.`);
      return false;
    } catch (error) {
      this.notification.error(`Failed to delete ${this.label}.`);
      console.error(error);
      return false;
    } finally {
      this.loadingSubject.next(false);
    }
  }

  /** Fuzzy search across the entity's searchable fields. */
  async search(query: string): Promise<T[]> {
    this.loadingSubject.next(true);
    try {
      const response = await this.wails.invoke<ListResponse<T>>(
        `Search${this.plural}`,
        query
      );
      if (response?.success) {
        this.itemsSubject.next(response.data ?? []);
        return response.data ?? [];
      }
      this.notification.error(response?.error ?? `Failed to search ${this.label}s.`);
      return [];
    } catch (error) {
      this.notification.error(`Failed to search ${this.label}s.`);
      console.error(error);
      return [];
    } finally {
      this.loadingSubject.next(false);
    }
  }

  protected async refresh(): Promise<void> {
    await this.load();
  }
}
