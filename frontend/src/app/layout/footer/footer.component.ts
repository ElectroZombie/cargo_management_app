import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

@Component({
  selector: 'app-footer',
  template: `
    <footer class="footer">
      <span>{{ appName }}</span>
      <span class="footer__meta">Local desktop edition · v{{ version }}</span>
    </footer>
  `,
  styleUrls: ['./footer.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FooterComponent {
  @Input() appName = 'Cargo Management App';
  @Input() version = '0.1.0';
}
