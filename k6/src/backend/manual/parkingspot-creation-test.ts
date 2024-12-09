import http from 'k6/http';
import { SharedArray } from 'k6/data';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.6.0/index.js';
import { check } from 'k6';
import { vu } from 'k6/execution';
import type { Options } from 'k6/options';
import { createClient } from '../api/client';
import type { paths } from '../api/parkeasy';
import { BASE_URL } from '../constants';

const api = createClient<paths>({ baseUrl: BASE_URL });

interface Data {
    cookies: http.CookieJarCookies;
}

const NUMBER_OF_USERS = 20;

export const options: Options = {
    thresholds: {
        http_req_failed: ['rate<0.01'], // http errors should be less than 1%
        http_req_duration: ['p(99)<1000'], // 99% of requests should be below 1s
    },

    scenarios: {
        // arbitrary name of scenario
        average_load: {
            executor: 'per-vu-iterations',
            vus: NUMBER_OF_USERS,
            iterations: 1,
        },
    },
};

const addresses = new SharedArray('addresses', function () {
    const f = JSON.parse(open('./addresses.json'));
    return f;
});

const ENTRIES_PER_USER = addresses.length / NUMBER_OF_USERS;

// Function to create a new user and log in
export function setup(): Data {
    const randomEmail = `landlord-${randomString(8)}@example.com`;
    const { response: resp } = api.post('/user', {
        body: {
            email: randomEmail,
            password: 'password123',
            full_name: `Landlord ${randomString(3)}`,
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

export default function (data: Data) {
    const jar = http.cookieJar();
    jar.set(BASE_URL, 'session', data.cookies.session[0]);

    const offset = (vu.idInTest - 1) * ENTRIES_PER_USER;
    for (
        let i = 0;
        i < ENTRIES_PER_USER && offset + i < addresses.length;
        ++i
    ) {
        const addr = addresses[offset + i];
        const start_time = new Date();
        start_time.setUTCHours(start_time.getUTCHours() + 1, 0);
        const end_time = new Date(start_time);
        end_time.setUTCMinutes(30);
        const { response: resp } = api.post('/spots', {
            body: {
                availability: [
                    {
                        start_time: start_time.toISOString(),
                        end_time: end_time.toISOString(),
                    },
                ],
                location: addr,
                price_per_hour: 1.0,
            },
        });
        check(resp, {
            'Spot created': (r) => {
                const result = r.status === 201;
                if (!result) {
                    console.log(r.json());
                }
                return result;
            },
        });
    }
}
