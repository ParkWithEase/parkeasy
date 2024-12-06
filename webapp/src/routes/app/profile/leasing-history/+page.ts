import { newClient } from '$lib/utils/client';
import { handleGetError } from '$lib/utils/error-handler';
import paginate from '$lib/utils/paginate';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    const client = newClient({ fetch });
    const paging = paginate(client, '/user/leasings', { params: { query: { count: 5 } } });
    const pageResult = await paging.next();
    const { data, error: err } = pageResult.value;
    handleGetError(err);

    return {
        leasing_transactions: data ?? [],
        hasNext: !pageResult.done,
        paging: paging
    };
};
