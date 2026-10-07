import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

@Component({
  selector: 'app-card',
  template: `
    <section class="card" [class.card--padded]="padded">
      <header class="card__header" *ngIf="title || subtitle">
        <div>
          <h3 class="card__title" *ngIf="title">{{ title }}</h3>
          <p class="card__subtitle" *ngIf="subtitle">{{ subtitle }}</p>
        </div>
        <div class="card__actions">
          <ng-content select="[cardActions]"></ng-content>
        </div>
      </header>
      <div class="card__body">
        <ng-content></ng-content>
      </div>
    </section>
  `,
  styleUrls: ['./card.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CardComponent {
  @Input() title = '';
  @Input() subtitle = '';
  @Input() padded = true;
}
