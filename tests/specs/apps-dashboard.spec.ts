import test, { expect } from '@playwright/test';
import authUtil from 'utils/auth.util';
import { oidcClients } from '../data';
import { cleanupBackend } from '../utils/cleanup.util';

test.beforeEach(async () => await cleanupBackend());

test('Dashboard shows all clients in the correct order', async ({ page }) => {
	const client1 = oidcClients.tailscale;
	const client2 = oidcClients.nextcloud;

	await page.goto('/settings/apps');

	await expect(page.getByTestId('authorized-oidc-client-card')).toHaveCount(6);

	// Should be first
	const card1 = page.getByTestId('authorized-oidc-client-card').first();

	await expect(card1.getByRole('heading')).toHaveText(client1.name);

	const card2 = page.getByTestId('authorized-oidc-client-card').nth(1);
	await expect(card2.getByRole('heading', { name: client2.name })).toBeVisible();
	await expect(card2.getByText(new URL(client2.launchURL).hostname)).toBeVisible();
});

test.describe('Dashboard shows only clients where user has access', () => {
	test("User can't see one client", async ({ page }) => {
		await authUtil.changeUser(page, 'craig');
		const notVisibleClient = oidcClients.immich;

		await page.goto('/settings/apps');

		const cards = page.getByTestId('authorized-oidc-client-card');

		await expect(cards).toHaveCount(5);

		const cardTexts = await cards.allTextContents();
		expect(cardTexts.some((text) => text.includes(notVisibleClient.name))).toBe(false);
	});
	test('User can see all clients', async ({ page }) => {
		await page.goto('/settings/apps');
		const cards = page.getByTestId('authorized-oidc-client-card');
		await expect(cards).toHaveCount(6);
	});
});

test('Revoke authorized client', async ({ page }) => {
	const client = oidcClients.tailscale;

	await page.goto('/settings/apps');

	const card = page.getByTestId('authorized-oidc-client-card').filter({ hasText: client.name });

	card.getByRole('button', { name: 'Toggle menu' }).click();

	await page.getByRole('menuitem', { name: 'Revoke' }).click();
	await page.getByRole('button', { name: 'Revoke' }).click();

	await expect(page.locator('[data-type="success"]')).toHaveText(
		`The access to ${client.name} has been successfully revoked.`
	);

	// The ... ago text should be gone as there is no last access anymore
	await expect(card).not.toContainText('ago');
});

test('Launch authorized client', async ({ page }) => {
	const client = oidcClients.nextcloud;

	await page.goto('/settings/apps');

	const card1 = page.getByTestId('authorized-oidc-client-card').first();
	await expect(card1.getByRole('button', { name: 'Launch' })).toBeDisabled();

	const card2 = page.getByTestId('authorized-oidc-client-card').nth(1);
	await expect(card2.getByRole('link', { name: 'Launch' })).toHaveAttribute(
		'href',
		client.launchURL
	);
});

test('Dashboard shows client description when present', async ({ page, request }) => {
	const clientWithDesc = oidcClients.withDescription;

	// Update the client to have a description via API
	await request.put(`/api/oidc/clients/${clientWithDesc.id}`, {
		data: {
			name: clientWithDesc.name,
			description: clientWithDesc.description,
			callbackURLs: [clientWithDesc.callbackUrl],
			logoutCallbackURLs: [],
			isPublic: false,
			pkceEnabled: false,
			requiresReauthentication: false,
			requiresPushedAuthorizationRequests: false,
			credentials: { federatedIdentities: [] },
			isGroupRestricted: false
		}
	});

	await page.goto('/settings/apps');

	// Find the card with the description
	const card = page.getByTestId('authorized-oidc-client-card').filter({ hasText: clientWithDesc.name });

	// Should show description instead of hostname
	await expect(card.getByText(clientWithDesc.description)).toBeVisible();
});

test('Dashboard shows hostname when description is not present', async ({ page }) => {
	const clientWithoutDesc = oidcClients.nextcloud;

	await page.goto('/settings/apps');

	const card = page.getByTestId('authorized-oidc-client-card').filter({ hasText: clientWithoutDesc.name });

	// Should show hostname
	await expect(card.getByText(new URL(clientWithoutDesc.launchURL).hostname)).toBeVisible();
});
