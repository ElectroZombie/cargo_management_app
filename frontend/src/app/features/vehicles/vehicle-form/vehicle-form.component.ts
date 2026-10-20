import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { Driver } from '../../../core/models/driver.model';
import { VehicleDTO } from '../../../core/models/vehicle.model';
import { DriverService } from '../../../core/services/driver.service';
import { VehicleService } from '../../../core/services/vehicle.service';

@Component({
  selector: 'app-vehicle-form',
  templateUrl: './vehicle-form.component.html',
  styleUrls: ['./vehicle-form.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class VehicleFormComponent implements OnInit {
  readonly form: FormGroup = this.fb.group({
    vec_id: ['', [Validators.required, Validators.pattern(/^[A-Za-z0-9._/-]+$/)]],
    license_id: ['', [Validators.required, Validators.pattern(/^[A-Za-z0-9-]{4,30}$/)]],
    vin: ['', [Validators.required, Validators.pattern(/^[A-HJ-NPR-Z0-9]{17}$/)]],
    license_expiration: ['', [Validators.required]],
    driver_id: [null, [Validators.required, Validators.min(1)]],
    trailer_number: [0, [Validators.min(0)]],
    overweight_permit_id: ['', [Validators.pattern(/^[A-Za-z0-9-._/]{2,64}$/)]],
  });

  readonly drivers$: Observable<Driver[]> = this.driverService.items();

  isEditing = false;
  vehicleId: number | null = null;
  saving = false;

  constructor(
    private readonly fb: FormBuilder,
    private readonly vehicleService: VehicleService,
    private readonly driverService: DriverService,
    private readonly route: ActivatedRoute,
    private readonly router: Router
  ) {}

  get f() {
    return this.form.controls;
  }

  ngOnInit(): void {
    void this.driverService.load({}, { page: 1, page_size: 500, sort_by: 'name', sort_dir: 'ASC' });

    const idParam = this.route.snapshot.paramMap.get('id');
    if (idParam) {
      this.vehicleId = Number(idParam);
      this.isEditing = true;
      void this.loadVehicle(this.vehicleId);
    }
  }

  async onSubmit(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving = true;
    const dto = this.form.getRawValue() as VehicleDTO;

    const result =
      this.isEditing && this.vehicleId
        ? await this.vehicleService.update(this.vehicleId, dto)
        : await this.vehicleService.create(dto);

    this.saving = false;
    if (result) {
      void this.router.navigate(['/vehicles']);
    }
  }

  onCancel(): void {
    void this.router.navigate(['/vehicles']);
  }

  private async loadVehicle(id: number): Promise<void> {
    const vehicle = await this.vehicleService.get(id);
    if (vehicle) {
      this.form.patchValue({
        vec_id: vehicle.vec_id,
        license_id: vehicle.license_id,
        vin: vehicle.vin,
        license_expiration: vehicle.license_expiration,
        driver_id: vehicle.driver_id,
        trailer_number: vehicle.trailer_number,
        overweight_permit_id: vehicle.overweight_permit_id,
      });
    }
  }
}
