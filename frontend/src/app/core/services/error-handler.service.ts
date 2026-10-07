import { Injectable, Injector } from '@angular/core';
import { ErrorHandler } from '@angular/core';
import { NotificationService } from './notification.service';

/** Centralized handler for uncaught Angular errors. */
@Injectable()
export class AppErrorHandler implements ErrorHandler {
  constructor(private readonly injector: Injector) {}

  handleError(error: unknown): void {
    const notification = this.injector.get(NotificationService);

    let message = 'An unexpected error occurred.';
    if (error instanceof Error) {
      message = error.message;
    } else if (typeof error === 'string') {
      message = error;
    }

    console.error('Unhandled application error:', error);
    notification.error(message);
  }
}
