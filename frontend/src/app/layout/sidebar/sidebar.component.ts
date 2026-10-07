import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

export interface NavItem {
  label: string;
  icon: string;
  route: string;
}

@Component({
  selector: 'app-sidebar',
  template: `
    <aside class="sidebar" [class.sidebar--collapsed]="collapsed">
      <nav class="sidebar__nav">
        <a
          *ngFor="let item of items"
          class="sidebar__link"
          [routerLink]="item.route"
          routerLinkActive="sidebar__link--active"
          [routerLinkActiveOptions]="{ exact: item.route === '/dashboard' }"
          [title]="item.label"
        >
          <span class="sidebar__icon" aria-hidden="true">{{ item.icon }}</span>
          <span class="sidebar__label" *ngIf="!collapsed">{{ item.label }}</span>
        </a>
      </nav>
    </aside>
  `,
  styleUrls: ['./sidebar.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SidebarComponent {
  readonly items: NavItem[] = [
    { label: 'Dashboard', icon: '▦', route: '/dashboard' },
    { label: 'Drivers', icon: '🧑', route: '/drivers' },
    { label: 'Loaders', icon: '🏗', route: '/loaders' },
    { label: 'Vehicles', icon: '🚚', route: '/vehicles' },
    { label: 'Wells', icon: '⛽', route: '/wells' },
    { label: 'Loads', icon: '📦', route: '/loads' },
    { label: 'Reports', icon: '📈', route: '/reports' },
    { label: 'Settings', icon: '⚙', route: '/settings' },
  ];

  @Input() collapsed = false;
}
