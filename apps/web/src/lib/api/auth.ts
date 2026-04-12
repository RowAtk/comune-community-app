import type { ApiEnvelope } from '$lib/api/types';

export type AuthUser = {
	id: string;
	email: string;
	first_name?: string;
	last_name?: string;
	phone?: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
	deleted_at?: string | null;
};

export type AuthAccount = {
	id: string;
	user_id: string;
	provider: string;
	provider_user_id: string;
	email_at_provider?: string;
	last_login_at?: string;
	created_at: string;
	updated_at: string;
};

export type AuthSessionView = {
	created_at: string;
	expires_at: string;
};

export type AuthPayload = {
	user: AuthUser;
	auth_account: AuthAccount;
	session: AuthSessionView;
};

export type DashboardOrganizationMembership = {
	organization_id: string;
	name: string;
	slug: string;
	role: string;
	status: string;
	joined_at: string;
};

export type DashboardCommunityMembership = {
	organization_id: string;
	organization_name: string;
	organization_slug: string;
	community_id: string;
	community_name: string;
	community_slug: string;
	role: string;
	status: string;
	joined_at: string;
};

export type DashboardPayload = {
	organizations: DashboardOrganizationMembership[];
	communities: DashboardCommunityMembership[];
};

export type OrganizationCreateValues = {
	name: string;
	slug: string;
	legal_name: string;
	billing_email: string;
	phone: string;
	country_code: string;
	timezone: string;
};

export type CommunityCreateValues = {
	organization_id: string;
	name: string;
	slug: string;
	address: string;
	timezone: string;
};

export type AuthResponse = ApiEnvelope<AuthPayload>;
export type DashboardResponse = ApiEnvelope<DashboardPayload>;

export type LoginPayload = {
	email: string;
	password: string;
};

export type SignupPayload = LoginPayload & {
	first_name: string;
	last_name: string;
	phone: string;
};
