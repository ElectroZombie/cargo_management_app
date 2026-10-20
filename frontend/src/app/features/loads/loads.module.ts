import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';

import { SharedModule } from '../../shared/shared.module';
import { LoadsRoutingModule } from './loads-routing.module';
import { LoadListComponent } from './load-list/load-list.component';
import { LoadFormComponent } from './load-form/load-form.component';

@NgModule({
  declarations: [LoadListComponent, LoadFormComponent],
  imports: [CommonModule, ReactiveFormsModule, SharedModule, LoadsRoutingModule],
})
export class LoadsModule {}
