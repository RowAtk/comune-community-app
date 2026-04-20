import { ApiError } from '$lib/api/http';
import type { Household, Resident, ResidentInvitation, Unit } from '$lib/api/types';
import { apiServerRequest } from '$lib/server/api';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const [resident, invitations] = await Promise.all([
			apiServerRequest<Resident>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/residents/${event.params.residentId}`
			),
			apiServerRequest<ResidentInvitation[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/residents/${event.params.residentId}/invitations`
			)
		]);

		const [unit, household] = await Promise.all([
			apiServerRequest<Unit>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/units/${resident.unit_id}`
			).catch(() => null),
			resident.household_id
				? apiServerRequest<Household>(
						event,
						`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/households/${resident.household_id}`
					).catch(() => null)
				: Promise.resolve(null)
		]);

		return { resident, invitations, unit, household, apiError: null };
	} catch (error) {
		return {
			resident: null,
			invitations: [],
			unit: null,
			household: null,
			apiError: error instanceof ApiError ? error.message : 'Unable to load resident details.'
		};
	}
};

export const actions: Actions = {
	createInvite: async (event) => {
		const form = await event.request.formData();
		const expires_at = String(form.get('expires_at') ?? '').trim();

		try {
			const invitation = await apiServerRequest<ResidentInvitation>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/residents/${event.params.residentId}/invitations`,
				{
					method: 'POST',
					body: JSON.stringify({
						expires_at: expires_at ? new Date(`${expires_at}T23:59:59Z`).toISOString() : null
					})
				}
			);

			return {
				createSuccess: 'Resident invitation generated successfully.',
				createdInvitation: invitation
			};
		} catch (error) {
			return fail(error instanceof ApiError ? error.status : 500, {
				createError:
					error instanceof ApiError ? error.message : 'Unable to generate resident invitation.',
				values: { expires_at }
			});
		}
	},
	unlinkUser: async (event) => {
		try {
			await apiServerRequest<Resident>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/residents/${event.params.residentId}/unlink-user`,
				{
					method: 'POST'
				}
			);

			throw redirect(303, event.url.pathname);
		} catch (error) {
			if (error instanceof ApiError) {
				return fail(error.status, {
					unlinkError: error.message
				});
			}

			throw error;
		}
	}
};
