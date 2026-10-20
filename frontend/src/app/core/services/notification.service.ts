import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
  id: number;
  type: ToastType;
  message: string;
  timeout: number;
}

/** Application-wide toast notification store. */
@Injectable({ providedIn: 'root' })
export class NotificationService {
  private readonly toasts$ = new BehaviorSubject<Toast[]>([]);
  private nextId = 1;

  getToasts(): Observable<Toast[]> {
    return this.toasts$.asObservable();
  }

  success(message: string, timeout = 3500): void {
    this.push('success', message, timeout);
  }

  error(message: string, timeout = 6000): void {
    this.push('error', message, timeout);
  }

  info(message: string, timeout = 3500): void {
    this.push('info', message, timeout);
  }

  warning(message: string, timeout = 4500): void {
    this.push('warning', message, timeout);
  }

  dismiss(id: number): void {
    this.toasts$.next(this.toasts$.value.filter((toast) => toast.id !== id));
  }

  clear(): void {
    this.toasts$.next([]);
  }

  private push(type: ToastType, message: string, timeout: number): void {
    if (!message) {
      return;
    }
    const toast: Toast = { id: this.nextId++, type, message, timeout };
    this.toasts$.next([...this.toasts$.value, toast]);

    if (timeout > 0) {
      setTimeout(() => this.dismiss(toast.id), timeout);
    }
  }
}
