import type { PageLoad } from './$types';
import { spots_data } from './mock_data';

export const load: PageLoad = async () => {
    return {
        spots: spots_data,
        hasNext: undefined,
        paging: undefined
    };
};
