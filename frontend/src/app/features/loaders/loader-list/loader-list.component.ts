import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { Loader } from '../../../core/models/loader.model';
import { LoaderService } from '../../../core/services/loader.service';
import { DataTableColumn } from '../../../shared/components/data-table/data-table-column.model';
import {
  PageEvent,
  SortEvent,
} from '../../../shared/components/data-table/data-table.component';

@Component({
  selector: 'app-loader-list',
  templateUrl: './loader-list.component.html',
  styleUrls: ['./loader-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoaderListComponent implements OnInit {
  readonly loaders$: Observable<Loader[]> = this.loaderService.items();
  readonly loading$: Observable<boolean> = this.loaderService.isLoading();
  readonly total$: Observable<number> = this.loaderService.total();

  readonly columns: DataTableColumn[] = [
    { key: 'name', header: 'Name', sortable: true },
    { key: 'company', header: 'Company' },
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
  pendingDelete: Loader | null = null;

  constructor(
    private readonly loaderService: LoaderService,
    private readonly router: Router
  ) {}

  ngOnInit(): void {
    void this.reload();
  }

  onSearch(term: string): void {
    this.searchTerm = term;
    if (term.trim()) {
      void this.loaderService.search(term);
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
    void this.router.navigate(['/loaders/new']);
  }

  onEdit(loader: Loader): void {
    void this.router.navigate(['/loaders', loader.id, 'edit']);
  }

  onDelete(loader: Loader): void {
    this.pendingDelete = loader;
    this.confirmOpen = true;
  }

  async confirmDelete(): Promise<void> {
    if (this.pendingDelete) {
      await this.loaderService.remove(this.pendingDelete.id);
    }
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  cancelDelete(): void {
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  private async reload(): Promise<void> {
    await this.loaderService.load(
      {},
      { page: this.page, page_size: this.pageSize, sort_by: this.sortBy, sort_dir: this.sortDir }
    );
  }
}
