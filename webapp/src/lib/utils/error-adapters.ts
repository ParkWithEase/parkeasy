import type { components } from '$lib/sdk/schema';

type ErrorModel = components['schemas']['ErrorModel'];

export function modelToError(err: ErrorModel): App.Error {
    return { message: err.detail || err.title || '', ...err };
}
