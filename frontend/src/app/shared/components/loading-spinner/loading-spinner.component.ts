import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

@Component({
  selector: 'app-loading-spinner',
  template: `
    <span class="spinner spinner-{{ size }}" role="status" aria-label="Loading">
      <span class="visually-hidden">Loading…</span>
    </span>
  `,
  styleUrls: ['./loading-spinner.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoadingSpinnerComponent {
  @Input() size: 'sm' | 'md' | 'lg' = 'md';
}
