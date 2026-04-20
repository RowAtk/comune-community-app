export type JsonObject = Record<string, unknown>;

export type ApiEnvelope<T> = {
	data: T;
	error?: {
		code: string;
		message: string;
	};
};

export type Community = {
	id: string;
	organization_id: string;
	name: string;
	slug: string;
	address?: string;
	timezone: string;
	status: string;
	settings: JsonObject | null;
	created_at: string;
	updated_at: string;
};

export interface CommunityResource {
	id: string;
	organization_id: string;
	community_id: string;
}

export type Unit = CommunityResource & {
	unit_number: string;
	block_floor?: string;
	unit_type?: string;
	status: string;
	created_at: string;
	updated_at: string;
	deleted_at?: string | null;
};

export type Household = CommunityResource & {
	unit_id: string;
	name?: string;
	created_at: string;
	updated_at: string;
};

export type Resident = CommunityResource & {
	unit_id: string;
	user_id?: string | null;
	linked_user_email?: string;
	linked_user_name?: string;
	linked_user_phone?: string;
	household_id?: string | null;
	first_name: string;
	last_name: string;
	email?: string;
	phone?: string;
	resident_type: string;
	household_role?: string | null;
	status: string;
	is_primary_contact: boolean;
	move_in_date?: string | null;
	move_out_date?: string | null;
	created_at: string;
	updated_at: string;
	deleted_at?: string | null;
};

export type ResidentInvitation = CommunityResource & {
	id: string;
	resident_id: string;
	token: string;
	expires_at: string;
	accepted_at?: string | null;
	invited_by?: string | null;
	accepted_by_user_id?: string | null;
	created_at: string;
};

export type ResidentInvitationPreview = {
	id: string;
	organization_id: string;
	organization: string;
	community_id: string;
	community: string;
	resident_id: string;
	resident_name: string;
	resident_type: string;
	household_id?: string | null;
	household_name?: string;
	unit_id: string;
	unit_number: string;
	household_role?: string | null;
	status: string;
	token: string;
	expires_at: string;
	accepted_at?: string | null;
	invited_by?: string | null;
};
