import { ChangeDetectionStrategy, Component } from '@angular/core';
import { Observable } from 'rxjs';
import { NotificationService, Toast } from '../../../core/services/notification.service';

@Component({
  selector: 'app-toast',
  template: `
    <div class="toast-container" aria-live="polite">
      <div
        *ngFor="let toast of toasts$ | async; trackBy: trackById"
        class="toast toast--{{ toast.type }}"
      >
        <span class="toast__message">{{ toast.message }}</span>
        <button
          type="button"
          class="toast__close"
          aria-label="Dismiss"
          (click)="dismiss(toast.id)"
        >
          &times;
        </button>
      </div>
    </div>
  `,
  styleUrls: ['./toast.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ToastComponent {
  readonly toasts$: Observable<Toast[]> = this.notification.getToasts();

  constructor(private readonly notification: NotificationService) {}

  trackById(_index: number, toast: Toast): number {
    return toast.id;
  }

  dismiss(id: number): void {
    this.notification.dismiss(id);
  }
}
