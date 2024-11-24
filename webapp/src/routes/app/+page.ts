import { newClient } from '$lib/utils/client';
import { handleGetError } from '$lib/utils/error-handler';
import type { PageLoad } from './$types';
import type { components } from '$lib/sdk/schema';

type ParkingSpot = components['schemas']['ParkingSpot'];

export const load: PageLoad = async ({ fetch }) => {
    const client = newClient({ fetch });

    const latitude = 49.88887;
    const longitude = -97.13449;
    const distance = 10000;

    // Fetch parking spots based on provided latitude and longitude
    const { data: spots, error: errorSpots } = await client.GET('/spots', {
        params: {
            query: {
                latitude: latitude,
                longitude: longitude,
                distance: distance
            }
        }
    });

    handleGetError(errorSpots);

    const uniqueSpots = coalesceListings(spots ?? []);

    return {
        spots: uniqueSpots
    };
};

function coalesceListings(listings: ParkingSpot[]): ParkingSpot[] {
    const uniqueListings: ParkingSpot[] = [];
    const seenIds = new Set<string>();

    for (const listing of listings) {
        if (!seenIds.has(listing.id)) {
            uniqueListings.push(listing);
            seenIds.add(listing.id);
        }
    }

    return uniqueListings;
}
