import { newClient } from '$lib/utils/client';
import { getMonday } from '$lib/utils/datetime-util';
import { handleGetError } from '$lib/utils/error-handler';
import type { PageLoad } from './$types';

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

    const { data: availability, error: errorAvailability } = await client.GET(
        '/spots/{id}/availability',
        {
            params: {
                path: {
                    id: params.id
                },
                query: {
                    availability_start: currentMonday.toISOString()
                }
            }
        }
    );

    handleGetError(errorAvailability);

    return {
        spot: spot_info,
        time_slots: availability
    };
};
