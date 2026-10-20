import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { DriverListComponent } from './driver-list/driver-list.component';
import { DriverFormComponent } from './driver-form/driver-form.component';

const routes: Routes = [
  { path: '', component: DriverListComponent },
  { path: 'new', component: DriverFormComponent },
  { path: ':id/edit', component: DriverFormComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class DriversRoutingModule {}
