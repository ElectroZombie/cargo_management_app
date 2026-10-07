export interface DataTableColumn {
  key: string;
  header: string;
  sortable?: boolean;
  align?: 'left' | 'center' | 'right';
  type?: 'text' | 'boolean' | 'date' | 'number' | 'currency';
}
