import { render } from '@testing-library/svelte';
import { expect, test, vi, describe, beforeAll, afterAll, afterEach } from 'vitest';
import UserProfilePage from './+page.svelte';
import { http, HttpResponse } from "msw";
import {load} from './+page';
import { setupServer } from "msw/node";
import { BACKEND_SERVER } from '$lib/constants';

const server = setupServer();

beforeAll(() => {
    // NOTE: server.listen must be called before `createClient` is used to ensure
    // the msw can inject its version of `fetch` to intercept the requests.
    server.listen({
      onUnhandledRequest: (request) => {
        throw new Error(
          `No request handler found for ${request.method} ${request.url}`
        );
      },
    });
  });

afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("fetch user information test", () => {
    test('test if information is loaded correctly', async () => {
        const raw_data = {email: "test@gmail.com", full_name: "name"};

        server.use(
            http.get(`${BACKEND_SERVER}/user`, () => 
                HttpResponse.json(raw_data, { status: 200 }))
        );

        const data = await load({fetch : global.fetch});
        expect(data.email).toBe(raw_data.email);
        expect(data.full_name).toBe(raw_data.full_name);
    })
})