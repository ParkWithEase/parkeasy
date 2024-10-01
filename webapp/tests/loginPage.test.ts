import { render, screen } from '@testing-library/svelte';
import { expect, test, vi, describe } from 'vitest';
import LoginPage from '../src/routes/auth/login/+page.svelte';

const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('render test', () => {
    test('all components should render correctly', async () => {
        mockFetch.mockResolvedValueOnce({
            ok: true,
            json: () => {}
        });

        render(LoginPage);
        const newPassword = screen.getByLabelText('Username');
        expect(newPassword).toBeDefined();
        const confirmPassword = screen.getByText('Password');
        expect(confirmPassword).toBeDefined();

        const button = screen.getByText('Submit');
        expect(button).toBeDefined();
    });
});
