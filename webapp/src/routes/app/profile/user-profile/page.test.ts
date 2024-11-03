import { expect, test, describe, beforeAll, afterAll, afterEach } from 'vitest';
import { http, HttpResponse } from 'msw';
import { load } from './+page';
import { setupServer } from 'msw/node';
import { BACKEND_SERVER } from '$lib/constants';
import type { PageData, PageLoadEvent } from './$types';
import { mock } from 'vitest-mock-extended';

const server = setupServer();

beforeAll(() => {
    // NOTE: server.listen must be called before `createClient` is used to ensure
    // the msw can inject its version of `fetch` to intercept the requests.
    server.listen({
        onUnhandledRequest: (request) => {
            throw new Error(`No request handler found for ${request.method} ${request.url}`);
        }
    });
});

afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe('fetch user information test', () => {
    test('test if information is loaded correctly', async () => {
        const rawData = { email: 'test@gmail.com', full_name: 'name' };

        server.use(
            http.get(`${BACKEND_SERVER}/user`, () => HttpResponse.json(rawData, { status: 200 }))
        );

        const loadEvent = mock<PageLoadEvent>({ fetch: global.fetch });
        const data = (await load(loadEvent)) as PageData;
        expect(data.email).toBe(rawData.email);
        expect(data.full_name).toBe(rawData.full_name);
    });
});
