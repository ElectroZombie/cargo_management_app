import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { Vehicle } from '../../../core/models/vehicle.model';
import { VehicleService } from '../../../core/services/vehicle.service';
import { DataTableColumn } from '../../../shared/components/data-table/data-table-column.model';
import {
  PageEvent,
  SortEvent,
} from '../../../shared/components/data-table/data-table.component';

@Component({
  selector: 'app-vehicle-list',
  templateUrl: './vehicle-list.component.html',
  styleUrls: ['./vehicle-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class VehicleListComponent implements OnInit {
  readonly vehicles$: Observable<Vehicle[]> = this.vehicleService.items();
  readonly loading$: Observable<boolean> = this.vehicleService.isLoading();
  readonly total$: Observable<number> = this.vehicleService.total();

  readonly columns: DataTableColumn[] = [
    { key: 'vec_id', header: 'Vehicle ID', sortable: true },
    { key: 'license_id', header: 'License Plate', sortable: true },
    { key: 'vin', header: 'VIN' },
    { key: 'license_expiration', header: 'License Expiration', type: 'date', sortable: true },
    { key: 'driver_id', header: 'Driver', sortable: true },
    { key: 'trailer_number', header: 'Trailer', type: 'number' },
    { key: 'overweight_permit_id', header: 'Permit' },
  ];

  page = 1;
  pageSize = 25;
  sortBy = 'id';
  sortDir: 'ASC' | 'DESC' = 'ASC';
  searchTerm = '';

  confirmOpen = false;
  pendingDelete: Vehicle | null = null;

  constructor(
    private readonly vehicleService: VehicleService,
    private readonly router: Router
  ) {}

  ngOnInit(): void {
    void this.reload();
  }

  onSearch(term: string): void {
    this.searchTerm = term;
    if (term.trim()) {
      void this.vehicleService.search(term);
    } else {
      void this.reload();
    }
  }

  onSort(event: SortEvent): void {
    this.sortBy = event.sortBy;
    this.sortDir = event.sortDir;
    void this.reload();
  }

  onPage(event: PageEvent): void {
    this.page = event.page;
    this.pageSize = event.pageSize;
    void this.reload();
  }

  onCreate(): void {
    void this.router.navigate(['/vehicles/new']);
  }

  onEdit(vehicle: Vehicle): void {
    void this.router.navigate(['/vehicles', vehicle.id, 'edit']);
  }

  onDelete(vehicle: Vehicle): void {
    this.pendingDelete = vehicle;
    this.confirmOpen = true;
  }

  async confirmDelete(): Promise<void> {
    if (this.pendingDelete) {
      await this.vehicleService.remove(this.pendingDelete.id);
    }
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  cancelDelete(): void {
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  private async reload(): Promise<void> {
    await this.vehicleService.load(
      {},
      { page: this.page, page_size: this.pageSize, sort_by: this.sortBy, sort_dir: this.sortDir }
    );
  }
}
