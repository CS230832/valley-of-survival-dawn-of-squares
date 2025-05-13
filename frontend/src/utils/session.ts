export const SESSION_KEY = 'vosdos-session-token'

export const getSessionID = () => localStorage.getItem(SESSION_KEY)
export const setSessionID = (id: string) => localStorage.setItem(SESSION_KEY, id)
export const removeSessionID = () => localStorage.removeItem(SESSION_KEY)