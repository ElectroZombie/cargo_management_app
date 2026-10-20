import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { Driver } from '../../../core/models/driver.model';
import { DriverService } from '../../../core/services/driver.service';
import { DataTableColumn } from '../../../shared/components/data-table/data-table-column.model';
import {
  PageEvent,
  SortEvent,
} from '../../../shared/components/data-table/data-table.component';

@Component({
  selector: 'app-driver-list',
  templateUrl: './driver-list.component.html',
  styleUrls: ['./driver-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DriverListComponent implements OnInit {
  readonly drivers$: Observable<Driver[]> = this.driverService.items();
  readonly loading$: Observable<boolean> = this.driverService.isLoading();
  readonly total$: Observable<number> = this.driverService.total();

  readonly columns: DataTableColumn[] = [
    { key: 'name', header: 'Name', sortable: true },
    { key: 'license_number', header: 'License #', sortable: true },
    { key: 'phone', header: 'Phone' },
    { key: 'email', header: 'Email' },
    { key: 'address', header: 'Address' },
  ];

  page = 1;
  pageSize = 25;
  sortBy = 'id';
  sortDir: 'ASC' | 'DESC' = 'ASC';
  searchTerm = '';

  confirmOpen = false;
  pendingDelete: Driver | null = null;

  constructor(
    private readonly driverService: DriverService,
    private readonly router: Router
  ) {}

  ngOnInit(): void {
    void this.reload();
  }

  onSearch(term: string): void {
    this.searchTerm = term;
    if (term.trim()) {
      void this.driverService.search(term);
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
    void this.router.navigate(['/drivers/new']);
  }

  onEdit(driver: Driver): void {
    void this.router.navigate(['/drivers', driver.id, 'edit']);
  }

  onDelete(driver: Driver): void {
    this.pendingDelete = driver;
    this.confirmOpen = true;
  }

  async confirmDelete(): Promise<void> {
    if (this.pendingDelete) {
      await this.driverService.remove(this.pendingDelete.id);
    }
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  cancelDelete(): void {
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  private async reload(): Promise<void> {
    await this.driverService.load(
      {},
      { page: this.page, page_size: this.pageSize, sort_by: this.sortBy, sort_dir: this.sortDir }
    );
  }
}
