import { randomString } from 'https://jslib.k6.io/k6-utils/1.6.0/index.js';
import { check, group } from 'k6';
import type { Options } from 'k6/options';
import { createClient } from './api/client';
import type { paths } from './api/parkeasy';
import { BASE_URL } from './constants';

const api = createClient<paths>({ baseUrl: BASE_URL });

export const options: Options = {
    thresholds: {
        http_req_failed: ['rate<0.01'], // http errors should be less than 1%
        http_req_duration: ['p(99)<700'], // 99% of requests should be below 700ms
    },

    scenarios: {
        // arbitrary name of scenario
        average_load: {
            executor: 'ramping-vus',
            stages: [
                // ramp up to average load of 20 virtual users
                { duration: '10s', target: 20 },
                // maintain load
                { duration: '50s', target: 20 },
                // ramp down to zero
                { duration: '5s', target: 0 },
            ],
        },
    },
};

export default function () {
    const uid = `user-${randomString(8)}`;
    const randomEmail = `${uid}@example.com`;
    const userData = {
        email: randomEmail,
        full_name: `Test User ${randomString(4)}`,
        password: randomString(8),
    };

    group('Create user', () => {
        const { response: resp } = api.post('/user', {
            body: userData,
        });
        check(resp, {
            '201': (r) => r.status === 201 && r.cookies['session'].length > 0,
        });
    });

    group('Get profile', () => {
        const { data: profile, response: resp } = api.get('/user', {});
        check(resp, {
            '200': (r) => r.status == 200,
        });

        check(profile, {
            'Correct name': (p) => p?.full_name === userData.full_name,
            'Correct email': (p) => p?.email === userData.email,
        });
    });
}
