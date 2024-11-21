import { render, screen } from '@testing-library/svelte';
import { expect, test, describe, beforeAll, afterAll } from 'vitest';
import userEvent from '@testing-library/user-event';

import PasswordForgot from './+page.svelte';
import { setupServer } from 'msw/node';
import { afterEach } from 'node:test';
import { http, HttpResponse } from 'msw';
import { BACKEND_SERVER } from '$lib/constants';

const server = setupServer();
const user = userEvent.setup();

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

describe('test password forgot service', () => {
    test('Successful link created when enter a valid password', async () => {
        server.use(
            http.post(`${BACKEND_SERVER}/auth/password:forgot`, () =>
                HttpResponse.json({ password_token: 'tano' }, { status: 200 })
            )
        );

        render(PasswordForgot);
        const emailField = screen.getByLabelText('Email');
        expect(emailField).toBeDefined();

        await user.click(emailField);
        await user.keyboard('random@gamil.com');

        const button = screen.getByRole('button', { name: /submit/i });
        await user.click(button);

        const passwordResetLink = screen.getByRole('link', { name: /password reset /i });
        expect(passwordResetLink).toBeDefined();
    });
});
