function contact(id: number, name: string, email: string) {
	return { id, name, email };
}

const data = [
	contact(0, 'First', 'first@mail.ru'),
	contact(1, 'Second', 'second@mail.ru'),
	contact(2, 'Third', 'third@mail.ru'),
	contact(3, 'Forth', 'forth@mail.ru'),
	contact(4, 'Fifth', 'fifth@mail.ru')
];

export const load = async () => {
	return { data };
};
