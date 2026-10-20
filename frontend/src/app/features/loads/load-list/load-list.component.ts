import { ChangeDetectionStrategy, Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Load } from '../../../core/models/load.model';
import { LoadService } from '../../../core/services/load.service';
import { DataTableColumn } from '../../../shared/components/data-table/data-table-column.model';
import {
  PageEvent,
  SortEvent,
} from '../../../shared/components/data-table/data-table.component';

@Component({
  selector: 'app-load-list',
  templateUrl: './load-list.component.html',
  styleUrls: ['./load-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoadListComponent implements OnInit, OnDestroy {
  readonly loads$: Observable<Load[]> = this.loadService.items();
  readonly loading$: Observable<boolean> = this.loadService.isLoading();
  readonly total$: Observable<number> = this.loadService.total();

  readonly columns: DataTableColumn[] = [
    { key: 'date', header: 'Date', type: 'date', sortable: true },
    { key: 'load_number', header: 'Load #', sortable: true },
    { key: 'ticket_number', header: 'Ticket' },
    { key: 'driver_id', header: 'Driver', sortable: true },
    { key: 'vec_id', header: 'Vehicle' },
    { key: 'net_weight', header: 'Net Weight', type: 'number', align: 'right' },
    { key: 'tons', header: 'Tons', type: 'number', align: 'right' },
    { key: 'total_value', header: 'Value', type: 'currency', align: 'right' },
    { key: 'status', header: 'Status', type: 'boolean' },
  ];

  page = 1;
  pageSize = 25;
  sortBy = 'id';
  sortDir: 'ASC' | 'DESC' = 'ASC';
  searchTerm = '';

  confirmOpen = false;
  pendingDelete: Load | null = null;

  private readonly destroy$ = new Subject<void>();

  constructor(
    private readonly loadService: LoadService,
    private readonly route: ActivatedRoute,
    private readonly router: Router
  ) {}

  ngOnInit(): void {
    this.route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((params) => {
      const q = params.get('q') ?? '';
      this.searchTerm = q;
      if (q.trim()) {
        void this.loadService.search(q);
      } else {
        void this.reload();
      }
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  onSearch(term: string): void {
    this.searchTerm = term;
    if (term.trim()) {
      void this.loadService.search(term);
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
    void this.router.navigate(['/loads/new']);
  }

  onEdit(load: Load): void {
    void this.router.navigate(['/loads', load.id, 'edit']);
  }

  onDelete(load: Load): void {
    this.pendingDelete = load;
    this.confirmOpen = true;
  }

  async confirmDelete(): Promise<void> {
    if (this.pendingDelete) {
      await this.loadService.remove(this.pendingDelete.id);
    }
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  cancelDelete(): void {
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  private async reload(): Promise<void> {
    await this.loadService.load(
      {},
      { page: this.page, page_size: this.pageSize, sort_by: this.sortBy, sort_dir: this.sortDir }
    );
  }
}
