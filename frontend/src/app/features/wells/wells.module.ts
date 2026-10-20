import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';

import { SharedModule } from '../../shared/shared.module';
import { WellsRoutingModule } from './wells-routing.module';
import { WellListComponent } from './well-list/well-list.component';
import { WellFormComponent } from './well-form/well-form.component';

@NgModule({
  declarations: [WellListComponent, WellFormComponent],
  imports: [CommonModule, ReactiveFormsModule, SharedModule, WellsRoutingModule],
})
export class WellsModule {}
