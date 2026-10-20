import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { Driver } from '../../../core/models/driver.model';
import { Loader } from '../../../core/models/loader.model';
import { LoadDTO } from '../../../core/models/load.model';
import { Vehicle } from '../../../core/models/vehicle.model';
import { Well } from '../../../core/models/well.model';
import { DriverService } from '../../../core/services/driver.service';
import { LoaderService } from '../../../core/services/loader.service';
import { LoadService } from '../../../core/services/load.service';
import { VehicleService } from '../../../core/services/vehicle.service';
import { WellService } from '../../../core/services/well.service';
import { positiveNumber } from '../../../shared/validators/positive-number.validator';

@Component({
  selector: 'app-load-form',
  templateUrl: './load-form.component.html',
  styleUrls: ['./load-form.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoadFormComponent implements OnInit {
  readonly form: FormGroup = this.fb.group({
    date: [new Date().toISOString().slice(0, 10), [Validators.required]],
    load_number: [
      '',
      [Validators.required, Validators.pattern(/^[A-Za-z0-9._/-]{3,64}$/)],
    ],
    loader_id: [null, [Validators.required, Validators.min(1)]],
    well_id: [null, [Validators.required, Validators.min(1)]],
    ticket_number: [
      '',
      [Validators.required, Validators.pattern(/^[A-Za-z0-9._/-]{3,64}$/)],
    ],
    miles: [null, [Validators.required, positiveNumber()]],
    driver_id: [null, [Validators.required, Validators.min(1)]],
    vec_id: ['', [Validators.required, Validators.pattern(/^[A-Za-z0-9._/-]+$/)]],
    trailer_number: [0, [Validators.min(0)]],
    net_weight: [null, [Validators.required, positiveNumber()]],
    tons: [null, [Validators.required, positiveNumber()]],
    ton_rate: [null, [Validators.required, positiveNumber()]],
    total_value: [null, [Validators.required, positiveNumber()]],
    status: [true],
    ticket_id: ['', [Validators.pattern(/^[A-Za-z0-9._/-]+$/)]],
  });

  readonly loaders$: Observable<Loader[]> = this.loaderService.items();
  readonly wells$: Observable<Well[]> = this.wellService.items();
  readonly drivers$: Observable<Driver[]> = this.driverService.items();
  readonly vehicles$: Observable<Vehicle[]> = this.vehicleService.items();

  isEditing = false;
  loadId: number | null = null;
  saving = false;

  constructor(
    private readonly fb: FormBuilder,
    private readonly loadService: LoadService,
    private readonly loaderService: LoaderService,
    private readonly wellService: WellService,
    private readonly driverService: DriverService,
    private readonly vehicleService: VehicleService,
    private readonly route: ActivatedRoute,
    private readonly router: Router
  ) {}

  get f() {
    return this.form.controls;
  }

  ngOnInit(): void {
    const options = { page: 1, page_size: 500, sort_by: 'name', sort_dir: 'ASC' as const };
    void this.loaderService.load({}, options);
    void this.wellService.load({}, options);
    void this.driverService.load({}, options);
    void this.vehicleService.load({}, { page: 1, page_size: 500, sort_by: 'vec_id', sort_dir: 'ASC' });

    const idParam = this.route.snapshot.paramMap.get('id');
    if (idParam) {
      this.loadId = Number(idParam);
      this.isEditing = true;
      void this.loadLoad(this.loadId);
    }
  }

  async onSubmit(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving = true;
    const dto = this.form.getRawValue() as LoadDTO;

    const result =
      this.isEditing && this.loadId
        ? await this.loadService.update(this.loadId, dto)
        : await this.loadService.create(dto);

    this.saving = false;
    if (result) {
      void this.router.navigate(['/loads']);
    }
  }

  onCancel(): void {
    void this.router.navigate(['/loads']);
  }

  private async loadLoad(id: number): Promise<void> {
    const load = await this.loadService.get(id);
    if (load) {
      this.form.patchValue({
        date: load.date,
        load_number: load.load_number,
        loader_id: load.loader_id,
        well_id: load.well_id,
        ticket_number: load.ticket_number,
        miles: load.miles,
        driver_id: load.driver_id,
        vec_id: load.vec_id,
        trailer_number: load.trailer_number,
        net_weight: load.net_weight,
        tons: load.tons,
        ton_rate: load.ton_rate,
        total_value: load.total_value,
        status: load.status,
        ticket_id: load.ticket_id,
      });
    }
  }
}
