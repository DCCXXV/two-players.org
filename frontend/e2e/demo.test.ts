import { expect, test } from '@playwright/test';

test('home page has expected h1', async ({ page }) => {
	await page.goto('/');
	await expect(page.locator('h1')).toBeVisible();
});

test('has title', async ({ page }) => {
	await page.goto('/');
	await expect(page).toHaveTitle(/two-players.org/);
});
