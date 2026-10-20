import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { DriverDTO } from '../../../core/models/driver.model';
import { DriverService } from '../../../core/services/driver.service';

@Component({
  selector: 'app-driver-form',
  templateUrl: './driver-form.component.html',
  styleUrls: ['./driver-form.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DriverFormComponent implements OnInit {
  readonly form: FormGroup = this.fb.group({
    name: ['', [Validators.required, Validators.maxLength(200)]],
    license_number: [
      '',
      [Validators.required, Validators.pattern(/^[A-Za-z0-9-]{4,30}$/)],
    ],
    phone: ['', [Validators.pattern(/^[0-9()+\s.-]{7,20}$/)]],
    email: ['', [Validators.email]],
    address: ['', [Validators.maxLength(500)]],
  });

  isEditing = false;
  driverId: number | null = null;
  saving = false;

  constructor(
    private readonly fb: FormBuilder,
    private readonly driverService: DriverService,
    private readonly route: ActivatedRoute,
    private readonly router: Router
  ) {}

  get f() {
    return this.form.controls;
  }

  ngOnInit(): void {
    const idParam = this.route.snapshot.paramMap.get('id');
    if (idParam) {
      this.driverId = Number(idParam);
      this.isEditing = true;
      void this.loadDriver(this.driverId);
    }
  }

  async onSubmit(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving = true;
    const dto = this.form.getRawValue() as DriverDTO;

    const result =
      this.isEditing && this.driverId
        ? await this.driverService.update(this.driverId, dto)
        : await this.driverService.create(dto);

    this.saving = false;
    if (result) {
      void this.router.navigate(['/drivers']);
    }
  }

  onCancel(): void {
    void this.router.navigate(['/drivers']);
  }

  private async loadDriver(id: number): Promise<void> {
    const driver = await this.driverService.get(id);
    if (driver) {
      this.form.patchValue({
        name: driver.name,
        license_number: driver.license_number,
        phone: driver.phone,
        email: driver.email,
        address: driver.address,
      });
    }
  }
}
