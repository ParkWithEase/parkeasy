// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
import type { components } from '$lib/sdk/schema';

declare global {
    namespace App {
        type Error = components['schemas']['ErrorModel'];
        // interface Locals {}
        // interface PageData {}
        // interface PageState {}
        // interface Platform {}
    }
}

export {};
