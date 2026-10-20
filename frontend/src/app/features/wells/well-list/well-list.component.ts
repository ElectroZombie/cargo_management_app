import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { Well } from '../../../core/models/well.model';
import { WellService } from '../../../core/services/well.service';
import { DataTableColumn } from '../../../shared/components/data-table/data-table-column.model';
import {
  PageEvent,
  SortEvent,
} from '../../../shared/components/data-table/data-table.component';

@Component({
  selector: 'app-well-list',
  templateUrl: './well-list.component.html',
  styleUrls: ['./well-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class WellListComponent implements OnInit {
  readonly wells$: Observable<Well[]> = this.wellService.items();
  readonly loading$: Observable<boolean> = this.wellService.isLoading();
  readonly total$: Observable<number> = this.wellService.total();

  readonly columns: DataTableColumn[] = [
    { key: 'name', header: 'Name', sortable: true },
    { key: 'client', header: 'Client', sortable: true },
    { key: 'location', header: 'Location', sortable: true },
  ];

  page = 1;
  pageSize = 25;
  sortBy = 'id';
  sortDir: 'ASC' | 'DESC' = 'ASC';
  searchTerm = '';

  confirmOpen = false;
  pendingDelete: Well | null = null;

  constructor(
    private readonly wellService: WellService,
    private readonly router: Router
  ) {}

  ngOnInit(): void {
    void this.reload();
  }

  onSearch(term: string): void {
    this.searchTerm = term;
    if (term.trim()) {
      void this.wellService.search(term);
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
    void this.router.navigate(['/wells/new']);
  }

  onEdit(well: Well): void {
    void this.router.navigate(['/wells', well.id, 'edit']);
  }

  onDelete(well: Well): void {
    this.pendingDelete = well;
    this.confirmOpen = true;
  }

  async confirmDelete(): Promise<void> {
    if (this.pendingDelete) {
      await this.wellService.remove(this.pendingDelete.id);
    }
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  cancelDelete(): void {
    this.confirmOpen = false;
    this.pendingDelete = null;
  }

  private async reload(): Promise<void> {
    await this.wellService.load(
      {},
      { page: this.page, page_size: this.pageSize, sort_by: this.sortBy, sort_dir: this.sortDir }
    );
  }
}
