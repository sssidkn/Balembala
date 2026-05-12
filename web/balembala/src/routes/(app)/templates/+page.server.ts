import type { Template } from '$lib/data';
import type { Actions } from './$types';

function template(id: string, title: string, message: string): Template {
	return { id, title, message };
}

const data = [
	template('0', 'First', 'message1'),
	template('1', 'Second', 'message2'),
	template('2', 'Third', 'message3'),
	template('3', 'Forth', 'message4'),
	template('4', 'Fifth', 'message5')
];

export const load = async () => {
	return { data };
};

export const actions = {
	delete: async ({ request }) => {
		const data = await request.formData();
		const id = data.get('id');
		// TODO: Подрубить бэк
		console.log(`Deleting contact with id: ${id}`);
		return { success: true };
	},
	send: async ({ request }) => {
		const data = await request.formData();
		const id = data.get('id');
		// TODO: Подрубить бэк
		console.log(`Sending message to: ${id}`);
		return { success: true };
	}
} satisfies Actions;
