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
            state: 'MB',
            country_code: 'CA',
            latitude: 49.808856,
            longitude: -97.13329,
            postal_code: 'RRRRRR',
            street_address: 'Engineering Building'
        },
        price_per_hour: 10,
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
            state: 'MB',
            country_code: 'CA',
            latitude: 49.811791,
            longitude: -97.199808,
            postal_code: 'RRRRRR',
            street_address: '230 Portage Avenue'
        },
        price_per_hour: 20,
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
            state: 'MB',
            country_code: 'CA',
            latitude: 49.9,
            longitude: -97,
            postal_code: 'RRRRRR',
            street_address: '230 Portage Avenue'
        },
        isListed: true
    },
    {
        id: '4',
        features: {
            charging_station: true,
            plug_in: false,
            shelter: true
        },
        location: {
            city: 'Winnipeg',
            state: 'MB',
            country_code: 'CA',
            latitude: 49.9,
            longitude: -97.2,
            postal_code: 'RRRRRR',
            street_address: '230 Portage Avenue'
        },
        price_per_hour: 30,
        isListed: true
    },
    {
        id: '5',
        features: {
            charging_station: true,
            plug_in: true,
            shelter: true
        },
        location: {
            city: 'Winnipeg',
            state: 'MB',
            country_code: 'CA',
            latitude: 49.808856,
            longitude: -97.13329,
            postal_code: 'RRRRRR',
            street_address: 'Engineering Building'
        },
        price_per_hour: 10,
        isListed: true
    },
    {
        id: '6',
        features: {
            charging_station: true,
            plug_in: true,
            shelter: true
        },
        location: {
            city: 'Winnipeg',
            state: 'MB',
            country_code: 'CA',
            latitude: 49.808856,
            longitude: -97.13329,
            postal_code: 'RRRRRR',
            street_address: 'Engineering Building'
        },
        price_per_hour: 10,
        isListed: true
    }
];

export const spots_time_slot = [
    {
        id: '1',
        time_slots: [
            {
                date: new Date(2024, 9, 19, 0, 0, 0),
                segment: 1,
                status: TimeSlotStatus.AVAILABLE
            },
            {
                date: new Date(2024, 9, 19, 0, 0, 0),
                segment: 2,
                status: TimeSlotStatus.AVAILABLE
            },
            {
                date: new Date(2024, 9, 19, 0, 0, 0),
                segment: 3,
                status: TimeSlotStatus.AVAILABLE
            },
            {
                date: new Date(2024, 9, 19, 0, 0, 0),
                segment: 4,
                status: TimeSlotStatus.AVAILABLE
            },
            {
                date: new Date(2024, 9, 19, 0, 0, 0),
                segment: 5,
                status: TimeSlotStatus.AVAILABLE
            },
            {
                date: new Date(2024, 9, 19, 0, 0, 0),
                segment: 6,
                status: TimeSlotStatus.AVAILABLE
            },
            {
                date: new Date(2024, 9, 19, 0, 0, 0),
                segment: 7,
                status: TimeSlotStatus.AVAILABLE
            },
            {
                date: new Date(2024, 9, 24, 0, 0, 0),
                segment: 1,
                status: TimeSlotStatus.BOOKED
            },
            {
                date: new Date(2024, 9, 24, 0, 0, 0),
                segment: 2,
                status: TimeSlotStatus.BOOKED
            },
            {
                date: new Date(2024, 9, 25, 0, 0, 0),
                segment: 1,
                status: TimeSlotStatus.AUCTIONED
            },
            {
                date: new Date(2024, 9, 25, 0, 0, 0),
                segment: 2,
                status: TimeSlotStatus.AUCTIONED
            }
        ]
    }
];
