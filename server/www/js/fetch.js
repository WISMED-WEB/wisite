const get = async (url, token) => {
    try {
        if (token !== undefined) {
            const header = {
                method: 'GET',
                withCredentials: true,
                credentials: 'include',
                headers: {
                    'Authorization': 'Bearer ' + token,
                    'Content-Type': 'application/json'
                }
            }
            return await fetch(url, header);
        }
        return await fetch(url);
    } catch (error) {
        console.log(error);
    }
}

export const get_json = async (url, token) => {
    const resp = get(url, token)
    return (await resp).json()
}

export const get_text = async (url, token) => {
    const resp = get(url, token)
    return (await resp).text()    
}