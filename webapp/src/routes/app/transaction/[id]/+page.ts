import { newClient } from '$lib/utils/client';
import { handleGetError } from '$lib/utils/error-handler';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
    const client = newClient({ fetch });
    const { data: transaction_info, error: error_transaction } = await client.GET(
        '/bookings/{id}',
        {
            params: {
                path: {
                    id: params.id
                }
            }
        }
    );

    handleGetError(error_transaction);

    const { data: spot_info, error: spot_error } = await client.GET('/spots/{id}', {
        params: { path: { id: transaction_info?.parkingspot_id ?? '' } }
    });

    handleGetError(spot_error);

    return {
        transaction_info: transaction_info,
        spot: spot_info,
        car: transaction_info?.car_details
    };
};
