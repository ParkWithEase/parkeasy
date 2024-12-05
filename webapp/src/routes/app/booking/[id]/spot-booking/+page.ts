import { DAY_IN_A_WEEK } from '$lib/constants';
import { newClient } from '$lib/utils/client';
import { getDateWithDayOffset, getMonday } from '$lib/utils/datetime-util';
import { handleGetError } from '$lib/utils/error-handler';
import paginate from '$lib/utils/paginate';
import type { PageLoad } from '../$types';

export const load: PageLoad = async ({ fetch, params }) => {
    const client = newClient({ fetch });
    const { data: spot_info, error: errorSpotInfo } = await client.GET('/spots/{id}', {
        params: {
            path: {
                id: params.id
            }
        }
    });

    handleGetError(errorSpotInfo);

    const currentMonday = getMonday(new Date(Date.now()));

    const nextMonday = getDateWithDayOffset(currentMonday, DAY_IN_A_WEEK);
    nextMonday.setMinutes(nextMonday.getMinutes() - 30);
    //This is supposed to be a booker perspective availability
    const { data: availability, error: errorAvailability } = await client.GET(
        '/spots/{id}/availability',
        {
            params: {
                path: {
                    id: params.id
                },
                query: {
                    availability_start: currentMonday.toISOString(),
                    availability_end: nextMonday.toISOString()
                }
            }
        }
    );

    handleGetError(errorAvailability);

    const { data: isPreferred, error: errorPreference } = await client.GET(
        '/spots/{id}/preference',
        {
            params: {
                path: {
                    id: params.id
                }
            }
        }
    );

    handleGetError(errorPreference);

    const paging = paginate(client, '/cars', { params: { query: { count: 5 } } });
    const pageResult = await paging.next();
    const { data: cars, error: errorCars } = pageResult.value;
    handleGetError(errorCars);

    return {
        paging: paging,
        cars: cars ?? [],
        hasNext: !pageResult.done,

        spot: spot_info,
        time_slots: availability,
        isPreferred: isPreferred ?? false
    };
};
