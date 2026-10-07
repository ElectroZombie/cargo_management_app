import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  Output,
} from '@angular/core';

@Component({
  selector: 'app-modal',
  template: `
    <div class="modal-backdrop" *ngIf="open" (click)="onBackdropClick()">
      <div
        class="modal"
        [class.modal--sm]="size === 'sm'"
        [class.modal--lg]="size === 'lg'"
        role="dialog"
        aria-modal="true"
        (click)="$event.stopPropagation()"
      >
        <header class="modal__header">
          <h3 class="modal__title">{{ title }}</h3>
          <button
            type="button"
            class="modal__close"
            aria-label="Close"
            (click)="closed.emit()"
          >
            &times;
          </button>
        </header>

        <div class="modal__body">
          <ng-content></ng-content>
        </div>

        <footer class="modal__footer">
          <ng-content select="[modalFooter]"></ng-content>
        </footer>
      </div>
    </div>
  `,
  styleUrls: ['./modal.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ModalComponent {
  @Input() open = false;
  @Input() title = '';
  @Input() size: 'sm' | 'md' | 'lg' = 'md';
  @Input() closeOnBackdrop = true;

  @Output() closed = new EventEmitter<void>();

  onBackdropClick(): void {
    if (this.closeOnBackdrop) {
      this.closed.emit();
    }
  }
}
