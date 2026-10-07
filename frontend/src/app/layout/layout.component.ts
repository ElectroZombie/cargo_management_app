import { ChangeDetectionStrategy, Component } from '@angular/core';
import { Router } from '@angular/router';
import { environment } from '../../../environments/environment';
import { ThemeMode, ThemeService } from '../core/services/theme.service';

@Component({
  selector: 'app-layout',
  templateUrl: './layout.component.html',
  styleUrls: ['./layout.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LayoutComponent {
  readonly appName = environment.appName;
  readonly theme$ = this.theme.getTheme();

  sidebarCollapsed = false;
  searchTerm = '';

  constructor(
    private readonly theme: ThemeService,
    private readonly router: Router
  ) {}

  toggleSidebar(): void {
    this.sidebarCollapsed = !this.sidebarCollapsed;
  }

  toggleTheme(): void {
    this.theme.toggle();
  }

  onSearch(term: string): void {
    this.searchTerm = term;
    void this.router.navigate(['/loads'], { queryParams: term ? { q: term } : {} });
  }

  trackTheme(_index: number, value: ThemeMode | null): ThemeMode {
    return value ?? 'light';
  }
}
