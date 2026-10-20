import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { DashboardStats } from '../../core/models/dashboard.model';
import { DashboardService } from '../../core/services/dashboard.service';
import { DataTableColumn } from '../../shared/components/data-table/data-table-column.model';

interface StatCard {
  label: string;
  value: string;
  hint: string;
  accent: 'primary' | 'secondary' | 'accent' | 'success';
}

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DashboardComponent implements OnInit {
  readonly stats$: Observable<DashboardStats | null> = this.dashboard.getStats();
  readonly loading$: Observable<boolean> = this.dashboard.isLoading();

  readonly recentColumns: DataTableColumn[] = [
    { key: 'load_number', header: 'Load #', sortable: true },
    { key: 'date', header: 'Date', type: 'date', sortable: true },
    { key: 'ticket_number', header: 'Ticket' },
    { key: 'net_weight', header: 'Net Weight', type: 'number', align: 'right' },
    { key: 'total_value', header: 'Value', type: 'currency', align: 'right' },
    { key: 'status', header: 'Status', type: 'boolean' },
  ];

  constructor(private readonly dashboard: DashboardService) {}

  ngOnInit(): void {
    void this.dashboard.loadStats();
  }

  toCards(stats: DashboardStats | null): StatCard[] {
    if (!stats) {
      return [];
    }
    return [
      {
        label: 'Total Loads',
        value: this.number(stats.total_loads),
        hint: 'All recorded loads',
        accent: 'primary',
      },
      {
        label: 'Active Loads',
        value: this.number(stats.active_loads),
        hint: 'Currently in progress',
        accent: 'secondary',
      },
      {
        label: 'Net Weight',
        value: `${this.number(stats.total_net_weight)} lb`,
        hint: 'Aggregate hauled weight',
        accent: 'accent',
      },
      {
        label: 'Total Value',
        value: this.currency(stats.total_value),
        hint: 'Aggregate load value',
        accent: 'success',
      },
    ];
  }

  countCards(stats: DashboardStats | null): StatCard[] {
    if (!stats) {
      return [];
    }
    return [
      { label: 'Drivers', value: this.number(stats.driver_count), hint: '', accent: 'primary' },
      { label: 'Loaders', value: this.number(stats.loader_count), hint: '', accent: 'secondary' },
      { label: 'Vehicles', value: this.number(stats.vehicle_count), hint: '', accent: 'accent' },
      { label: 'Wells', value: this.number(stats.well_count), hint: '', accent: 'success' },
    ];
  }

  private number(value: number | undefined): string {
    return (value ?? 0).toLocaleString();
  }

  private currency(value: number | undefined): string {
    return (value ?? 0).toLocaleString(undefined, {
      style: 'currency',
      currency: 'USD',
      maximumFractionDigits: 0,
    });
  }
}
