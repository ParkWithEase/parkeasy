import { TimeSlotStatus } from '$lib/enum/timeslot-status';
export const spots_data = [
    {
        id: '1',
        features: {
            charging_station: true,
            plug_in: true,
            shelter: true
        },
        location: {
            city: 'Winnipeg',
            country_code: 'CA',
            latitude: 1,
            longitude: 1,
            postal_code: 'RRRRRR',
            street_address: '230 Portage Avenue'
        },
        isListed: true
    },
    {
        id: '2',
        features: {
            charging_station: true,
            plug_in: false,
            shelter: false
        },
        location: {
            city: 'Winnipeg',
            country_code: 'CA',
            latitude: 1,
            longitude: 1,
            postal_code: 'RRRRRR',
            street_address: '230 Portage Avenue'
        },
        isListed: false
    },
    {
        id: '3',
        features: {
            charging_station: true,
            plug_in: false,
            shelter: true
        },
        location: {
            city: 'Winnipeg',
            country_code: 'CA',
            latitude: 1,
            longitude: 1,
            postal_code: 'RRRRRR',
            street_address: '230 Portage Avenue'
        },
        isListed: true
    }
];

export const spots_time_slot = [
    {
        id: '1',
        time_slots: [
            { date: new Date(Date.UTC(2024,9,19,0,0,0)), segment: 1, status: TimeSlotStatus.AVAILABLE },
            { date: new Date(Date.UTC(2024,9,19,0,0,0)), segment: 2, status: TimeSlotStatus.AVAILABLE },
            { date: new Date(Date.UTC(2024,9,19,0,0,0)), segment: 3, status: TimeSlotStatus.AVAILABLE },
            { date: new Date(Date.UTC(2024,9,19,0,0,0)), segment: 4, status: TimeSlotStatus.AVAILABLE },
            { date: new Date(Date.UTC(2024,9,19,0,0,0)), segment: 5, status: TimeSlotStatus.AVAILABLE },
            { date: new Date(Date.UTC(2024,9,19,0,0,0)), segment: 6, status: TimeSlotStatus.AVAILABLE },
            { date: new Date(Date.UTC(2024,9,19,0,0,0)), segment: 7, status: TimeSlotStatus.AVAILABLE },
            { date: new Date(Date.UTC(2024,9,24,0,0,0)), segment: 1, status: TimeSlotStatus.BOOKED },
            { date: new Date(Date.UTC(2024,9,24,0,0,0)), segment: 2, status: TimeSlotStatus.BOOKED },
            { date: new Date(Date.UTC(2024,9,25,0,0,0)), segment: 1, status: TimeSlotStatus.AUCTIONED },
            { date: new Date(Date.UTC(2024,9,25,0,0,0)), segment: 2, status: TimeSlotStatus.AUCTIONED }
        ]
    }
];
