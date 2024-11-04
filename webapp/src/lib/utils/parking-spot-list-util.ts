import type { components } from '$lib/sdk/schema';
type ParkingSpot = components['schemas']['ParkingSpot'];

export function coalesceListings(listings: ParkingSpot[]): ParkingSpot[] {
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
