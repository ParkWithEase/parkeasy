import { newClient } from '$lib/utils/client';
import { handleGetError } from '$lib/utils/error-handler';
import type { PageLoad } from './$types';
import { LATITUDE, LONGITUDE, DISTANCE } from '$lib/constants';


export const load: PageLoad = async ({ fetch }) => {
    const client = newClient({ fetch });

    // Fetch parking spots based on provided latitude and longitude
    const { data: spots, error: errorSpots } = await client.GET('/spots', {
        params: {
            query: {
                latitude: LATITUDE,
                longitude: LONGITUDE,
                distance: DISTANCE
            }
        }
    });

    handleGetError(errorSpots);

    return {
        spots: spots
    };
};
