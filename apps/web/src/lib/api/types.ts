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
	household_id?: string | null;
	first_name: string;
	last_name: string;
	email?: string;
	phone?: string;
	resident_type: string;
	status: string;
	is_primary_contact: boolean;
	move_in_date?: string | null;
	move_out_date?: string | null;
	created_at: string;
	updated_at: string;
	deleted_at?: string | null;
};
