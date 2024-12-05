import { error, redirect } from '@sveltejs/kit';
import type { components } from '$lib/sdk/schema';
import { INTERNAL_SERVER_ERROR } from '$lib/constants';

type ErrorModel = components['schemas']['ErrorModel'];

export function handleGetError(errorDetails?: ErrorModel) {
    if (!errorDetails) {
        return;
    }
    switch (errorDetails.status) {
        case 401:
            redirect(307, '/auth/login');
            break;
        case 404:
        case 422:
            error(errorDetails.status, {
                message: errorDetails.detail ?? INTERNAL_SERVER_ERROR,
                ...errorDetails
            });
            break;
        default:
            error(500, { message: errorDetails.detail ?? INTERNAL_SERVER_ERROR, ...errorDetails });
            break;
    }
}

export function getErrorMessage(errorDetails?: ErrorModel): string {
    if (!errorDetails) {
        return '';
    }
    switch (errorDetails.status) {
        case 401:
            redirect(307, '/auth/login');
            break;
        case 422:
        case 500:
            error(500, { message: errorDetails.detail ?? INTERNAL_SERVER_ERROR, ...errorDetails });
            break;
        default:
            return errorDetails.detail || INTERNAL_SERVER_ERROR;
    }
}
