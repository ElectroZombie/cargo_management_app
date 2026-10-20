import { ChangeDetectionStrategy, Component, OnDestroy, OnInit } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { filter, Subject, takeUntil } from 'rxjs';

interface Crumb {
  label: string;
  url: string;
}

@Component({
  selector: 'app-breadcrumb',
  template: `
    <nav class="breadcrumb" aria-label="Breadcrumb">
      <a class="breadcrumb__item" routerLink="/dashboard">Home</a>
      <ng-container *ngFor="let crumb of crumbs; let last = last">
        <span class="breadcrumb__separator">/</span>
        <a
          class="breadcrumb__item"
          [class.breadcrumb__item--current]="last"
          [routerLink]="last ? null : crumb.url"
        >
          {{ crumb.label }}
        </a>
      </ng-container>
    </nav>
  `,
  styleUrls: ['./breadcrumb.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class BreadcrumbComponent implements OnInit, OnDestroy {
  crumbs: Crumb[] = [];
  private readonly destroy$ = new Subject<void>();

  constructor(private readonly router: Router) {}

  ngOnInit(): void {
    this.router.events
      .pipe(
        filter((event): event is NavigationEnd => event instanceof NavigationEnd),
        takeUntil(this.destroy$)
      )
      .subscribe((event) => this.build(event.urlAfterRedirects));
    this.build(this.router.url);
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private build(url: string): void {
    const segments = url.split('?')[0].split('/').filter((segment) => segment && segment !== 'dashboard');
    let path = '';
    this.crumbs = segments.map((segment) => {
      path += `/${segment}`;
      const isNumeric = /^\d+$/.test(segment);
      return {
        label: isNumeric ? `#${segment}` : this.titleCase(segment),
        url: path,
      };
    });
  }

  private titleCase(value: string): string {
    return value
      .split('-')
      .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
      .join(' ');
  }
}
