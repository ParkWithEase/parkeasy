import { newClient } from '$lib/utils/client';
import type { PageLoad } from './$types';
import { handleGetError } from '$lib/utils/error-handler';
import paginate from '$lib/utils/paginate';

export const load: PageLoad = async ({ fetch }) => {
    const client = newClient({ fetch });
    const paging = paginate(client, '/spots/preference', { params: { query: { count: 5 } } });
    const pageResult = await paging.next();
    const { data, error: err } = pageResult.value;
    if (err)
    {
        handleGetError(err);
    }
   
    return {
        spots: data ?? [],
        hasNext: !pageResult.done,
        paging: paging
    };
};
