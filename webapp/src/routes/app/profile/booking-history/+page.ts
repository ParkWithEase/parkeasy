import type { PageLoad } from './$types';
import { newClient } from '$lib/utils/client';
import { handleGetError } from '$lib/utils/error-handler';
import paginate from '$lib/utils/paginate';

export const load: PageLoad = async ({ fetch }) => {
    const client = newClient({ fetch });
    const paging = paginate(client, '/user/bookings', {params: {query: {count: 5}, path: {id: ''}}});
    const pageResult = await paging.next();
    const { data, error: err } = pageResult.value;
    handleGetError(err);

    return {
        booking_transactions: data ?? [],
        hasNext: !pageResult.done,
        paging: paging
    };
};
