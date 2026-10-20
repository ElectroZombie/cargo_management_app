import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { WellDTO } from '../../../core/models/well.model';
import { WellService } from '../../../core/services/well.service';

@Component({
  selector: 'app-well-form',
  templateUrl: './well-form.component.html',
  styleUrls: ['./well-form.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class WellFormComponent implements OnInit {
  readonly form: FormGroup = this.fb.group({
    name: ['', [Validators.required, Validators.maxLength(200)]],
    client: ['', [Validators.maxLength(200)]],
    location: ['', [Validators.required, Validators.maxLength(500)]],
  });

  isEditing = false;
  wellId: number | null = null;
  saving = false;

  constructor(
    private readonly fb: FormBuilder,
    private readonly wellService: WellService,
    private readonly route: ActivatedRoute,
    private readonly router: Router
  ) {}

  get f() {
    return this.form.controls;
  }

  ngOnInit(): void {
    const idParam = this.route.snapshot.paramMap.get('id');
    if (idParam) {
      this.wellId = Number(idParam);
      this.isEditing = true;
      void this.loadWell(this.wellId);
    }
  }

  async onSubmit(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving = true;
    const dto = this.form.getRawValue() as WellDTO;

    const result =
      this.isEditing && this.wellId
        ? await this.wellService.update(this.wellId, dto)
        : await this.wellService.create(dto);

    this.saving = false;
    if (result) {
      void this.router.navigate(['/wells']);
    }
  }

  onCancel(): void {
    void this.router.navigate(['/wells']);
  }

  private async loadWell(id: number): Promise<void> {
    const well = await this.wellService.get(id);
    if (well) {
      this.form.patchValue({
        name: well.name,
        client: well.client,
        location: well.location,
      });
    }
  }
}
