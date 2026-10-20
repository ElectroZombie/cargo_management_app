import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';

import { SharedModule } from '../../shared/shared.module';
import { LoadersRoutingModule } from './loaders-routing.module';
import { LoaderListComponent } from './loader-list/loader-list.component';
import { LoaderFormComponent } from './loader-form/loader-form.component';

@NgModule({
  declarations: [LoaderListComponent, LoaderFormComponent],
  imports: [CommonModule, ReactiveFormsModule, SharedModule, LoadersRoutingModule],
})
export class LoadersModule {}
