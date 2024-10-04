import { render, screen } from '@testing-library/svelte';
import { expect, test, vi, describe } from 'vitest';
import userEvent from '@testing-library/user-event';

import PasswordReset from './password-reset.svelte';
import { PASSWORD_NOT_MATCH } from '$lib/constants';

const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('fetchData test', () => {
    test('No error for matching password', async () => {
        mockFetch.mockResolvedValueOnce({
            ok: true,
            json: () => {}
        });

        render(PasswordReset, { resetToken: 'random' });
        const user = userEvent.setup();
        const newPassword = screen.getByLabelText('New Password');
        expect(newPassword).toBeDefined();
        const confirmPassword = screen.getByLabelText('Confirm Password');
        expect(confirmPassword).toBeDefined();

        //Test match password
        await user.click(newPassword);
        await user.keyboard('random');
        await user.click(confirmPassword);
        await user.keyboard('random');

        let errorMessage = screen.queryByText("password doesn't match");
        expect(errorMessage).toBeNull();

        const button = screen.getByRole('button', { name: /submit/i });
        await user.click(button);

        errorMessage = screen.queryByText('Failure');
        expect(errorMessage).toBeNull();
    });

    test('Error is shown properly for non matching password field', async () => {
        mockFetch.mockResolvedValueOnce({
            ok: false,
            json: () => {}
        });

        render(PasswordReset, { resetToken: 'random' });
        const user = userEvent.setup();
        const newPassword = screen.getByLabelText('New Password');
        expect(newPassword).toBeDefined();
        const confirmPassword = screen.getByLabelText('Confirm Password');
        expect(confirmPassword).toBeDefined();

        //Test unmatch password
        await user.click(newPassword);
        await user.keyboard('random');
        await user.click(confirmPassword);
        await user.keyboard('random unmatch');

        let errorMessage = screen.queryByText(PASSWORD_NOT_MATCH);
        expect(errorMessage).toBeDefined();

        const button = screen.getByRole('button', { name: /submit/i });
        await user.click(button);

        errorMessage = screen.queryByText('Failure');
        expect(errorMessage).toBeDefined();
    });
});
