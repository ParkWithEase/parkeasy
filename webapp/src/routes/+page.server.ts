import { redirect } from '@sveltejs/kit';

export function load({ cookies }) {
    const session = cookies.get('session');
    if (session == undefined) {
        redirect(307, '/auth/login');
    }
}
