import { expect, test } from 'playwright-test-coverage';

test('home page has expected h1', async ({ page }) => {
    await page.goto('/auth/password-reset');
    await expect(page.locator('h1')).toBeVisible();
});
