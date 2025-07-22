import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
	try {
		const roomId = params.id;
		const storedRoom = sessionStorage.getItem(`room_${roomId}`);

		if (storedRoom) {
			console.log('Room data found in sessionStorage');
			return { room: JSON.parse(storedRoom), error: null };
		}

		const res = await fetch(import.meta.env.VITE_SOCKET_URL + `/api/v1/rooms/${roomId}`);

		if (!res.ok) {
			const errorData = await res
				.json()
				.catch(() => ({ message: `HTTP error! status: ${res.status}` }));
			return { room: null, error: errorData.error || errorData.message };
		}

		const room = await res.json();
		return { room, error: null };
	} catch (e) {
		return { room: null, error: e instanceof Error ? e.message : 'An unknown error occurred' };
	}
};
