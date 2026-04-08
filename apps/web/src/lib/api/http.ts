import { env } from '$env/dynamic/public';
import type { ApiEnvelope } from '$lib/api/types';

export class ApiError extends Error {
	code: string;
	status: number;

	constructor(message: string, code = 'api_error', status = 500) {
		super(message);
		this.name = 'ApiError';
		this.code = code;
		this.status = status;
	}
}

function normalizeBaseUrl() {
	return (env.PUBLIC_API_BASE_URL || 'http://localhost:8080').replace(/\/$/, '');
}

export async function apiRequest<T>(
	fetcher: typeof fetch,
	path: string,
	init?: RequestInit
): Promise<T> {
	const response = await fetcher(`${normalizeBaseUrl()}${path}`, {
		...init,
		headers: {
			'content-type': 'application/json',
			...(init?.headers ?? {})
		}
	});

	const payload = (await response.json().catch(() => null)) as ApiEnvelope<T> | null;

	if (!response.ok) {
		throw new ApiError(
			payload?.error?.message ?? `Request failed with status ${response.status}`,
			payload?.error?.code ?? 'request_failed',
			response.status
		);
	}

	if (!payload) {
		throw new ApiError('API returned an empty response');
	}

	return payload.data;
}
