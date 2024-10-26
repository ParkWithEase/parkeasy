import type { PageLoad } from './$types';
const test_data = [
    {
        id: '1',
        spot: {
            street_address: '230 Portage',
            postal_code: 'RRRRRR',
            city: 'Winnipeg',
            country_code: 'CA',
            state: 'MB'
        },
        car: {
            license_plate: 'ABCDE',
            color: 'Red',
            model: 'Honda',
            make: 'random make'
        },
        time: {
            date: new Date('2024-10-1'),
            start: '12:00',
            end: '14:00'
        },
        paid_amount: 100
    },
    {
        id: '2',
        spot: {
            street_address: '230 Portage',
            postal_code: 'RRRRRR',
            city: 'Winnipeg',
            country_code: 'CA',
            state: 'MB'
        },
        car: {
            license_plate: 'VCDES',
            color: 'Red',
            model: 'Honda',
            make: 'random make'
        },
        time: {
            date: new Date('2024-10-4'),
            start: '2:00',
            end: '14:00'
        },
        paid_amount: 150
    },
    {
        id: '3',
        spot: {
            street_address: '230 Portage',
            postal_code: 'RRRRRR',
            city: 'Winnipeg',
            country_code: 'CA',
            state: 'MB'
        },
        car: {
            license_plate: 'RacerX',
            color: 'Blue',
            model: 'Honda',
            make: 'random make'
        },
        time: {
            date: new Date('2024-10-10'),
            start: '12:00',
            end: '14:00'
        },
        paid_amount: 240
    }
];

export const load: PageLoad = async ({ fetch }) => {
    return {
        transactions: test_data,
        hasNext: undefined,
        paging: undefined
    };
};
