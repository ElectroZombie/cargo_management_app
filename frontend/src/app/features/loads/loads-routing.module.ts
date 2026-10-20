import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { LoadListComponent } from './load-list/load-list.component';
import { LoadFormComponent } from './load-form/load-form.component';

const routes: Routes = [
  { path: '', component: LoadListComponent },
  { path: 'new', component: LoadFormComponent },
  { path: ':id/edit', component: LoadFormComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class LoadsRoutingModule {}
