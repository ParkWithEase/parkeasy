import { randomString } from 'https://jslib.k6.io/k6-utils/1.6.0/index.js';
import { check, group } from 'k6';
import type { Options } from 'k6/options';
import http from 'k6/http';
import { createClient } from './api/client';
import type { paths } from './api/parkeasy';
import { BASE_URL } from './constants';

const api = createClient<paths>({ baseUrl: BASE_URL });

interface Data {
    cookies: http.CookieJarCookies;
}

// Function to create a new user and log in
export function setup(): Data {
    const randomEmail = `user-${randomString(8)}@example.com`;
    const { response: resp } = api.post('/user', {
        body: {
            email: randomEmail,
            password: 'password123',
            full_name: 'Test User',
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

    return { cookies: cookiesForURL };
}

export const options: Options = {
    thresholds: {
        http_req_failed: ['rate<0.01'], // http errors should be less than 1%
        http_req_duration: ['p(99)<300'], // 99% of requests should be below 300ms
    },

    scenarios: {
        // arbitrary name of scenario
        average_load: {
            executor: 'ramping-vus',
            stages: [
                // ramp up to average load of 20 virtual users
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

    group('Create car', () => {
        const { response: resp } = api.post('/cars', {
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
    });

    group('Get 10 car', () => {
        const { data, response: resp } = api.get('/cars', {
            params: {
                query: {
                    count: 10,
                },
            },
        });

        check(resp, {
            '200': (r) => r.status === 200,
        });

        check(data, {
            'at least one car': (d) => d !== undefined && d.length >= 1,
        });
    });
}
