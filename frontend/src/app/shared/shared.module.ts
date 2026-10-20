import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ButtonComponent } from './components/button/button.component';
import { CardComponent } from './components/card/card.component';
import { ModalComponent } from './components/modal/modal.component';
import { LoadingSpinnerComponent } from './components/loading-spinner/loading-spinner.component';
import { ToastComponent } from './components/toast/toast.component';
import { ConfirmDialogComponent } from './components/confirm-dialog/confirm-dialog.component';
import { DataTableComponent } from './components/data-table/data-table.component';
import { ComingSoonComponent } from './components/coming-soon/coming-soon.component';

const COMPONENTS = [
  ButtonComponent,
  CardComponent,
  ModalComponent,
  LoadingSpinnerComponent,
  ToastComponent,
  ConfirmDialogComponent,
  DataTableComponent,
  ComingSoonComponent,
];

@NgModule({
  declarations: [...COMPONENTS],
  imports: [CommonModule],
  exports: [...COMPONENTS, CommonModule],
})
export class SharedModule {}
