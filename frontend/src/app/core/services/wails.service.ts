import { Injectable } from '@angular/core';

/**
 * Thin bridge over the Wails runtime.
 *
 * Supports both the generated `window.go.main.App` bindings and the
 * lower-level `window.runtime.Call.ByName` runtime API so the app works
 * whether or not `wails dev` has regenerated the bindings.
 */
@Injectable({ providedIn: 'root' })
export class WailsService {
  /** Invoke a bound Go method by name and resolve its result. */
  async invoke<T = unknown>(method: string, ...args: unknown[]): Promise<T> {
    const runtime = this.getRuntime();

    if (!runtime) {
      throw new Error(
        'Wails runtime is not available. Run the app through "wails dev" or "wails build".'
      );
    }

    if (typeof runtime === 'object' && typeof runtime[method] === 'function') {
      return (await runtime[method](...args)) as T;
    }

    const callByName = (globalThis as any)?.window?.runtime?.Call?.ByName;
    if (typeof callByName === 'function') {
      return (await callByName(method, ...args)) as T;
    }

    throw new Error(`Wails binding "${method}" is not available.`);
  }

  /** Whether the Wails runtime is currently reachable. */
  isAvailable(): boolean {
    return this.getRuntime() !== null;
  }

  private getRuntime(): Record<string, (...args: unknown[]) => Promise<unknown>> | null {
    const w = (globalThis as any)?.window;
    if (!w) {
      return null;
    }
    return w?.go?.main?.App ?? w?.go?.app?.App ?? w?.go?.App ?? null;
  }
}
