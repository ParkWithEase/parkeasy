import { spots_data, spots_time_slot } from "../../mock_data";
import type { PageLoad } from "./$types";


export const load: PageLoad = async ({ fetch, params }) => {
    return {
        spot: spots_data.find((spot) => spot.id === params.id),
        timeslots: spots_time_slot.find((timeslots) => timeslots.id === params.id)?.time_slots
    }
}