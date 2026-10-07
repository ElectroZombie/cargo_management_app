import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  Output,
} from '@angular/core';

@Component({
  selector: 'app-confirm-dialog',
  template: `
    <app-modal
      [open]="open"
      [title]="title"
      size="sm"
      [closeOnBackdrop]="false"
      (closed)="cancelled.emit()"
    >
      <p class="confirm-dialog__message">{{ message }}</p>
      <div modalFooter class="confirm-dialog__actions">
        <app-button variant="ghost" (clicked)="cancelled.emit()">
          {{ cancelText }}
        </app-button>
        <app-button
          [variant]="destructive ? 'danger' : 'primary'"
          [loading]="loading"
          (clicked)="confirmed.emit()"
        >
          {{ confirmText }}
        </app-button>
      </div>
    </app-modal>
  `,
  styleUrls: ['./confirm-dialog.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ConfirmDialogComponent {
  @Input() open = false;
  @Input() title = 'Confirm';
  @Input() message = 'Are you sure?';
  @Input() confirmText = 'Confirm';
  @Input() cancelText = 'Cancel';
  @Input() destructive = false;
  @Input() loading = false;

  @Output() confirmed = new EventEmitter<void>();
  @Output() cancelled = new EventEmitter<void>();
}
