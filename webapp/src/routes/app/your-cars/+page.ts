import { newClient } from '$lib/utils/client';
import paginate from '$lib/utils/paginate';
import { error, redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { modelToError } from '$lib/utils/error-adapters';

export const ssr = false;

export const load: PageLoad = async ({ fetch }) => {
    const client = newClient({ fetch });
    const paging = paginate(client, '/cars', { params: { query: { count: 5 } } });
    const pageResult = await paging.next();
    const { data, error: err } = pageResult.value;
    if (err) {
        switch (err.status) {
            case 401:
                redirect(307, '/auth/login');
                break;
            default:
                error(err.status ?? 500, modelToError(err));
        }
    }
    return {
        cars: data ?? [],
        hasNext: !pageResult.done,
        paging: paging
    };
};
