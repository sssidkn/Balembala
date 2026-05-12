import { bound_contacts, contacts } from '$lib/data';
import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { resolve } from '$app/paths';

export const load = async () => {
	// TODO: бэк. Здесь мы аж два раза в бэк сходим, потому что мне лень общее состояние делать.
	// По хорошему, это надо в layout.server делать. А ещё лучше на коленях перед бэком ползать,
	// чтоб нужный эндпоинт дали. Не хочу. Хочу мороженое.
	return { data: new Set(contacts).difference(new Set(bound_contacts)) };
};

export const actions = {
	default: async ({ request, params }) => {
		const id = params.id;
		const data = await request.formData();
		console.log(`binding to template with ${id}`);
		const contacts = data.get('contacts');
		// TODO: Подрубить бэк
		console.log(contacts);
		redirect(303, resolve(`/templates/${id}`));
	}
} satisfies Actions;
