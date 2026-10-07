import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  Output,
} from '@angular/core';
import { DataTableColumn } from './data-table-column.model';

export interface SortEvent {
  sortBy: string;
  sortDir: 'ASC' | 'DESC';
}

export interface PageEvent {
  page: number;
  pageSize: number;
}

@Component({
  selector: 'app-data-table',
  templateUrl: './data-table.component.html',
  styleUrls: ['./data-table.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DataTableComponent<T extends Record<string, unknown>> {
  @Input() columns: DataTableColumn[] = [];
  @Input() rows: T[] = [];
  @Input() loading = false;
  @Input() total = 0;
  @Input() page = 1;
  @Input() pageSize = 25;
  @Input() selectable = false;
  @Input() showActions = true;
  @Input() emptyMessage = 'No records found.';
  @Input() sortBy = 'id';
  @Input() sortDir: 'ASC' | 'DESC' = 'ASC';

  @Output() sortChange = new EventEmitter<SortEvent>();
  @Output() pageChange = new EventEmitter<PageEvent>();
  @Output() selectionChange = new EventEmitter<T[]>();
  @Output() edit = new EventEmitter<T>();
  @Output() delete = new EventEmitter<T>();
  @Output() rowClick = new EventEmitter<T>();

  private selection = new Set<unknown>();
  private allSelected = false;

  get totalPages(): number {
    return Math.max(1, Math.ceil(this.total / this.pageSize));
  }

  get rangeStart(): number {
    return this.total === 0 ? 0 : (this.page - 1) * this.pageSize + 1;
  }

  get rangeEnd(): number {
    return Math.min(this.page * this.pageSize, this.total);
  }

  trackById(index: number, row: T): unknown {
    return row['id'] ?? index;
  }

  isSelected(row: T): boolean {
    return this.selection.has(row['id']);
  }

  onSort(column: DataTableColumn): void {
    if (!column.sortable) {
      return;
    }
    const sortDir: 'ASC' | 'DESC' =
      this.sortBy === column.key && this.sortDir === 'ASC' ? 'DESC' : 'ASC';
    this.sortChange.emit({ sortBy: column.key, sortDir });
  }

  onRowClick(row: T): void {
    this.rowClick.emit(row);
  }

  onEdit(row: T, event: Event): void {
    event.stopPropagation();
    this.edit.emit(row);
  }

  onDelete(row: T, event: Event): void {
    event.stopPropagation();
    this.delete.emit(row);
  }

  onToggleRow(row: T, event: Event): void {
    event.stopPropagation();
    const id = row['id'];
    if (this.selection.has(id)) {
      this.selection.delete(id);
    } else {
      this.selection.add(id);
    }
    this.emitSelection();
  }

  onToggleAll(event: Event): void {
    event.stopPropagation();
    this.allSelected = (event.target as HTMLInputElement).checked;
    this.selection.clear();
    if (this.allSelected) {
      for (const row of this.rows) {
        this.selection.add(row['id']);
      }
    }
    this.emitSelection();
  }

  get isAllSelected(): boolean {
    return this.rows.length > 0 && this.rows.every((row) => this.selection.has(row['id']));
  }

  clearSelection(): void {
    this.selection.clear();
    this.emitSelection();
  }

  previousPage(): void {
    if (this.page > 1) {
      this.pageChange.emit({ page: this.page - 1, pageSize: this.pageSize });
    }
  }

  nextPage(): void {
    if (this.page < this.totalPages) {
      this.pageChange.emit({ page: this.page + 1, pageSize: this.pageSize });
    }
  }

  onPageSizeChange(value: string): void {
    this.pageChange.emit({ page: 1, pageSize: Number(value) });
  }

  formatCell(column: DataTableColumn, row: T): string {
    const value = row[column.key];
    if (value === null || value === undefined || value === '') {
      return '—';
    }
    switch (column.type) {
      case 'date':
        return this.formatDate(String(value));
      case 'currency':
        return this.formatCurrency(Number(value));
      case 'number':
        return this.formatNumber(Number(value));
      case 'boolean':
        return value ? 'Active' : 'Inactive';
      default:
        return String(value);
    }
  }

  isBoolean(column: DataTableColumn): boolean {
    return column.type === 'boolean';
  }

  getBooleanValue(column: DataTableColumn, row: T): boolean {
    return Boolean(row[column.key]);
  }

  private emitSelection(): void {
    const selected = this.rows.filter((row) => this.selection.has(row['id']));
    this.selectionChange.emit(selected);
  }

  private formatDate(value: string): string {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return value;
    }
    return date.toLocaleDateString();
  }

  private formatNumber(value: number): string {
    if (Number.isNaN(value)) {
      return '—';
    }
    return value.toLocaleString(undefined, { maximumFractionDigits: 2 });
  }

  private formatCurrency(value: number): string {
    if (Number.isNaN(value)) {
      return '—';
    }
    return value.toLocaleString(undefined, { style: 'currency', currency: 'USD' });
  }
}
