import {
    randomString,
    randomItem,
} from 'https://jslib.k6.io/k6-utils/1.6.0/index.js';
import { check, group, fail } from 'k6';
import type { Options } from 'k6/options';
import http from 'k6/http';
import { createClient } from '../api/client';
import type { paths } from '../api/parkeasy';
import { BASE_URL } from '../constants';

const api = createClient<paths>({ baseUrl: BASE_URL });

interface Data {
    car_id: string;
    cookies: http.CookieJarCookies;
}

// Function to create a new user and log in
export function setup(): Data {
    const randomEmail = `user-${randomString(8)}@example.com`;
    const { response: resp } = api.post('/user', {
        body: {
            email: randomEmail,
            password: 'password123',
            full_name: `Booker ${randomString(8)}`,
        },
    });
    check(resp, {
        '201': (r) => r.status === 201 && r.cookies['session'].length > 0,
    });

    const jar = http.cookieJar();
    const cookiesForURL = jar.cookiesForURL(BASE_URL);
    check(null, {
        "vu jar has cookie 'session'": () =>
            (cookiesForURL.session?.length ?? 0) > 0,
    });

    const car_id = group('Create car', () => {
        const { data, response: resp } = api.post('/cars', {
            body: {
                license_plate: `ABC ${randomString(3)}`,
                make: 'Toyota',
                model: 'Camry',
                color: 'Blue',
            },
        });
        check(resp, {
            '201': (r) => r.status === 201,
        });
        if (
            !check(data, {
                'has id': (d) => d !== undefined && d.id !== '',
            })
        ) {
            fail('Could not create a car');
        }
        return data?.id ?? '';
    });

    return {
        cookies: cookiesForURL,
        car_id,
    };
}

function getRandomInRange(from: number, to: number, fixed: number): number {
    return Number((Math.random() * (to - from) + from).toFixed(fixed));
}

export const options: Options = {
    thresholds: {
        http_req_duration: ['p(99)<500'], // 99% of requests should be below 300ms
    },

    scenarios: {
        // arbitrary name of scenario
        average_load: {
            executor: 'ramping-vus',
            stages: [
                // ramp up to average load of 100 virtual users
                { duration: '10s', target: 100 },
                // maintain load
                { duration: '50s', target: 100 },
                // ramp down to zero
                { duration: '5s', target: 0 },
            ],
        },
    },
};

export default function (data: Data) {
    const jar = http.cookieJar();
    jar.set(BASE_URL, 'session', data.cookies.session[0]);

    const spot = group('Get spot to book', () => {
        const latitude = getRandomInRange(
            49.79776937321384,
            49.949069874886376,
            5,
        );
        const longitude = getRandomInRange(
            -97.29122835966817,
            -97.0318636215253,
            5,
        );
        const { data, response: resp } = api.get('/spots', {
            params: {
                query: {
                    latitude,
                    longitude,
                    distance: 250,
                },
            },
        });

        check(resp, {
            '200': (r) => r.status === 200,
        });

        if (data === undefined || data.length === 0) {
            return null;
        }

        return randomItem(data);
    });

    if (spot === null) {
        return;
    }

    const slot = group('Get slot to book', () => {
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        const { data, response: resp } = api.get('/spots/{id}/availability', {
            params: {
                path: {
                    id: spot.id,
                },
                query: {
                    availability_start: today.toISOString(),
                },
            },
        });

        check(resp, {
            '200': (r) => r.status === 200,
        });

        if (data === undefined || data.length === 0) {
            return null;
        }
        return data[0];
    });

    if (slot === null) {
        return;
    }

    group('Book it', () => {
        const { error, response: resp } = api.post('/spots/{id}/bookings', {
            params: {
                path: {
                    id: spot.id,
                },
            },
            body: {
                car_id: data.car_id,
                booked_times: [slot],
            },
        });

        check(
            { error, resp },
            {
                'Booked successfully or duplicate': ({ error, resp }) =>
                    resp.status === 201 ||
                    error?.type ===
                        'tag:parkwithease.github.io,2024-10-13:problem:duplicate-entity',
            },
        );
    });

    group('Get bookings', () => {
        const { error, response: resp } = api.get('/user/bookings', {
            params: {
                path: {
                    id: '',
                },
            },
        });

        check(
            { error, resp },
            {
                'Succesffuly retreived bookings': ({ resp }) =>
                    resp.status === 200,
            },
        );
    });
}
