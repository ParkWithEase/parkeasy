import { render, screen } from '@testing-library/svelte';
import { test, describe, beforeAll, afterAll } from 'vitest';
import userEvent from '@testing-library/user-event';
import SignUp from './+page.svelte';
import { BACKEND_SERVER, PASSWORD_NOT_MATCH } from '$lib/constants';
import { setupServer } from 'msw/node';
import { afterEach } from 'node:test';
import { http, HttpResponse } from 'msw';

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

describe('Sign in page test', () => {
    test('Successful create account', async () => {
        server.use(http.post(`${BACKEND_SERVER}/user`, () => HttpResponse.json({ status: 200 })));
        render(SignUp);

        const firstName = screen.getByLabelText('First Name');
        const lastName = screen.getByLabelText('Last Name');
        const email = screen.getByLabelText('Email');
        const password = screen.getByLabelText('Password');
        const passwordConfirm = screen.getByLabelText('Conform Password');

        await user.click(firstName);
        await user.keyboard('Robin');

        await user.click(lastName);
        await user.keyboard('Hood');

        await user.click(email);
        await user.keyboard('Robin@Hood');

        await user.click(password);
        await user.keyboard('VeryNicePassword');

        await user.click(passwordConfirm);
        await user.keyboard('VeryNicePassword');

        const submitButton = screen.getByRole('button', { name: /Sign Up/i });
        await user.click(submitButton);

        screen.getByRole('paragraph', { name: /sucess message/i });
    });

    test('password mismatch error should show up as intended', async () => {
        render(SignUp);

        const password = screen.getByLabelText('Password');
        const passwordConfirm = screen.getByLabelText('Conform Password');
        await user.click(password);
        await user.keyboard('VeryNicePassword');

        await user.click(passwordConfirm);
        await user.keyboard('VeryNicePasswordButNotMatching');

        screen.getByText(PASSWORD_NOT_MATCH);
    });
});
