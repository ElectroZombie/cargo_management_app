import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

export type ThemeMode = 'light' | 'dark';

const STORAGE_KEY = 'cargo.theme';

/** Manages the light/dark theme and persists the preference. */
@Injectable({ providedIn: 'root' })
export class ThemeService {
  private readonly theme$ = new BehaviorSubject<ThemeMode>(this.readInitial());

  constructor() {
    this.apply(this.theme$.value);
  }

  getTheme(): Observable<ThemeMode> {
    return this.theme$.asObservable();
  }

  current(): ThemeMode {
    return this.theme$.value;
  }

  toggle(): void {
    this.set(this.theme$.value === 'light' ? 'dark' : 'light');
  }

  set(mode: ThemeMode): void {
    this.theme$.next(mode);
    this.apply(mode);
    try {
      localStorage.setItem(STORAGE_KEY, mode);
    } catch {
      // Storage may be unavailable; the theme still applies for the session.
    }
  }

  private apply(mode: ThemeMode): void {
    const root = document.documentElement;
    root.setAttribute('data-theme', mode);
    root.style.colorScheme = mode;
  }

  private readInitial(): ThemeMode {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored === 'light' || stored === 'dark') {
        return stored;
      }
    } catch {
      // Ignore storage access errors.
    }
    return 'light';
  }
}
