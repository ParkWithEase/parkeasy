import { render, screen } from '@testing-library/svelte';
import { test, vi, describe } from 'vitest';
import userEvent from '@testing-library/user-event';
import SignUp from './+page.svelte';
import { PASSWORD_NOT_MATCH } from '$lib/constants';

const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('Sign in page test', () => {
    test('Successful create account', async () => {
        mockFetch.mockResolvedValueOnce({
            ok: true,
            json: () => {}
        });
        render(SignUp);
        const user = userEvent.setup();

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
        const user = userEvent.setup();

        const password = screen.getByLabelText('Password');
        const passwordConfirm = screen.getByLabelText('Conform Password');
        await user.click(password);
        await user.keyboard('VeryNicePassword');

        await user.click(passwordConfirm);
        await user.keyboard('VeryNicePasswordButNotMatching');

        screen.getByText(PASSWORD_NOT_MATCH);
    });
});
