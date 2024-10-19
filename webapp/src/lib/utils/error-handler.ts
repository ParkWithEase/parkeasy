import { redirect } from '@sveltejs/kit';
import type { components } from '$lib/sdk/schema';

type ErrorModel = components['schemas']['ErrorModel'];

export function getErrorMessage(errorDetails: ErrorModel): string {
    switch (errorDetails.status) {
        case 401:
            redirect(307, '/auth/login');
            break;
        case 422:
        case 500:
            return (
                errorDetails.errors?.[0].location +
                ' : ' +
                errorDetails.detail +
                ' with value ' +
                errorDetails.errors?.[0].value
            );
        default:
            return 'Something wrong happen';
    }
}
