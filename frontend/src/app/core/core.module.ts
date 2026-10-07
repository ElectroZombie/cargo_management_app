import { NgModule, Optional, SkipSelf } from '@angular/core';
import { CommonModule } from '@angular/common';

/**
 * Singleton services are provided at root level. This module exists to
 * group core concerns and prevent accidental re-imports.
 */
@NgModule({
  imports: [CommonModule],
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parent: CoreModule | null) {
    if (parent) {
      throw new Error('CoreModule is already loaded. Import it only in AppModule.');
    }
  }
}
