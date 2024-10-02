import { render, screen } from '@testing-library/svelte';
import { expect, test, vi, describe } from 'vitest';
import LoginPage from './+page.svelte';

const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('render test', () => {
    test('all components should render correctly', async () => {
        mockFetch.mockResolvedValueOnce({
            ok: true,
            json: () => {}
        });

        render(LoginPage);
        const newPassword = screen.getByLabelText('Email');
        expect(newPassword).toBeDefined();
        const confirmPassword = screen.getByText('Password');
        expect(confirmPassword).toBeDefined();

        const button = screen.getByText('Login');
        expect(button).toBeDefined();
    });
});
