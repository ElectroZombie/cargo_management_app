import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  Output,
} from '@angular/core';
import { ThemeMode } from '../../core/services/theme.service';

@Component({
  selector: 'app-navbar',
  template: `
    <header class="navbar">
      <div class="navbar__left">
        <button
          type="button"
          class="navbar__icon-btn"
          (click)="toggleSidebar.emit()"
          title="Toggle navigation"
        >
          ☰
        </button>
        <span class="navbar__brand">{{ appName }}</span>
      </div>

      <div class="navbar__search">
        <span class="navbar__search-icon" aria-hidden="true">⌕</span>
        <input
          type="search"
          class="navbar__search-input"
          placeholder="Search loads, drivers, vehicles…"
          [value]="searchTerm"
          (input)="onSearch($any($event.target).value)"
        />
      </div>

      <div class="navbar__right">
        <button
          type="button"
          class="navbar__icon-btn"
          (click)="toggleTheme.emit()"
          [title]="theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'"
        >
          {{ theme === 'dark' ? '☀' : '☾' }}
        </button>
        <a class="navbar__icon-btn" routerLink="/settings" title="Settings">⚙</a>
      </div>
    </header>
  `,
  styleUrls: ['./navbar.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NavbarComponent {
  @Input() appName = 'Cargo Management App';
  @Input() theme: ThemeMode = 'light';
  @Input() searchTerm = '';

  @Output() toggleSidebar = new EventEmitter<void>();
  @Output() toggleTheme = new EventEmitter<void>();
  @Output() search = new EventEmitter<string>();

  onSearch(value: string): void {
    this.search.emit(value);
  }
}
