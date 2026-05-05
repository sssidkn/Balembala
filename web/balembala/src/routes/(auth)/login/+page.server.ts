import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';
export const actions = {
    default: async ({request}) => {
        const data = await request.formData();
        const email = data.get("email");
        const password = data.get("password");
        if (password === "error") {
            return fail(400, {email, error: "Generic error"})
        }
        console.log({email, password})
        return { success: true }
    }
} satisfies Actions;