import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ApiResponse } from '../models/api-response.model';
import { DashboardStats } from '../models/dashboard.model';
import { NotificationService } from './notification.service';
import { WailsService } from './wails.service';

@Injectable({ providedIn: 'root' })
export class DashboardService {
  private readonly stats$ = new BehaviorSubject<DashboardStats | null>(null);
  private readonly loading$ = new BehaviorSubject<boolean>(false);

  constructor(
    private readonly wails: WailsService,
    private readonly notification: NotificationService
  ) {}

  getStats(): Observable<DashboardStats | null> {
    return this.stats$.asObservable();
  }

  isLoading(): Observable<boolean> {
    return this.loading$.asObservable();
  }

  async loadStats(): Promise<void> {
    this.loading$.next(true);
    try {
      const response = await this.wails.invoke<ApiResponse<DashboardStats>>('GetDashboardStats');
      if (response?.success && response.data) {
        this.stats$.next(response.data);
      } else {
        this.notification.error(response?.error ?? 'Failed to load dashboard statistics.');
      }
    } catch (error) {
      this.notification.error('Failed to load dashboard statistics.');
      console.error(error);
    } finally {
      this.loading$.next(false);
    }
  }
}
