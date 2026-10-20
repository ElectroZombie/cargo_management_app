import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ComingSoonComponent } from './shared/components/coming-soon/coming-soon.component';

const routes: Routes = [
  { path: '', pathMatch: 'full', redirectTo: 'dashboard' },
  {
    path: 'dashboard',
    loadChildren: () =>
      import('./features/dashboard/dashboard.module').then((m) => m.DashboardModule),
  },
  {
    path: 'drivers',
    loadChildren: () => import('./features/drivers/drivers.module').then((m) => m.DriversModule),
  },
  {
    path: 'loaders',
    loadChildren: () => import('./features/loaders/loaders.module').then((m) => m.LoadersModule),
  },
  {
    path: 'vehicles',
    loadChildren: () =>
      import('./features/vehicles/vehicles.module').then((m) => m.VehiclesModule),
  },
  {
    path: 'wells',
    loadChildren: () => import('./features/wells/wells.module').then((m) => m.WellsModule),
  },
  {
    path: 'loads',
    loadChildren: () => import('./features/loads/loads.module').then((m) => m.LoadsModule),
  },
  { path: 'reports', component: ComingSoonComponent, data: { title: 'Reports' } },
  { path: 'settings', component: ComingSoonComponent, data: { title: 'Settings' } },
  { path: '**', redirectTo: 'dashboard' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { scrollPositionRestoration: 'enabled' })],
  exports: [RouterModule],
})
export class AppRoutingModule {}
