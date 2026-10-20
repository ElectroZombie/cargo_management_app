import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';

import { SharedModule } from '../../shared/shared.module';
import { DriversRoutingModule } from './drivers-routing.module';
import { DriverListComponent } from './driver-list/driver-list.component';
import { DriverFormComponent } from './driver-form/driver-form.component';

@NgModule({
  declarations: [DriverListComponent, DriverFormComponent],
  imports: [CommonModule, ReactiveFormsModule, SharedModule, DriversRoutingModule],
})
export class DriversModule {}
