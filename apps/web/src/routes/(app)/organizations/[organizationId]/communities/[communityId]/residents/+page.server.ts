import { ApiError } from '$lib/api/http';
import type { Household, Resident, Unit } from '$lib/api/types';
import { apiServerRequest } from '$lib/server/api';
import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const [residents, units, households] = await Promise.all([
			apiServerRequest<Resident[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/residents`
			),
			apiServerRequest<Unit[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/units`
			),
			apiServerRequest<Household[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/households`
			)
		]);

		return { residents, units, households, apiError: null };
	} catch (error) {
		return {
			residents: [],
			units: [],
			households: [],
			apiError: error instanceof ApiError ? error.message : 'Unable to load residents.'
		};
	}
};

export const actions: Actions = {
	create: async (event) => {
		const form = await event.request.formData();
		const first_name = String(form.get('first_name') ?? '').trim();
		const last_name = String(form.get('last_name') ?? '').trim();
		const unit_id = String(form.get('unit_id') ?? '').trim();
		const household_id = String(form.get('household_id') ?? '').trim();
		const email = String(form.get('email') ?? '').trim();
		const phone = String(form.get('phone') ?? '').trim();
		const move_in_date = String(form.get('move_in_date') ?? '').trim();
		const resident_type = String(form.get('resident_type') ?? '').trim().toUpperCase();
		const household_role = String(form.get('household_role') ?? '').trim().toUpperCase();

		const values = {
			first_name,
			last_name,
			unit_id,
			household_id,
			email,
			phone,
			move_in_date,
			resident_type,
			household_role
		};

		if (!first_name || !last_name || !unit_id || !resident_type) {
			return fail(400, {
				createError: 'First name, last name, unit, and role are required.',
				values
			});
		}

		try {
			const households = await apiServerRequest<Household[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/households`
			);

			const unitHouseholds = households.filter((household) => household.unit_id === unit_id);
			const resolvedHouseholdId =
				household_id || (unitHouseholds.length === 1 ? unitHouseholds[0].id : '');

			if (unitHouseholds.length > 1 && !resolvedHouseholdId) {
				return fail(400, {
					createError: 'Select a household for this unit.',
					values
				});
			}

			if (
				resolvedHouseholdId &&
				!unitHouseholds.some((household) => household.id === resolvedHouseholdId)
			) {
				return fail(400, {
					createError: 'The selected household does not belong to that unit.',
					values
				});
			}

			await apiServerRequest<Resident>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/residents`,
				{
					method: 'POST',
					body: JSON.stringify({
						first_name,
						last_name,
						unit_id,
						household_id: resolvedHouseholdId || null,
						email: email || null,
						phone: phone || null,
						move_in_date: move_in_date || null,
						resident_type,
						household_role: household_role || null,
						status: 'ACTIVE'
					})
				}
			);

			return {
				createSuccess: `Resident ${first_name} ${last_name} created successfully.`
			};
		} catch (error) {
			return fail(error instanceof ApiError ? error.status : 500, {
				createError: error instanceof ApiError ? error.message : 'Unable to create resident.',
				values
			});
		}
	}
};
