import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

/**
 * Validates that a numeric control holds a value greater than zero.
 * Empty values are left to `Validators.required` so the two messages stay
 * independent.
 */
export function positiveNumber(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const value = control.value;
    if (value === null || value === undefined || value === '') {
      return null;
    }
    return Number(value) > 0 ? null : { positive: true };
  };
}
