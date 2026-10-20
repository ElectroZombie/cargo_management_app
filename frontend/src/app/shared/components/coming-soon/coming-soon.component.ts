import { ChangeDetectionStrategy, Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

/**
 * Lightweight placeholder for sections that are intentionally not implemented
 * yet (for example Reports and Settings).
 */
@Component({
  selector: 'app-coming-soon',
  template: `
    <div class="page">
      <div class="page__header">
        <div>
          <h1>{{ title }}</h1>
          <p class="text-secondary">This section is not available yet.</p>
        </div>
      </div>
      <app-card>
        <p>Coming soon.</p>
      </app-card>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ComingSoonComponent {
  readonly title = (this.route.snapshot.data['title'] as string) || 'Coming Soon';

  constructor(private readonly route: ActivatedRoute) {}
}
