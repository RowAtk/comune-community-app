import type { AuthPayload, LoginPayload, SignupPayload } from '$lib/api/auth';
import { ApiError } from '$lib/api/http';
import { apiServerRequest, apiServerResponse, clearSessionCookie, syncSessionCookie } from '$lib/server/api';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const load = async (event) => {
	try {
		const auth = await apiServerRequest<AuthPayload>(event, '/v1/auth/me');
		return { auth };
	} catch {
		return { auth: null };
	}
};

export const actions: Actions = {
	login: async (event) => {
		const form = await event.request.formData();
		const payload: LoginPayload = {
			email: String(form.get('email') ?? '').trim(),
			password: String(form.get('password') ?? '')
		};

		if (!payload.email || !payload.password) {
			return fail(400, {
				loginError: 'Email and password are required.',
				loginValues: { email: payload.email }
			});
		}

		const response = await apiServerResponse(event, '/v1/auth/login', {
			method: 'POST',
			body: JSON.stringify(payload)
		});

		if (!response.ok) {
			const errorPayload = (await response.json().catch(() => null)) as
				| { error?: { message?: string } }
				| null;

			return fail(response.status, {
				loginError: errorPayload?.error?.message ?? 'Unable to log in.',
				loginValues: { email: payload.email }
			});
		}

		syncSessionCookie(event.cookies, response);
		return { loginSuccess: true };
	},
	signup: async (event) => {
		const form = await event.request.formData();
		const payload: SignupPayload = {
			email: String(form.get('signup_email') ?? '').trim(),
			password: String(form.get('signup_password') ?? ''),
			first_name: String(form.get('first_name') ?? '').trim(),
			last_name: String(form.get('last_name') ?? '').trim(),
			phone: String(form.get('phone') ?? '').trim()
		};

		if (!payload.email || !payload.password) {
			return fail(400, {
				signupError: 'Email and password are required.',
				signupValues: {
					email: payload.email,
					first_name: payload.first_name,
					last_name: payload.last_name,
					phone: payload.phone
				}
			});
		}

		const response = await apiServerResponse(event, '/v1/auth/signup', {
			method: 'POST',
			body: JSON.stringify(payload)
		});

		if (!response.ok) {
			const errorPayload = (await response.json().catch(() => null)) as
				| { error?: { message?: string } }
				| null;

			return fail(response.status, {
				signupError: errorPayload?.error?.message ?? 'Unable to create account.',
				signupValues: {
					email: payload.email,
					first_name: payload.first_name,
					last_name: payload.last_name,
					phone: payload.phone
				}
			});
		}

		syncSessionCookie(event.cookies, response);
		return { signupSuccess: true };
	},
	logout: async (event) => {
		try {
			await apiServerResponse(event, '/v1/auth/logout', { method: 'POST' });
		} finally {
			clearSessionCookie(event.cookies);
		}

		return { logoutSuccess: true };
	},
	openWorkspace: async ({ request }) => {
		const form = await request.formData();
		const organizationId = String(form.get('organizationId') ?? '').trim();
		const communityId = String(form.get('communityId') ?? '').trim();

		if (!organizationId || !communityId) {
			return { error: 'Both organization ID and community ID are required.' };
		}

		throw redirect(
			303,
			`/organizations/${encodeURIComponent(organizationId)}/communities/${encodeURIComponent(communityId)}`
		);
	}
};
