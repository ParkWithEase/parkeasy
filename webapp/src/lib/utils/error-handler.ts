import { redirect } from '@sveltejs/kit';

export function getErrorMessage(errorDetails): string {
    switch (errorDetails.status) {
        case 401:
            redirect(307, '/auth/login');
            break;
        case 422:
        case 500:
            return (
                errorDetails.errors[0].location +
                ' : ' +
                errorDetails.errors[0].message +
                ' with value ' +
                errorDetails.errors[0].value
            );
        default:
            return 'Something wrong happen';
    }
}
