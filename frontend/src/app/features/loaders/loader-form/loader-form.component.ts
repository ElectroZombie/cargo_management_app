import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { LoaderDTO } from '../../../core/models/loader.model';
import { LoaderService } from '../../../core/services/loader.service';

@Component({
  selector: 'app-loader-form',
  templateUrl: './loader-form.component.html',
  styleUrls: ['./loader-form.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoaderFormComponent implements OnInit {
  readonly form: FormGroup = this.fb.group({
    name: ['', [Validators.required, Validators.maxLength(200)]],
    company: ['', [Validators.maxLength(200)]],
    phone: ['', [Validators.pattern(/^[0-9()+\s.-]{7,20}$/)]],
    email: ['', [Validators.email]],
    address: ['', [Validators.maxLength(500)]],
  });

  isEditing = false;
  loaderId: number | null = null;
  saving = false;

  constructor(
    private readonly fb: FormBuilder,
    private readonly loaderService: LoaderService,
    private readonly route: ActivatedRoute,
    private readonly router: Router
  ) {}

  get f() {
    return this.form.controls;
  }

  ngOnInit(): void {
    const idParam = this.route.snapshot.paramMap.get('id');
    if (idParam) {
      this.loaderId = Number(idParam);
      this.isEditing = true;
      void this.loadLoader(this.loaderId);
    }
  }

  async onSubmit(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving = true;
    const dto = this.form.getRawValue() as LoaderDTO;

    const result =
      this.isEditing && this.loaderId
        ? await this.loaderService.update(this.loaderId, dto)
        : await this.loaderService.create(dto);

    this.saving = false;
    if (result) {
      void this.router.navigate(['/loaders']);
    }
  }

  onCancel(): void {
    void this.router.navigate(['/loaders']);
  }

  private async loadLoader(id: number): Promise<void> {
    const loader = await this.loaderService.get(id);
    if (loader) {
      this.form.patchValue({
        name: loader.name,
        company: loader.company,
        phone: loader.phone,
        email: loader.email,
        address: loader.address,
      });
    }
  }
}
