import { render, screen } from '@testing-library/svelte';
import { expect, test, vi, describe } from 'vitest';
import userEvent from '@testing-library/user-event';

import PasswordForgot from './+page.svelte';

const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('test password forgot service', () => {
    test('Successful link created when enter a valid password', async () => {
        mockFetch.mockResolvedValueOnce({
            ok: true,
            json() {
                return {
                    password_token: 'tano'
                };
            }
        });
        render(PasswordForgot);
        const user = userEvent.setup();
        const emailField = screen.getByLabelText('Email');
        expect(emailField).toBeDefined();

        await user.click(emailField);
        await user.keyboard('random@gamil.com');

        const button = screen.getByRole('button', { name: /submit/i });
        await user.click(button);

        const passwordResetLink = screen.getByRole('link', { name: /password reset/i });
        expect(passwordResetLink).toBeDefined();
    });
});
