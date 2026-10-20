import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { LoaderListComponent } from './loader-list/loader-list.component';
import { LoaderFormComponent } from './loader-form/loader-form.component';

const routes: Routes = [
  { path: '', component: LoaderListComponent },
  { path: 'new', component: LoaderFormComponent },
  { path: ':id/edit', component: LoaderFormComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class LoadersRoutingModule {}
