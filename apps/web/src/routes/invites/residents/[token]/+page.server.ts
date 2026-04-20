import type { AuthPayload, LoginPayload, SignupPayload } from '$lib/api/auth';
import { ApiError } from '$lib/api/http';
import type { Resident, ResidentInvitationPreview } from '$lib/api/types';
import { apiServerRequest, apiServerResponse, syncSessionCookie } from '$lib/server/api';
import { fail, redirect, type RequestEvent } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	const parent = await event.parent();

	try {
		const invite = await apiServerRequest<ResidentInvitationPreview>(
			event,
			`/v1/resident-invitations/${event.params.token}`
		);

		return {
			invite,
			auth: parent.auth
		};
	} catch (error) {
		return {
			invite: null,
			auth: parent.auth,
			apiError: error instanceof ApiError ? error.message : 'Unable to load resident invitation.'
		};
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
			return fail(await toStatus(response), {
				loginError: await toMessage(response, 'Unable to log in.'),
				loginValues: { email: payload.email }
			});
		}

		syncSessionCookie(event.cookies, response);
		return acceptAfterAuth(event);
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
			return fail(await toStatus(response), {
				signupError: await toMessage(response, 'Unable to create account.'),
				signupValues: {
					email: payload.email,
					first_name: payload.first_name,
					last_name: payload.last_name,
					phone: payload.phone
				}
			});
		}

		syncSessionCookie(event.cookies, response);
		return acceptAfterAuth(event);
	},
	accept: async (event) => {
		return acceptAfterAuth(event);
	}
};

async function acceptAfterAuth(event: RequestEvent) {
	const preview = await apiServerRequest<ResidentInvitationPreview>(
		event,
		`/v1/resident-invitations/${event.params.token}`
	);

	try {
		await apiServerRequest<Resident>(event, `/v1/resident-invitations/${event.params.token}/accept`, {
			method: 'POST'
		});
		throw redirect(
			303,
			`/organizations/${preview.organization_id}/communities/${preview.community_id}/residents/${preview.resident_id}`
		);
	} catch (error) {
		if (error instanceof ApiError) {
			return fail(error.status, {
				acceptError: error.message
			});
		}

		throw error;
	}
}

async function toMessage(response: Response, fallback: string) {
	const errorPayload = (await response.json().catch(() => null)) as { error?: { message?: string } } | null;
	return errorPayload?.error?.message ?? fallback;
}

function toStatus(response: Response) {
	return response.status || 500;
}
