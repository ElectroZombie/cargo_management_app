import { ChangeDetectionStrategy, Component, NgZone, OnDestroy, OnInit } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { filter, Subject, takeUntil } from 'rxjs';
import { environment } from '../../environments/environment';
import { ThemeMode, ThemeService } from '../core/services/theme.service';

const MOBILE_QUERY = '(max-width: 768px)';

@Component({
  selector: 'app-layout',
  templateUrl: './layout.component.html',
  styleUrls: ['./layout.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LayoutComponent implements OnInit, OnDestroy {
  readonly appName = environment.appName;
  readonly theme$ = this.theme.getTheme();

  sidebarCollapsed = false;
  isMobile = false;
  searchTerm = '';

  private readonly mobileQuery = window.matchMedia(MOBILE_QUERY);
  private readonly destroy$ = new Subject<void>();

  private readonly onMediaChange = (event: MediaQueryListEvent): void =>
    this.zone.run(() => this.applyViewport(event.matches));

  constructor(
    private readonly theme: ThemeService,
    private readonly router: Router,
    private readonly zone: NgZone
  ) {}

  ngOnInit(): void {
    this.applyViewport(this.mobileQuery.matches);
    this.mobileQuery.addEventListener('change', this.onMediaChange);

    this.router.events
      .pipe(
        filter((event): event is NavigationEnd => event instanceof NavigationEnd),
        takeUntil(this.destroy$)
      )
      .subscribe(() => {
        if (this.isMobile) {
          this.sidebarCollapsed = true;
        }
      });
  }

  ngOnDestroy(): void {
    this.mobileQuery.removeEventListener('change', this.onMediaChange);
    this.destroy$.next();
    this.destroy$.complete();
  }

  toggleSidebar(): void {
    this.sidebarCollapsed = !this.sidebarCollapsed;
  }

  closeSidebar(): void {
    this.sidebarCollapsed = true;
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

  private applyViewport(isMobile: boolean): void {
    this.isMobile = isMobile;
    this.sidebarCollapsed = isMobile;
  }
}
