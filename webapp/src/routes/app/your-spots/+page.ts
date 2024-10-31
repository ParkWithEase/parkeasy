import { newClient } from '$lib/utils/client';
import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';
import { INTERNAL_SERVER_ERROR } from '$lib/constants';

export const load: PageLoad = async ({ fetch }) => {
    const client = newClient({ fetch });
    const { data, error: err } = await client.GET('/user/spots');

    if (err) {
        switch (err.status) {
            case 401:
                redirect(307, '/auth/login');
                break;
            case 500:
                error(500, {
                    message: err.title || INTERNAL_SERVER_ERROR
                });
                break;
            default:
                error(500, {
                    message: INTERNAL_SERVER_ERROR
                });
        }
    }

    return {
        spots: data,
        hasNext: undefined,
        paging: undefined
    };
};
