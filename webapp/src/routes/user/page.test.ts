import { render, screen, act } from '@testing-library/svelte';
import { expect, test, vi, describe } from 'vitest';

import Page from './+page.svelte';

const mockedGetUser = () =>
    Promise.resolve({
        ok: true,
        json() {
            return {
                full_name: 'tano',
                email: 'test@email.com'
            };
        }
    });

describe('fetchData test', () => {
    const fetchSpy = vi.spyOn(global, 'fetch');
    fetchSpy.mockImplementation(mockedGetUser);
    test('The page should get the correct data', async () => {
        render(Page);
        await act();
        const fullName = screen.getByText('tano');
        expect(fullName).toBeDefined();
        const email = screen.getByText('test@email.com');
        expect(email).toBeDefined();
    });
});
