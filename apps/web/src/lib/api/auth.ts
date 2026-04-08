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

export type AuthResponse = ApiEnvelope<AuthPayload>;

export type LoginPayload = {
	email: string;
	password: string;
};

export type SignupPayload = LoginPayload & {
	first_name: string;
	last_name: string;
	phone: string;
};
