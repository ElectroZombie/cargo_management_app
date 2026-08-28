import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

interface Cargo {
  id: number;
  name: string;
  status: string;
  weight: number;
  destination: string;
  date: string;
}

@Component({
  selector: 'app-cargo-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './cargo-list.component.html',
  styleUrls: ['./cargo-list.component.scss']
})
export class CargoListComponent {
  cargoList: Cargo[] = [
    {
      id: 1,
      name: 'Electronics shipment',
      status: 'In Transit',
      weight: 150,
      destination: 'New York',
      date: '2026-08-28'
    },
    {
      id: 2,
      name: 'Furniture package',
      status: 'Delivered',
      weight: 500,
      destination: 'Los Angeles',
      date: '2026-08-25'
    },
    {
      id: 3,
      name: 'Books and documents',
      status: 'Pending',
      weight: 75,
      destination: 'Chicago',
      date: '2026-08-29'
    }
  ];

  newCargo: Cargo = {
    id: 0,
    name: '',
    status: 'Pending',
    weight: 0,
    destination: '',
    date: new Date().toISOString().split('T')[0]
  };

  showForm = false;

  toggleForm(): void {
    this.showForm = !this.showForm;
  }

  addCargo(): void {
    if (this.newCargo.name && this.newCargo.destination) {
      this.newCargo.id = Math.max(...this.cargoList.map(c => c.id)) + 1;
      this.cargoList.push({ ...this.newCargo });
      this.resetForm();
    }
  }

  resetForm(): void {
    this.newCargo = {
      id: 0,
      name: '',
      status: 'Pending',
      weight: 0,
      destination: '',
      date: new Date().toISOString().split('T')[0]
    };
    this.showForm = false;
  }

  deleteCargo(id: number): void {
    this.cargoList = this.cargoList.filter(c => c.id !== id);
  }

  getStatusClass(status: string): string {
    return `status-${status.toLowerCase().replace(/\s+/g, '-')}`;
  }
}
