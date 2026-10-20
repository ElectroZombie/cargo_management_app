import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { WellListComponent } from './well-list/well-list.component';
import { WellFormComponent } from './well-form/well-form.component';

const routes: Routes = [
  { path: '', component: WellListComponent },
  { path: 'new', component: WellFormComponent },
  { path: ':id/edit', component: WellFormComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class WellsRoutingModule {}
