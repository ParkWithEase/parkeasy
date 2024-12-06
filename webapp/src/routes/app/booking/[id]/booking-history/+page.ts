import type { PageLoad } from './$types';
import { newClient } from '$lib/utils/client';
import { handleGetError } from '$lib/utils/error-handler';
import paginate from '$lib/utils/paginate';

export const load: PageLoad = async ({ fetch, params }) => {
    const client = newClient({ fetch });
    const paging = paginate(client, '/spots/{id}/bookings', {
        params: { path: { id: params.id }, query: { count: 5 } }
    });
    const pageResult = await paging.next();
    const { data, error: err } = pageResult.value;
    handleGetError(err);

    return {
        booking_transactions: data ?? [],
        hasNext: !pageResult.done,
        paging: paging
    };
};
