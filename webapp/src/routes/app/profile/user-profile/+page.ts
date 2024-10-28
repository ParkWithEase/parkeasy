import type { PageLoad } from './$types';
import { newClient } from '$lib/utils/client';
import { redirect, error } from '@sveltejs/kit';

export const load: PageLoad = async ({ fetch }) => {
    const client = newClient({ fetch });
    const { data, error: err, response } = await client.GET('/user');

    if (err) {
        switch (err.status || response.status) {
            case 401:
                redirect(307, '/auth/login');
                break;
            default:
                error(err.status ?? 500, err.detail);
                break;
        }
    }

    return {
        full_name: data?.full_name,
        email: data?.email
    };
};
