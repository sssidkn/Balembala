import type { Cookies } from '@sveltejs/kit';

export const ACCESS_TOKEN = 'token';
export function setAuth(cookies: Cookies, token: string) {
	const expirationDate = new Date(Date.now() + 60 * 60 * 24 * 1000);
	cookies.set(ACCESS_TOKEN, token, { path: '/', expires: expirationDate });
}
export function getAuth(cookies: Cookies): string | undefined {
	return cookies.get(ACCESS_TOKEN);
}
export function deleteAuth(cookies: Cookies) {
	cookies.delete(ACCESS_TOKEN, { path: '/' });
}
