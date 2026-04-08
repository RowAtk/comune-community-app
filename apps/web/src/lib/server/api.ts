import { env } from '$env/dynamic/public';
import { ApiError } from '$lib/api/http';
import type { ApiEnvelope } from '$lib/api/types';
import type { Cookies, RequestEvent } from '@sveltejs/kit';

const SESSION_COOKIE_NAME = 'comune_session';

function normalizeBaseUrl() {
	return (env.PUBLIC_API_BASE_URL || 'http://localhost:8080').replace(/\/$/, '');
}

function buildHeaders(event: RequestEvent, headers?: HeadersInit) {
	const requestHeaders = new Headers(headers);
	const sessionToken = event.cookies.get(SESSION_COOKIE_NAME);

	if (!requestHeaders.has('content-type')) {
		requestHeaders.set('content-type', 'application/json');
	}

	if (sessionToken) {
		requestHeaders.set('cookie', `${SESSION_COOKIE_NAME}=${sessionToken}`);
	}

	return requestHeaders;
}

export async function apiServerRequest<T>(
	event: RequestEvent,
	path: string,
	init?: RequestInit
): Promise<T> {
	const response = await event.fetch(`${normalizeBaseUrl()}${path}`, {
		...init,
		headers: buildHeaders(event, init?.headers)
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

export async function apiServerResponse(
	event: RequestEvent,
	path: string,
	init?: RequestInit
): Promise<Response> {
	return event.fetch(`${normalizeBaseUrl()}${path}`, {
		...init,
		headers: buildHeaders(event, init?.headers)
	});
}

export function syncSessionCookie(cookies: Cookies, response: Response) {
	const setCookie = response.headers.get('set-cookie');
	if (!setCookie) {
		return;
	}

	const parsed = parseSessionCookie(setCookie);
	if (!parsed) {
		return;
	}

	cookies.set(SESSION_COOKIE_NAME, parsed.value, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure: false,
		expires: parsed.expires
	});
}

export function clearSessionCookie(cookies: Cookies) {
	cookies.delete(SESSION_COOKIE_NAME, {
		path: '/'
	});
}

function parseSessionCookie(header: string) {
	const [cookiePart, ...attributes] = header.split(';').map((part) => part.trim());
	const [name, ...rest] = cookiePart.split('=');
	if (name !== SESSION_COOKIE_NAME || rest.length === 0) {
		return null;
	}

	const expiresAttribute = attributes.find((attribute) =>
		attribute.toLowerCase().startsWith('expires=')
	);

	return {
		value: rest.join('='),
		expires: expiresAttribute ? new Date(expiresAttribute.slice(8)) : undefined
	};
}
